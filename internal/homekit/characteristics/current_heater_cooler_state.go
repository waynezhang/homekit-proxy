package characteristics

import (
	"github.com/brutella/hap/characteristic"
	g "github.com/maragudk/gomponents"
	"github.com/waynezhang/homekit-proxy/internal/config"
	"github.com/waynezhang/homekit-proxy/internal/homekit/stat"
)

func init() {
	const cType = "CurrentHeaterCoolerState"

	registerCConstructor(cType, func(cc config.CharacteristicsConfig) *characteristic.C {
		return characteristic.NewCurrentHeaterCoolerState().C
	})
	registerConverterFromCommandLine(cType, func(v string) any {
		return map[string]interface{}{
			"CurrentHeaterCoolerStateInactive": characteristic.CurrentHeaterCoolerStateInactive,
			"CurrentHeaterCoolerStateIdle":     characteristic.CurrentHeaterCoolerStateIdle,
			"CurrentHeaterCoolerStateHeating":  characteristic.CurrentHeaterCoolerStateHeating,
			"CurrentHeaterCoolerStateCooling":  characteristic.CurrentHeaterCoolerStateCooling,
		}[v]
	})
	registerConverterToCommandLine(cType, func(v any) string {
		return map[any]string{
			characteristic.CurrentHeaterCoolerStateInactive: "CurrentHeaterCoolerStateInactive",
			characteristic.CurrentHeaterCoolerStateIdle:     "CurrentHeaterCoolerStateIdle",
			characteristic.CurrentHeaterCoolerStateHeating:  "CurrentHeaterCoolerStateHeating",
			characteristic.CurrentHeaterCoolerStateCooling:  "CurrentHeaterCoolerStateCooling",
		}[v]
	})
	registerHTMLElBuilderFunc(cType, func(name string, v string, id string, cst *stat.CharacteristicsStat) g.Node {
		return radioGroup(
			name,
			[]string{
				"CurrentHeaterCoolerStateInactive",
				"CurrentHeaterCoolerStateIdle",
				"CurrentHeaterCoolerStateHeating",
				"CurrentHeaterCoolerStateCooling",
			},
			v,
			id,
		)
	})
}
