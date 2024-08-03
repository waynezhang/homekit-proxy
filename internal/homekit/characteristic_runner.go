package homekit

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/brutella/hap/characteristic"
	"github.com/waynezhang/homekit-proxy/internal/config"
	"github.com/waynezhang/homekit-proxy/internal/utils"
)

type characteristicRunner struct {
	name   string
	config *config.CharacteristicsConfig
	c      *characteristic.C
}

func newCharacteristicRunner(name string, config *config.CharacteristicsConfig, c *characteristic.C) *characteristicRunner {
	r := &characteristicRunner{
		name:   name,
		config: config,
		c:      c,
	}

	r.c.OnCValueUpdate(func(c *characteristic.C, new, old interface{}, req *http.Request) {
		slog.Info("Remote value changed on "+r.name, "new", new, "req", req)
		if req == nil {
			return
		}

		param := ""
		switch any(new).(type) {
		case int:
			param = strconv.Itoa(new.(int))
		case string:
			param = new.(string)
		case bool:
			param = strconv.FormatBool(new.(bool))
		case float64:
			param = strconv.FormatFloat(new.(float64), 'f', 2, 64)
		default:
			slog.Error("Unsupported value", "value", new)
		}

		cmd := r.config.Set + " " + param
		utils.Exec(cmd)
	})

	return r
}

func (r *characteristicRunner) start() {
	go func() {
		for {
			slog.Info("Updating status of " + r.name)

			cmd := r.config.Get
			output := utils.Exec(cmd)

			var val interface{}
			if len(output) > 0 {
				switch r.config.Type {
				case "Active":
					if i, e := strconv.Atoi(output); e == nil {
						val = i
					}
				case "On":
					if output == "true" {
						val = true
					} else {
						val = false
					}
				case "CurrentTemperature":
					if f, e := strconv.ParseFloat(output, 64); e == nil {
						val = f
					}
				case "TargetHeaterCoolerState":
					switch output {
					case "Auto":
						val = characteristic.TargetHeaterCoolerStateAuto
					case "Cool":
						val = characteristic.TargetHeaterCoolerStateCool
					case "Heat":
						val = characteristic.TargetHeaterCoolerStateHeat
					}
				case "CurrentHeaterCoolerState":
					switch output {
					case "Inactive":
						val = characteristic.CurrentHeaterCoolerStateInactive
					case "Idle":
						val = characteristic.CurrentHeaterCoolerStateIdle
					case "Heating":
						val = characteristic.CurrentHeaterCoolerStateHeating
					case "Cooling":
						val = characteristic.CurrentHeaterCoolerStateCooling
					}
				default:
				}
				if val != nil {
					slog.Info("Setting value", "name", r.name, "val", val)
					r.c.SetValueRequest(val, nil)
				} else {
					slog.Error("No value parsed")
				}
			}
			time.Sleep(time.Duration(r.config.Poll) * time.Second)
		}
	}()
}
