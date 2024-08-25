package homekit

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/brutella/hap/characteristic"
	"github.com/waynezhang/homekit-proxy/internal/config"
	ch "github.com/waynezhang/homekit-proxy/internal/homekit/characteristics"
	"github.com/waynezhang/homekit-proxy/internal/utils"
)

type characteristicRunner struct {
	id        int
	name      string
	config    *config.CharacteristicsConfig
	c         *characteristic.C
	lastValue any
}

func newCharacteristicRunner(name string, config *config.CharacteristicsConfig, c *characteristic.C) *characteristicRunner {
	r := &characteristicRunner{
		name:   name,
		config: config,
		c:      c,
	}

	r.c.OnCValueUpdate(func(c *characteristic.C, new, old interface{}, req *http.Request) {
		slog.Info("[Characteristcs] Remote value changed", "name", r.name, "new", new, "type", r.config.Type, "hasReq", req != nil)
		if req == nil {
			return
		}

		param := ch.ToString(new, r.config.Type)
		r.runSetter(param)
	})

	return r
}

func (r *characteristicRunner) start() {
	if len(r.config.Get) == 0 {
		slog.Info("[Characteristcs] No Getter, skip")
		return
	}

	go func() {
		for {
			r.runGetter()
			time.Sleep(time.Duration(r.config.Poll) * time.Second)
		}
	}()
}

func (r *characteristicRunner) runGetter() {
	slog.Info("[Characteristcs] Updating status of " + r.name)

	cmd := r.config.Get
	output, _ := utils.Exec(cmd)

	val := ch.ParseValue(output, r.config.Type)
	if val != nil {
		if r.lastValue != val {
			slog.Info("[Characteristcs] Setting remote value", "name", r.name, "val", val)
			r.lastValue = val
			r.c.SetValueRequest(val, nil)
		}
	} else {
		slog.Error("[Characteristcs] No value parsed")
	}
}

func (r *characteristicRunner) runSetter(param string) {
	cmd := r.config.Set + " " + param
	utils.Exec(cmd)
}

func (r *characteristicRunner) getLastValue() any {
	if r.lastValue == nil {
		r.runGetter()
	}

	return r.lastValue
}

func characteristicFromConfig(cc config.CharacteristicsConfig) *characteristic.C {
	return ch.NewCharacteristic(cc)
}
