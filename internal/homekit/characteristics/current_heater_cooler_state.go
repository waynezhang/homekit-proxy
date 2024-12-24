package characteristics

import (
	"github.com/brutella/hap/characteristic"
	"github.com/waynezhang/homekit-proxy/internal/config"
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
}
