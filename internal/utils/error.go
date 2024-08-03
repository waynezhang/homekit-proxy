package utils

import (
	"fmt"
	"log/slog"
	"os"
)

func CheckFatalError(err error, msg string, args ...any) {
	if err == nil {
		return
	}
	str := fmt.Sprintf(msg, args...)
	slog.Error(str, "err", err)
	os.Exit(1)
}
