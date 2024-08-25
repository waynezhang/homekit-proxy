package characteristics

import (
	"log/slog"

	"github.com/brutella/hap/characteristic"
	"github.com/waynezhang/homekit-proxy/internal/config"
)

type parseFunc func(v string) any
type convertFunc func(v any) string
type cConstructorFunc func(config.CharacteristicsConfig) *characteristic.C

var parserMap = map[string]parseFunc{}
var converterMap = map[string]convertFunc{}
var cConstructorMap = map[string]cConstructorFunc{}

func addParser(typ string, fn parseFunc) {
	parserMap[typ] = fn
}

func addConverter(typ string, fn convertFunc) {
	converterMap[typ] = fn
}

func addCConstructor(typ string, fn cConstructorFunc) {
	cConstructorMap[typ] = fn
}

func ParseValue(v string, typ string) any {
	fn := parserMap[typ]
	if fn == nil {
		slog.Error("[Parser] No parser found", "type", typ)
		return nil
	}

	return fn(v)
}

func ToString(v any, typ string) string {
	fn := converterMap[typ]
	if fn == nil {
		slog.Error("[Converter] No converter found", "type", typ)
		return ""
	}

	return fn(v)
}

func NewCharacteristic(cc config.CharacteristicsConfig) *characteristic.C {
	fn := cConstructorMap[cc.Type]
	if fn == nil {
		slog.Error("[C Constructor] No C Constructor found", "type", cc.Type)
		return nil
	}

	return fn(cc)
}
