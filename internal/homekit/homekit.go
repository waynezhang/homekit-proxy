package homekit

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/service"
	"github.com/waynezhang/homekit-proxy/internal/config"
	"github.com/waynezhang/homekit-proxy/internal/utils"
)

type HMManager struct {
	cfg *config.Config
}

func New(cfg *config.Config) *HMManager {
	return &HMManager{cfg: cfg}
}

func (m *HMManager) Start(dbPath string) {
	bridge, accessories, runners := parseConfig(m.cfg)

	store := hap.NewFsStore(dbPath)
	server, err := hap.NewServer(store, bridge.A, accessories...)
	utils.CheckFatalError(err, "Failed to build server")

	iface := os.Getenv("HOMEKIT_PROXY_IFACE")
	if iface != "" {
		slog.Info("Binding to iface", "iface", iface)
		server.Ifaces = []string{iface}
	}

	addr := os.Getenv("HOMEKIT_PROXY_BINDADDR")
	if addr != "" {
		slog.Info("Binding to addr", "addr", addr)
		server.Addr = addr
	}

	server.Pin = m.cfg.Bridge.PinCode
	slog.Info("PIN code", "pin", m.cfg.Bridge.PinCode)

	for _, r := range runners {
		r.start()
	}
	startHealthCheckHandler(server)
	err = server.ListenAndServe(context.Background())
	utils.CheckFatalError(err, "Failed to start server")
}

func startHealthCheckHandler(server *hap.Server) {
	server.ServeMux().HandleFunc("/health", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("OK"))
	})
}

func parseConfig(cfg *config.Config) (*accessory.Bridge, []*accessory.A, []*characteristicRunner) {
	bridge := accessory.NewBridge(accessory.Info{
		Name:         cfg.Bridge.Name,
		Model:        cfg.Bridge.Model,
		Manufacturer: cfg.Bridge.Manufacturer,
		Firmware:     cfg.Bridge.Firmware,
	})
	bridge.Id = 1
	slog.Info("New bridge", "name", bridge.Name())

	accessories := []*accessory.A{}
	runners := []*characteristicRunner{}
	for _, ac := range cfg.Accessories {
		a := accessoryFromConfig(&ac)
		slog.Info("  New accessory", "name", a.Name(), "type", a.Type)
		for _, sc := range ac.Services {
			s := service.New(sc.TypeString)
			a.Id = uint64(ac.Id)
			a.AddS(s)
			slog.Info("    New Service", "name", s.Type)

			for _, cc := range sc.Characteristics {
				c := characteristicFromType(cc.Type)
				if c != nil {
					s.AddC(c)
					runners = append(runners, newCharacteristicRunner(
						ac.Name+" - "+cc.Type,
						&cc,
						c,
					))
					slog.Info("      New Characteristic", "type", cc.Type, "polling", cc.Poll)
				} else {
					slog.Error("Unsupported Characteristics type", "type", cc.Type)
				}
			}
		}
		accessories = append(accessories, a)
	}

	return bridge, accessories, runners
}

func accessoryFromConfig(ac *config.AccessoriesConfig) *accessory.A {
	return accessory.New(
		accessory.Info{
			Name:         ac.Name,
			Model:        ac.Model,
			Manufacturer: ac.Manufacturer,
			Firmware:     ac.Firmware,
		},
		byte(ac.TypeByte),
	)
}
