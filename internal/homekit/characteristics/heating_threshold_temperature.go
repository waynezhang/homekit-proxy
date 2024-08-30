package characteristics

import (
	"strconv"

	"github.com/brutella/hap/characteristic"
	g "github.com/maragudk/gomponents"
	"github.com/waynezhang/homekit-proxy/internal/config"
	"github.com/waynezhang/homekit-proxy/internal/homekit/stat"
	"github.com/waynezhang/homekit-proxy/internal/html"
)

func init() {
	const cType = "HeatingThresholdTemperature"
	registerCConstructor(cType, func(cc config.CharacteristicsConfig) *characteristic.C {
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
	registerConverterFromCommandLine(cType, func(v string) any {
		f, _ := strconv.ParseFloat(v, 64)
		return f
	})
	registerConverterToCommandLine(cType, func(v any) string {
		return strconv.FormatFloat(v.(float64), 'f', 2, 64)
	})
	registerHTMLElBuilderFunc(cType, func(name string, v string, id string, cst *stat.CharacteristicsStat) g.Node {
		return html.Slider(cst.Min, cst.Max, cst.Step, v, id)
	})
}
