package homekit

import (
	"log/slog"

	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/service"
	"github.com/waynezhang/homekit-proxy/internal/actionlog"
	"github.com/waynezhang/homekit-proxy/internal/config"
	"github.com/waynezhang/homekit-proxy/internal/homekit/characteristics"
	"github.com/waynezhang/homekit-proxy/internal/homekit/runner"
)

type rootBridge struct {
	b           *accessory.Bridge
	accessories []*accessory.A
	runners     []*runner.CharacteristicRunner
}

func parseConfig(cfg *config.Config, actionLog *actionlog.ActionLog) *rootBridge {
	bridge := accessory.NewBridge(accessory.Info{
		Name: cfg.Bridge.Name,
	})
	bridge.Id = 1

	accessories := []*accessory.A{}
	runners := []*runner.CharacteristicRunner{}
	nextId := 1
	for _, ac := range cfg.Accessories {
		a := accessoryFromConfig(ac)
		accessories = append(accessories, a)

		for _, sc := range ac.Services {
			s := service.New(sc.TypeString)
			a.Id = uint64(ac.Id)
			a.AddS(s)

			for _, cc := range sc.Characteristics {
				c := characteristics.NewCharacteristic(cc)
				if c == nil {
					slog.Error("Unsupported Characteristics type", "type", cc.Type)
					continue
				}

				s.AddC(c)

				name := ac.Name + " - " + cc.Type
				runner := runner.NewCharacteristicRunner(name, ac.Area, ac.Icon, &cc, c, actionLog)
				runner.Id = nextId
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
			Name: ac.Name,
		},
		byte(ac.TypeByte),
	)
}

func automationRunnersFromConfig(cs []*config.AutomationConfig, ntfyConfig *config.NtfyConfig, actionLog *actionlog.ActionLog) []*runner.AutomationRunner {
	runners := []*runner.AutomationRunner{}

	for _, a := range cs {
		r := runner.AutomationRunner{Config: a, ActionLog: actionLog, NtfyTopic: ntfyConfig.Topic}
		r.Config.Enabled = a.Enabled
		runners = append(runners, &r)
	}

	return runners
}
