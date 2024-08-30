package characteristics

import (
	"log/slog"

	"github.com/brutella/hap/characteristic"
	g "github.com/maragudk/gomponents"
	"github.com/waynezhang/homekit-proxy/internal/config"
	"github.com/waynezhang/homekit-proxy/internal/homekit/stat"
)

type parseFunc func(v string) any
type convertFunc func(v any) string
type cConstructorFunc func(config.CharacteristicsConfig) *characteristic.C
type htmlElFunc func(name string, val string, id string, cst *stat.CharacteristicsStat) g.Node

var parserMap = map[string]parseFunc{}
var converterMap = map[string]convertFunc{}
var cConstructorMap = map[string]cConstructorFunc{}
var htmlElFuncMap = map[string]htmlElFunc{}

const (
	ExtraTypeCharacteristic string = "C"
	ExtraTypeAutomation     string = "A"
)

func NewCharacteristic(cc config.CharacteristicsConfig) *characteristic.C {
	fn := cConstructorMap[cc.Type]
	if fn == nil {
		slog.Error("[C Constructor] No C Constructor found", "type", cc.Type)
		return nil
	}

	return fn(cc)
}

func ParseValueFromCommandLine(v string, typ string) any {
	fn := parserMap[typ]
	if fn == nil {
		slog.Error("[Parser] No parser found", "type", typ)
		return nil
	}

	return fn(v)
}

func ConvertValueToCommandLine(v any, typ string) string {
	fn := converterMap[typ]
	if fn == nil {
		slog.Error("[Converter] No converter found", "type", typ)
		return ""
	}

	return fn(v)
}

func BuildHtmlEl(name string, v string, id string, cst *stat.CharacteristicsStat) g.Node {
	fn := htmlElFuncMap[cst.Type]
	if fn == nil {
		slog.Error("[HTML] No HTML El func found", "type", cst.Type)
		return g.Text(v)
	}

	return fn(name, v, id, cst)
}

func registerCConstructor(typ string, fn cConstructorFunc) {
	cConstructorMap[typ] = fn
}

func registerConverterFromCommandLine(typ string, fn parseFunc) {
	parserMap[typ] = fn
}

func registerConverterToCommandLine(typ string, fn convertFunc) {
	converterMap[typ] = fn
}

func registerHTMLElBuilderFunc(typ string, fn htmlElFunc) {
	htmlElFuncMap[typ] = fn
}
