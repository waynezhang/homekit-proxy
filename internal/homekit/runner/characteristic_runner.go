package runner

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/brutella/hap/characteristic"
	"github.com/waynezhang/homekit-proxy/internal/config"
	ch "github.com/waynezhang/homekit-proxy/internal/homekit/characteristics"
	"github.com/waynezhang/homekit-proxy/internal/utils"
)

type CharacteristicRunner struct {
	Id        int
	Name      string
	Config    *config.CharacteristicsConfig
	C         *characteristic.C
	LastValue any
}

func NewCharacteristicRunner(name string, config *config.CharacteristicsConfig, c *characteristic.C) *CharacteristicRunner {
	r := &CharacteristicRunner{
		Name:   name,
		Config: config,
		C:      c,
	}

	r.C.OnCValueUpdate(func(c *characteristic.C, new, old interface{}, req *http.Request) {
		slog.Info("[Characteristcs] Remote value changed", "name", r.Name, "new", new, "type", r.Config.Type, "hasReq", req != nil)
		if req == nil {
			return
		}

		param := ch.ConvertValueToCommandLine(new, r.Config.Type)
		r.RunSetter(param)
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
	output, _ := utils.Exec(cmd)

	val := ch.ParseValueFromCommandLine(output, r.Config.Type)
	if val != nil {
		if r.LastValue != val {
			slog.Info("[Characteristcs] Setting remote value", "name", r.Name, "val", val)
			r.LastValue = val
			r.C.SetValueRequest(val, nil)
		}
	} else {
		slog.Error("[Characteristcs] No value parsed")
	}
}

func (r *CharacteristicRunner) RunSetter(param string) {
	cmd := r.Config.Set + " " + param
	utils.Exec(cmd)
}

func (r *CharacteristicRunner) GetLastValue() any {
	if r.LastValue == nil {
		r.runGetter()
	}

	return r.LastValue
}
