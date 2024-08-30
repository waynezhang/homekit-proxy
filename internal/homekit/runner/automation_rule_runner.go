package runner

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	"github.com/bradhe/cadence"
	"github.com/waynezhang/homekit-proxy/internal/config"
	"github.com/waynezhang/homekit-proxy/internal/utils"
)

type AutomationRunner struct {
	Id        int
	Config    *config.AutomationConfig
	LastRun   time.Time
	LastError error
	NextRun   time.Time
}

func (r *AutomationRunner) Start(t time.Time, ctx context.Context) {
	now := time.Now()

	next, ref, err := nextRunTime(r.Config.Cron, r.Config.Tolerance, t)
	if err != nil {
		slog.Error("[Automation] Failed to start automation rule", "name", r.Config.Name, "err", err)
		return
	}

	r.NextRun = next

	slog.Info("[Automation] Scheduling next run time", "rule", r.Config.Name, "time", next)
	go func() {
		duration := next.Sub(now)
		time.AfterFunc(duration, func() {
			select {
			case <-ctx.Done():
				slog.Info("[Automation] Context is done, cancelling", "name", r.Config.Name, "cmd", r.Config.Cmd)
			default:
				if r.Config.Enabled {
					slog.Info("[Automation] Running automtion task", "name", r.Config.Name, "cmd", r.Config.Cmd)

					_, err := utils.Exec(r.Config.Cmd)
					r.LastRun = time.Now()
					r.LastError = err
				} else {
					slog.Info("[Automation] Skipping automtion task", "name", r.Config.Name, "cmd", r.Config.Cmd)
				}

				r.Start(ref, ctx)
			}
		})
	}()
}

func nextRunTime(cron string, tolerance int, ref time.Time) (time.Time, time.Time, error) {
	next, err := cadence.Next(cron, ref)
	if err != nil {
		return next, next, err
	}

	m := 0
	if tolerance != 0 {
		random := rand.New(rand.NewSource(ref.UnixNano()))
		m = random.Intn(tolerance*2) - tolerance
	}

	runTime := next.Add(time.Duration(m) * time.Second)

	return runTime, next, nil
}
