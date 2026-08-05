package runner

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/brutella/hap/characteristic"
	"github.com/waynezhang/homekit-proxy/internal/actionlog"
	"github.com/waynezhang/homekit-proxy/internal/config"
	ch "github.com/waynezhang/homekit-proxy/internal/homekit/characteristics"
	"github.com/waynezhang/homekit-proxy/internal/utils"
)

type CharacteristicRunner struct {
	Id        int
	Name      string
	Area      string
	Icon      string
	Config    *config.CharacteristicsConfig
	C         *characteristic.C
	LastValue any
	ActionLog *actionlog.ActionLog
}

func NewCharacteristicRunner(name string, area string, icon string, config *config.CharacteristicsConfig, c *characteristic.C, actionLog *actionlog.ActionLog) *CharacteristicRunner {
	r := &CharacteristicRunner{
		Name:      name,
		Area:      area,
		Icon:      icon,
		Config:    config,
		C:         c,
		ActionLog: actionLog,
	}

	r.C.OnCValueUpdate(func(c *characteristic.C, new, old interface{}, req *http.Request) {
		slog.Info("[Characteristcs] Remote value changed", "name", r.Name, "new", new, "type", r.Config.Type, "hasReq", req != nil)
		if req == nil {
			return
		}

		param := ch.ConvertValueToCommandLine(new, r.Config.Type)
		if err := r.RunSetter(param); err != nil {
			slog.Error("[Characteristcs] Setter failed, not logging action", "name", r.Name, "err", err)
			return
		}

		r.ActionLog.LogAction(r.Name, r.Config.Type, param)
	})

	return r
}

func (r *CharacteristicRunner) Start(ctx context.Context) {
	if len(r.Config.Get) == 0 {
		slog.Info("[Characteristcs] No Getter, skip")
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				slog.Info("[Characteristcs] Context is done, cancelling", "name", r.Name)
			default:
				r.runGetter()
			}
			time.Sleep(time.Duration(r.Config.Poll) * time.Second)
		}
	}()
}

func (r *CharacteristicRunner) runGetter() {
	slog.Info("[Characteristcs] Updating status of " + r.Name)

	cmd := r.Config.Get
	output, err := utils.Exec(cmd)
	if err != nil {
		slog.Error("[Characteristcs] Getter failed, skipping update", "name", r.Name, "err", err)
		return
	}

	val := ch.ParseValueFromCommandLine(output, r.Config.Type)
	if val != nil {
		if r.LastValue != val {
			slog.Info("[Characteristcs] Setting remote value", "name", r.Name, "val", val)
			r.LastValue = val
			r.C.SetValueRequest(val, nil)
		}
		valStr := ch.ConvertValueToCommandLine(val, r.Config.Type)
		r.ActionLog.LogAction(r.Name, r.Config.Type, valStr)
	} else {
		slog.Error("[Characteristcs] No value parsed")
	}
}

func (r *CharacteristicRunner) RunSetter(param string) error {
	cmd := r.Config.Set + " " + param
	_, err := utils.Exec(cmd)
	return err
}

func (r *CharacteristicRunner) GetLastValue() any {
	if r.LastValue == nil {
		r.runGetter()
	}

	return r.LastValue
}
