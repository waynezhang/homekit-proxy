package characteristics

import (
	"github.com/brutella/hap/characteristic"
	"github.com/waynezhang/homekit-proxy/internal/config"
)

func init() {
	addParser("On", func(v string) any {
		return map[string]interface{}{
			"true":  true,
			"false": false,
		}[v]
	})
	addConverter("On", func(v any) string {
		return map[any]string{
			true:  "true",
			false: "false",
		}[v]
	})
	addCConstructor("On", func(cc config.CharacteristicsConfig) *characteristic.C {
		return characteristic.NewOn().C
	})
}
