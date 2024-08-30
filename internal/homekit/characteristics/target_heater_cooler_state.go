package characteristics

import (
	"github.com/brutella/hap/characteristic"
	g "github.com/maragudk/gomponents"

	"github.com/waynezhang/homekit-proxy/internal/config"
	"github.com/waynezhang/homekit-proxy/internal/homekit/stat"
	"github.com/waynezhang/homekit-proxy/internal/html"
)

func init() {
	const cType = "TargetHeaterCoolerState"

	registerCConstructor(cType, func(cc config.CharacteristicsConfig) *characteristic.C {
		return characteristic.NewTargetHeaterCoolerState().C
	})
	registerConverterFromCommandLine(cType, func(v string) any {
		return map[string]interface{}{
			"TargetHeaterCoolerStateAuto": characteristic.TargetHeaterCoolerStateAuto,
			"TargetHeaterCoolerStateHeat": characteristic.TargetHeaterCoolerStateHeat,
			"TargetHeaterCoolerStateCool": characteristic.TargetHeaterCoolerStateCool,
		}[v]
	})
	registerConverterToCommandLine(cType, func(v any) string {
		return map[any]string{
			characteristic.TargetHeaterCoolerStateAuto: "TargetHeaterCoolerStateAuto",
			characteristic.TargetHeaterCoolerStateHeat: "TargetHeaterCoolerStateHeat",
			characteristic.TargetHeaterCoolerStateCool: "TargetHeaterCoolerStateCool",
		}[v]
	})
	registerHTMLElBuilderFunc(cType, func(name string, v string, id string, cst *stat.CharacteristicsStat) g.Node {
		return html.RadioGroup(
			name,
			[]string{"TargetHeaterCoolerStateAuto", "TargetHeaterCoolerStateHeat", "TargetHeaterCoolerStateCool"},
			v,
			id,
			ExtraTypeCharacteristic,
		)
	})
}
