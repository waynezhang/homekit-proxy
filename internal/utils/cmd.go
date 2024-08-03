package utils

import (
	"log/slog"
	"os/exec"
	"strings"
)

func Exec(cmd string) string {
	output, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		slog.Error("Error to exec", "cmd", cmd, "err", err)
		return ""
	}
	s := strings.TrimSpace(string(output))
	slog.Info("Exec", "cmd", cmd, "out", s)
	return s
}
