package characteristics

import (
	"github.com/brutella/hap/characteristic"
	"github.com/waynezhang/homekit-proxy/internal/config"
)

func init() {
	addParser("Active", func(v string) any {
		return map[string]interface{}{
			"ActiveInactive": characteristic.ActiveInactive,
			"ActiveActive":   characteristic.ActiveActive,
		}[v]
	})
	addConverter("Active", func(v any) string {
		return map[any]string{
			characteristic.ActiveInactive: "ActiveInactive",
			characteristic.ActiveActive:   "ActiveActive",
		}[v]
	})
	addCConstructor("Active", func(cc config.CharacteristicsConfig) *characteristic.C {
		return characteristic.NewActive().C
	})
}
