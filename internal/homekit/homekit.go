package homekit

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/brutella/hap"
	"github.com/fsnotify/fsnotify"
	"github.com/waynezhang/homekit-proxy/internal/config"
	"github.com/waynezhang/homekit-proxy/internal/utils"
)

type HMManager struct {
	server      *hap.Server
	root        *rootBridge
	automations []*automationRunner
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
	cfg := config.Parse(cfgFile)
	root := parseConfig(&cfg)
	automations := automationRunnersFromConfig(cfg.Automations)

	var w = slog.Info
	w("[Config] Bridge: ", "name", root.b.Name())
	for _, c := range root.runners {
		w("[Config]   Characteristic", "name", c.name, "get", c.config.Get, "set", c.config.Set, "Poll", c.config.Poll)
	}
	w("[Config] Automations:")
	for _, a := range automations {
		w("[Config]   Rule", "name", a.config.Name, "cmd", a.config.Cmd, "cron", a.config.Cron, "margin", a.config.Margin)
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
		server:      server,
		root:        root,
		automations: automations,
	}
}

func (m *HMManager) start() {
	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel

	for _, r := range m.root.runners {
		r.start()
	}
	for _, r := range m.automations {
		r.start(time.Now())
	}

	m.startHealthCheckHandler()
	m.startAPIHandler()

	err := m.server.ListenAndServe(ctx)
	slog.Info("[Server] Server is stopped", "reason", err)
}

func (m *HMManager) stop() {
	m.cancel()
}

func startWatchingConfigFile(cfgFile string, ch chan serverEvent) {
	watcher, err := fsnotify.NewWatcher()
	utils.CheckFatalError(err, "[FS] Failed to start filesystem watcher")
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					slog.Info("[FS] Config file changed")
					ch <- serverEventConfigFileChanged
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				slog.Error("[FS] Received error", "err", err)
			}
		}

	}()

	err = watcher.Add(cfgFile)
	utils.CheckFatalError(err, "[FS] Failed to watch file")
	<-make(chan struct{})
}
