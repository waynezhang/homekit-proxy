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

			val := parseValueFromCharacteristicType(output, r.config.Type)
			if val != nil {
				slog.Info("Setting value", "name", r.name, "val", val)
				r.c.SetValueRequest(val, nil)
			} else {
				slog.Error("No value parsed")
			}
			time.Sleep(time.Duration(r.config.Poll) * time.Second)
		}
	}()
}

func parseValueFromCharacteristicType(output string, typ string) any {
	switch typ {
	case "Active":
		i, _ := strconv.Atoi(output)
		return i
	case "On":
		return map[string]interface{}{
			"true":  true,
			"false": false,
		}[output]
	case "CurrentTemperature", "TargetTemperature", "RotationSpeed", "CoolingThresholdTemperature", "HeatingThresholdTemperature":
		f, _ := strconv.ParseFloat(output, 64)
		return f
	case "TargetHeaterCoolerState":
		return map[string]interface{}{
			"TargetHeaterCoolerStateAuto": characteristic.TargetHeaterCoolerStateAuto,
			"TargetHeaterCoolerStateHeat": characteristic.TargetHeaterCoolerStateHeat,
			"TargetHeaterCoolerStateCool": characteristic.TargetHeaterCoolerStateCool,
		}[output]
	case "CurrentHeaterCoolerState":
		return map[string]interface{}{
			"CurrentHeaterCoolerStateInactive": characteristic.CurrentHeaterCoolerStateInactive,
			"CurrentHeaterCoolerStateIdle":     characteristic.CurrentHeaterCoolerStateIdle,
			"CurrentHeaterCoolerStateHeating":  characteristic.CurrentHeaterCoolerStateHeating,
			"CurrentHeaterCoolerStateCooling":  characteristic.CurrentHeaterCoolerStateCooling,
		}[output]
	default:
	}
	return nil
}

func characteristicFromType(typ string) *characteristic.C {
	switch typ {
	case "Active":
		return characteristic.NewActive().C
	case "On":
		return characteristic.NewOn().C
	case "CurrentTemperature":
		return characteristic.NewCurrentTemperature().C
	case "TargetTemperature":
		return characteristic.NewTargetTemperature().C
	case "TargetHeaterCoolerState":
		return characteristic.NewTargetHeaterCoolerState().C
	case "CurrentHeaterCoolerState":
		return characteristic.NewCurrentHeaterCoolerState().C
	case "RotationSpeed":
		return characteristic.NewRotationSpeed().C
	case "CoolingThresholdTemperature":
		return characteristic.NewCoolingThresholdTemperature().C
	case "HeatingThresholdTemperature":
		return characteristic.NewHeatingThresholdTemperature().C
	default:
		return nil
	}
}
