package homekit

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/brutella/hap"
	"github.com/radovskyb/watcher"
	"github.com/waynezhang/homekit-proxy/internal/config"
	"github.com/waynezhang/homekit-proxy/internal/homekit/runner"
	"github.com/waynezhang/homekit-proxy/internal/utils"
)

type HMManager struct {
	config      *config.Config
	server      *hap.Server
	root        *rootBridge
	automations []*runner.AutomationRunner
	cancel      context.CancelFunc
}

const (
	serverEventServerInit serverEvent = iota + 1
	serverEventServerPreparing
	serverEventServerRunning
	serverEventServerStopped
	serverEventConfigFileChanged
)

type serverEvent int

func Serve(cfgFile string, dbPath string) {
	ch := make(chan serverEvent, 1)
	ch <- serverEventServerInit
	defer close(ch)

	var m *HMManager
	for {
		select {
		case e := <-ch:
			slog.Info("[State] State changed", "state", e)
			switch e {
			case serverEventServerInit:
				go func() {
					startWatchingConfigFile(cfgFile, ch)
				}()
				ch <- serverEventServerPreparing
			case serverEventServerPreparing:
				m = new(cfgFile, dbPath)
				go func() {
					m.start()
					ch <- serverEventServerStopped
				}()
				ch <- serverEventServerRunning
			case serverEventServerStopped:
				ch <- serverEventServerPreparing
			case serverEventConfigFileChanged:
				m.cancel()
			}
		}
	}
}

func new(cfgFile string, dbPath string) *HMManager {
	cfg := config.Parse(cfgFile, dbPath)
	root := parseConfig(&cfg)
	automations := automationRunnersFromConfig(cfg.Automations)

	var w = slog.Info
	w("[Config] Bridge: ", "name", root.b.Name())
	for _, c := range root.runners {
		w("[Config]   Characteristic", "name", c.Name, "get", c.Config.Get, "set", c.Config.Set, "Poll", c.Config.Poll)
	}
	w("[Config] Automations:")
	for _, a := range automations {
		w("[Config]   Rule", "name", a.Config.Name, "cmd", a.Config.Cmd, "cron", a.Config.Cron, "tolerance", a.Config.Tolerance)
	}

	store := hap.NewFsStore(dbPath)
	server, err := hap.NewServer(store, root.b.A, root.accessories...)
	utils.CheckFatalError(err, "[Server] Failed to build server")

	iface := os.Getenv("HOMEKIT_PROXY_IFACE")
	if iface != "" {
		slog.Info("[Server] Binding to iface", "iface", iface)
		server.Ifaces = []string{iface}
	}

	addr := os.Getenv("HOMEKIT_PROXY_BINDADDR")
	if addr != "" {
		slog.Info("[Server] Binding to addr", "addr", addr)
		server.Addr = addr
	}

	server.Pin = cfg.Bridge.PinCode
	slog.Info("[Server] PIN code", "pin", cfg.Bridge.PinCode)

	return &HMManager{
		config:      &cfg,
		server:      server,
		root:        root,
		automations: automations,
	}
}

func (m *HMManager) start() {
	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel

	for _, r := range m.root.runners {
		r.Start(ctx)
	}
	for _, r := range m.automations {
		r.Start(time.Now(), ctx)
	}

	m.startHealthCheckHandler()
	m.startAPIHandler()
	m.startUIHandler()

	err := m.server.ListenAndServe(ctx)
	slog.Info("[Server] Server is stopped", "reason", err)
}

func (m *HMManager) stop() {
	m.cancel()
}

func startWatchingConfigFile(cfgFile string, ch chan serverEvent) {
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Write)
	err := w.Add(cfgFile)
	utils.CheckFatalError(err, "[FS] Failed to watch file")

	go func() {
		for {
			select {
			case event := <-w.Event:
				slog.Info("[FS] Config file changed", "event", event)
				ch <- serverEventConfigFileChanged
			case err := <-w.Error:
				slog.Error("[FS] Received error", "err", err)
			}
		}
	}()

	err = w.Start(time.Millisecond * 100)
	utils.CheckFatalError(err, "[FS] Failed to start watching file")
}
