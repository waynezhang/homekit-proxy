package homekit

import (
	"log/slog"

	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/service"
	"github.com/waynezhang/homekit-proxy/internal/config"
)

type rootBridge struct {
	b           *accessory.Bridge
	accessories []*accessory.A
	runners     []*characteristicRunner
}

func parseConfig(cfg *config.Config) *rootBridge {
	bridge := accessory.NewBridge(accessory.Info{
		Name:         cfg.Bridge.Name,
		Model:        cfg.Bridge.Model,
		Manufacturer: cfg.Bridge.Manufacturer,
		Firmware:     cfg.Bridge.Firmware,
	})
	bridge.Id = 1

	accessories := []*accessory.A{}
	runners := []*characteristicRunner{}
	nextId := 1
	for _, ac := range cfg.Accessories {
		a := accessoryFromConfig(ac)
		accessories = append(accessories, a)

		for _, sc := range ac.Services {
			s := service.New(sc.TypeString)
			a.Id = uint64(ac.Id)
			a.AddS(s)

			for _, cc := range sc.Characteristics {
				c := characteristicFromConfig(cc)
				if c == nil {
					slog.Error("Unsupported Characteristics type", "type", cc.Type)
					continue
				}

				s.AddC(c)

				name := ac.Name + " - " + cc.Type
				runner := newCharacteristicRunner(name, &cc, c)
				runner.id = nextId
				nextId++
				runners = append(runners, runner)
			}
		}
	}

	return &rootBridge{
		b:           bridge,
		accessories: accessories,
		runners:     runners,
	}
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

func automationRunnersFromConfig(cs []*config.AutomationConfig) []*automationRunner {
	runners := []*automationRunner{}

	nextId := 1
	for _, a := range cs {
		r := automationRunner{config: a}
		r.id = nextId
		nextId++
		runners = append(runners, &r)
	}

	return runners
}
