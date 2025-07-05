package characteristics

import (
	"strconv"

	"github.com/brutella/hap/characteristic"
	"github.com/waynezhang/homekit-proxy/internal/config"
)

func init() {
	const cType = "CurrentRelativeHumidity"

	registerCConstructor(cType, func(cc config.CharacteristicsConfig) *characteristic.C {
		return characteristic.NewCurrentRelativeHumidity().C
	})
	registerConverterFromCommandLine(cType, func(v string) any {
		f, _ := strconv.ParseFloat(v, 64)
		return f
	})
	registerConverterToCommandLine(cType, func(v any) string {
		return strconv.FormatFloat(v.(float64), 'f', 2, 64)
	})
}