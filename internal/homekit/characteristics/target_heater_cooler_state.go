package characteristics

import (
	"github.com/brutella/hap/characteristic"

	"github.com/waynezhang/homekit-proxy/internal/config"
)

func init() {
	addParser("TargetHeaterCoolerState", func(v string) any {
		return map[string]interface{}{
			"TargetHeaterCoolerStateAuto": characteristic.TargetHeaterCoolerStateAuto,
			"TargetHeaterCoolerStateHeat": characteristic.TargetHeaterCoolerStateHeat,
			"TargetHeaterCoolerStateCool": characteristic.TargetHeaterCoolerStateCool,
		}[v]
	})
	addConverter("TargetHeaterCoolerState", func(v any) string {
		return map[any]string{
			characteristic.TargetHeaterCoolerStateAuto: "TargetHeaterCoolerStateAuto",
			characteristic.TargetHeaterCoolerStateHeat: "TargetHeaterCoolerStateHeat",
			characteristic.TargetHeaterCoolerStateCool: "TargetHeaterCoolerStateCool",
		}[v]
	})
	addCConstructor("TargetHeaterCoolerState", func(cc config.CharacteristicsConfig) *characteristic.C {
		return characteristic.NewTargetHeaterCoolerState().C
	})
}
