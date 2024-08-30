package utils

import (
	"fmt"
	"log/slog"
	"strconv"
)

func ParseFloat(v string) float64 {
	f, _ := strconv.ParseFloat(v, 64)
	return f
}

func TruncateFloat(s string) string {
	f := ParseFloat(s)
	return strconv.FormatInt(int64(f), 10)
}

func NumberToString(n interface{}) string {
	switch v := n.(type) {
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', 2, 64)
	case nil:
		return ""
	default:
		slog.Error("[API] Unhandled type", "type", v, "n", n)
		return fmt.Sprintf("%v", v)
	}
}

func ErrStringOrEmpty(e error) string {
	if e != nil {
		return e.Error()
	}

	return ""
}
