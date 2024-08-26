package characteristics

import (
	"strconv"

	"github.com/brutella/hap/characteristic"
	g "github.com/maragudk/gomponents"
	"github.com/waynezhang/homekit-proxy/internal/config"
	"github.com/waynezhang/homekit-proxy/internal/homekit/stat"
)

func init() {
	const cType = "RotationSpeed"

	registerCConstructor(cType, func(cc config.CharacteristicsConfig) *characteristic.C {
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
	})
	registerConverterFromCommandLine(cType, func(v string) any {
		return parseFloat(v)
	})
	registerConverterToCommandLine(cType, func(v any) string {
		return strconv.FormatFloat(v.(float64), 'f', 2, 64)
	})
	registerHTMLElBuilderFunc(cType, func(name string, v string, id string, cst *stat.CharacteristicsStat) g.Node {
		return slider(cst.Min, cst.Max, cst.Step, v, id)
	})
}
