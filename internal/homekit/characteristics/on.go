package characteristics

import (
	"github.com/brutella/hap/characteristic"
	"github.com/waynezhang/homekit-proxy/internal/config"
)

func init() {
	const cType = "On"

	registerCConstructor(cType, func(cc config.CharacteristicsConfig) *characteristic.C {
		return characteristic.NewOn().C
	})
	registerConverterFromCommandLine(cType, func(v string) any {
		return map[string]interface{}{
			"true":  true,
			"false": false,
		}[v]
	})
	registerConverterToCommandLine(cType, func(v any) string {
		return map[any]string{
			true:  "true",
			false: "false",
		}[v]
	})
}
