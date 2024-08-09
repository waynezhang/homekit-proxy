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
	name      string
	config    *config.CharacteristicsConfig
	c         *characteristic.C
	lastValue any
}

func newCharacteristicRunner(name string, config *config.CharacteristicsConfig, c *characteristic.C) *characteristicRunner {
	r := &characteristicRunner{
		name:   name,
		config: config,
		c:      c,
	}

	r.c.OnCValueUpdate(func(c *characteristic.C, new, old interface{}, req *http.Request) {
		slog.Info("[Characteristcs] Remote value changed", "name", r.name, "new", new, "type", r.config.Type, "hasReq", req != nil)
		if req == nil {
			return
		}

		param := valueStringFromCharacteristicType(new, r.config.Type)
		cmd := r.config.Set + " " + param
		utils.Exec(cmd)
	})

	return r
}

func (r *characteristicRunner) start() {
	cmd := r.config.Get
	if len(cmd) == 0 {
		slog.Info("[Characteristcs] No Getter, skip")
		return
	}

	go func() {
		for {
			slog.Info("[Characteristcs] Updating status of " + r.name)

			output, _ := utils.Exec(cmd)

			val := parseValueFromCharacteristicType(output, r.config.Type)
			if val != nil {
				if r.lastValue != val {
					slog.Info("[Characteristcs] Setting remote value", "name", r.name, "val", val)
					r.lastValue = val
					r.c.SetValueRequest(val, nil)
				}
			} else {
				slog.Error("[Characteristcs] No value parsed")
			}
			time.Sleep(time.Duration(r.config.Poll) * time.Second)
		}
	}()
}

func valueStringFromCharacteristicType(v any, typ string) string {
	switch typ {
	case "Active":
		return map[any]string{
			characteristic.ActiveInactive: "ActiveInactive",
			characteristic.ActiveActive:   "ActiveActive",
		}[v]
	case "On":
		return map[any]string{
			true:  "true",
			false: "false",
		}[v]
	case "CurrentTemperature", "RotationSpeed", "CoolingThresholdTemperature", "HeatingThresholdTemperature":
		return strconv.FormatFloat(v.(float64), 'f', 2, 64)
	case "TargetHeaterCoolerState":
		return map[any]string{
			characteristic.TargetHeaterCoolerStateAuto: "TargetHeaterCoolerStateAuto",
			characteristic.TargetHeaterCoolerStateHeat: "TargetHeaterCoolerStateHeat",
			characteristic.TargetHeaterCoolerStateCool: "TargetHeaterCoolerStateCool",
		}[v]
	case "CurrentHeaterCoolerState":
		return map[any]string{
			characteristic.CurrentHeaterCoolerStateInactive: "CurrentHeaterCoolerStateInactive",
			characteristic.CurrentHeaterCoolerStateIdle:     "CurrentHeaterCoolerStateIdle",
			characteristic.CurrentHeaterCoolerStateHeating:  "CurrentHeaterCoolerStateHeating",
			characteristic.CurrentHeaterCoolerStateCooling:  "CurrentHeaterCoolerStateCooling",
		}[v]
	default:
	}
	return ""
}

func parseValueFromCharacteristicType(v string, typ string) any {
	switch typ {
	case "Active":
		return map[string]interface{}{
			"ActiveInactive": characteristic.ActiveInactive,
			"ActiveActive":   characteristic.ActiveActive,
		}[v]
	case "On":
		return map[string]interface{}{
			"true":  true,
			"false": false,
		}[v]
	case "CurrentTemperature", "RotationSpeed", "CoolingThresholdTemperature", "HeatingThresholdTemperature":
		f, _ := strconv.ParseFloat(v, 64)
		return f
	case "TargetHeaterCoolerState":
		return map[string]interface{}{
			"TargetHeaterCoolerStateAuto": characteristic.TargetHeaterCoolerStateAuto,
			"TargetHeaterCoolerStateHeat": characteristic.TargetHeaterCoolerStateHeat,
			"TargetHeaterCoolerStateCool": characteristic.TargetHeaterCoolerStateCool,
		}[v]
	case "CurrentHeaterCoolerState":
		return map[string]interface{}{
			"CurrentHeaterCoolerStateInactive": characteristic.CurrentHeaterCoolerStateInactive,
			"CurrentHeaterCoolerStateIdle":     characteristic.CurrentHeaterCoolerStateIdle,
			"CurrentHeaterCoolerStateHeating":  characteristic.CurrentHeaterCoolerStateHeating,
			"CurrentHeaterCoolerStateCooling":  characteristic.CurrentHeaterCoolerStateCooling,
		}[v]
	default:
	}
	return nil
}

func characteristicFromConfig(cc config.CharacteristicsConfig) *characteristic.C {
	switch cc.Type {
	case "Active":
		return characteristic.NewActive().C
	case "On":
		return characteristic.NewOn().C
	case "CurrentTemperature":
		return characteristic.NewCurrentTemperature().C
	case "TargetHeaterCoolerState":
		return characteristic.NewTargetHeaterCoolerState().C
	case "CurrentHeaterCoolerState":
		return characteristic.NewCurrentHeaterCoolerState().C
	case "RotationSpeed":
		c := characteristic.NewRotationSpeed()
		if cc.Min > 0 {
			c.SetMinValue(float64(cc.Min))
		}
		if cc.Max > 0 {
			c.SetMaxValue(float64(cc.Max))
		}
		if cc.Step > 0 {
			c.SetStepValue(float64(cc.Step))
		}
		return c.C
	case "CoolingThresholdTemperature":
		c := characteristic.NewCoolingThresholdTemperature()
		if cc.Min > 0 {
			c.SetMinValue(float64(cc.Min))
		}
		if cc.Max > 0 {
			c.SetMaxValue(float64(cc.Max))
		}
		if cc.Step > 0 {
			c.SetStepValue(float64(cc.Step))
		}
		return c.C
	case "HeatingThresholdTemperature":
		c := characteristic.NewHeatingThresholdTemperature()
		if cc.Min > 0 {
			c.SetMinValue(float64(cc.Min))
		}
		if cc.Max > 0 {
			c.SetMaxValue(float64(cc.Max))
		}
		if cc.Step > 0 {
			c.SetStepValue(float64(cc.Step))
		}
		return c.C
	default:
		return nil
	}
}
