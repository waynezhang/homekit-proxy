package homekit

import (
	"log/slog"
	"math/rand"
	"time"

	"github.com/bradhe/cadence"
	"github.com/waynezhang/homekit-proxy/internal/config"
	"github.com/waynezhang/homekit-proxy/internal/utils"
)

type automationRunner struct {
	config    *config.AutomationConfig
	lastRun   time.Time
	lastError error
	nextRun   time.Time
}

func (r *automationRunner) start(t time.Time) {
	now := time.Now()

	next, ref, err := nextRunTime(r.config.Cron, r.config.Margin, t)
	if err != nil {
		slog.Error("[Automation] Failed to start automation rule", "name", r.config.Name, "err", err)
		return
	}

	r.nextRun = next

	slog.Info("[Automation] Scheduling next run time", "rule", r.config.Name, "time", next)
	go func() {
		duration := next.Sub(now)
		time.AfterFunc(duration, func() {
			_, err := utils.Exec(r.config.Cmd)

			slog.Info("[Automation] Running automtion task", "name", r.config.Name, "cmd", r.config.Cmd)

			r.lastRun = time.Now()
			r.lastError = err

			r.start(ref)
		})
	}()
}

func nextRunTime(cron string, margin int, ref time.Time) (time.Time, time.Time, error) {
	next, err := cadence.Next(cron, ref)
	if err != nil {
		return next, next, err
	}

	m := 0
	if margin != 0 {
		random := rand.New(rand.NewSource(ref.UnixNano()))
		m = random.Intn(margin*2) - margin
	}

	runTime := next.Add(time.Duration(m) * time.Second)

	return runTime, next, nil
}
