package utils

import "strconv"

func ParseFloat(v string) float64 {
	f, _ := strconv.ParseFloat(v, 64)
	return f
}

func TruncateFloat(s string) string {
	f := ParseFloat(s)
	return strconv.FormatInt(int64(f), 10)
}
