package characteristics

import (
	"github.com/brutella/hap/characteristic"
	g "github.com/maragudk/gomponents"
	"github.com/waynezhang/homekit-proxy/internal/config"
	"github.com/waynezhang/homekit-proxy/internal/homekit/stat"
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
	registerHTMLElBuilderFunc(cType, func(name string, v string, id string, cst *stat.CharacteristicsStat) g.Node {
		return radioGroup(name, []string{"ActiveActive", "ActiveInactive"}, v, id)
	})
}
