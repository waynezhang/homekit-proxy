package characteristics

import "strconv"

func parseFloat(v string) float64 {
	f, _ := strconv.ParseFloat(v, 64)
	return f
}

func truncateFloat(s string) string {
	f := parseFloat(s)
	return strconv.FormatInt(int64(f), 10)
}
