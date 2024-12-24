package characteristics

import (
	"github.com/brutella/hap/characteristic"
	"github.com/waynezhang/homekit-proxy/internal/config"
)

func init() {
	const cType = "Active"

	registerCConstructor(cType, func(cc config.CharacteristicsConfig) *characteristic.C {
		return characteristic.NewActive().C
	})
	registerConverterFromCommandLine(cType, func(v string) any {
		return map[string]interface{}{
			"ActiveInactive": characteristic.ActiveInactive,
			"ActiveActive":   characteristic.ActiveActive,
		}[v]
	})
	registerConverterToCommandLine(cType, func(v any) string {
		return map[any]string{
			characteristic.ActiveInactive: "ActiveInactive",
			characteristic.ActiveActive:   "ActiveActive",
		}[v]
	})
}
