package runner

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/bradhe/cadence"
	"github.com/waynezhang/homekit-proxy/internal/actionlog"
	"github.com/waynezhang/homekit-proxy/internal/config"
	"github.com/waynezhang/homekit-proxy/internal/utils"
)

type AutomationRunner struct {
	Id        int
	Config    *config.AutomationConfig
	LastRun   time.Time
	LastError error
	NextRun   time.Time
	ActionLog *actionlog.ActionLog
	NtfyTopic string
}

func (r *AutomationRunner) Start(t time.Time, ctx context.Context) {
	now := time.Now()

	next, ref, err := nextRunTime(r.Config.Cron, r.Config.Offset, t)
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
					r.ActionLog.LogAction(fmt.Sprintf("Automation %d", r.Config.Id), "automation_run", "success")
					r.sendNtfyNotification()
				} else {
					slog.Info("[Automation] Skipping automtion task", "name", r.Config.Name, "cmd", r.Config.Cmd)
				}

				r.Start(ref, ctx)
			}
		})
	}()
}

func (r *AutomationRunner) sendNtfyNotification() {
	if r.NtfyTopic == "" {
		return
	}

	slog.Info("[Automation] Sending ntfy notification", "name", r.Config.Name, "topic", r.NtfyTopic)
	_, err := http.Post(
		fmt.Sprintf("https://ntfy.sh/%s", r.NtfyTopic),
		"text/plain",
		strings.NewReader(fmt.Sprintf("Automation %s executed", r.Config.Name)),
	)
	if err != nil {
		slog.Error("[Automation] Failed to send ntfy notification", "name", r.Config.Name, "err", err)
	}
}

func nextRunTime(cron string, offset int, ref time.Time) (time.Time, time.Time, error) {
	next, err := cadence.Next(cron, ref)
	if err != nil {
		return next, next, err
	}

	m := 0
	if offset != 0 {
		random := rand.New(rand.NewSource(ref.UnixNano()))
		m = random.Intn(offset*2) - offset
	}

	runTime := next.Add(time.Duration(m) * time.Second)

	return runTime, next, nil
}
