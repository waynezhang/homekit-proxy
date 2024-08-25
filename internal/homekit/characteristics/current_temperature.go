package characteristics

import (
	"strconv"

	"github.com/brutella/hap/characteristic"
	"github.com/waynezhang/homekit-proxy/internal/config"
)

func init() {
	addParser("CurrentTemperature", func(v string) any {
		f, _ := strconv.ParseFloat(v, 64)
		return f
	})
	addConverter("CurrentTemperature", func(v any) string {
		return strconv.FormatFloat(v.(float64), 'f', 2, 64)
	})
	addCConstructor("CurrentTemperature", func(cc config.CharacteristicsConfig) *characteristic.C {
		return characteristic.NewCurrentTemperature().C
	})
}
