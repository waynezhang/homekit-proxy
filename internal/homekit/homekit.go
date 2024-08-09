package homekit

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/brutella/hap"
	"github.com/waynezhang/homekit-proxy/internal/config"
	"github.com/waynezhang/homekit-proxy/internal/utils"
)

type HMManager struct {
	server      *hap.Server
	root        *rootBridge
	automations []*automationRunner
}

func New(cfg *config.Config, dbPath string) *HMManager {
	root := parseConfig(cfg)
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

func (m *HMManager) Start() {
	for _, r := range m.root.runners {
		r.start()
	}
	for _, r := range m.automations {
		r.start(time.Now())
	}

	m.startHealthCheckHandler()
	m.startAPIHandler()

	err := m.server.ListenAndServe(context.Background())
	utils.CheckFatalError(err, "[Server] Failed to start server")
}
