package characteristics

import (
	"strconv"

	"github.com/brutella/hap/characteristic"
	"github.com/waynezhang/homekit-proxy/internal/config"
)

func init() {
	addParser("HeatingThresholdTemperature", func(v string) any {
		f, _ := strconv.ParseFloat(v, 64)
		return f
	})
	addConverter("HeatingThresholdTemperature", func(v any) string {
		return strconv.FormatFloat(v.(float64), 'f', 2, 64)
	})
	addCConstructor("HeatingThresholdTemperature", func(cc config.CharacteristicsConfig) *characteristic.C {
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
	})
}
