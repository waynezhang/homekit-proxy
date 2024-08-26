package homekit

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/waynezhang/homekit-proxy/internal/homekit/characteristics"
	"github.com/waynezhang/homekit-proxy/internal/homekit/stat"
)

func (m *HMManager) startHealthCheckHandler() {
	m.server.ServeMux().HandleFunc("/health", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("OK"))
	})
}

func (m *HMManager) startAPIHandler() {
	handleGetAll(m)
	handleUpdate(m)
}

func handleGetAll(m *HMManager) {
	m.server.ServeMux().HandleFunc("/s/all", func(res http.ResponseWriter, req *http.Request) {
		st := m.getAllStat()
		j, _ := json.MarshalIndent(st, "", "  ")
		res.Write(j)
	})
}

func handleUpdate(m *HMManager) {
	m.server.ServeMux().HandleFunc("/s/c/{id}", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("Invalid method"))
			return
		}

		path := strings.TrimPrefix(req.URL.Path, "/s/c/")
		id := strings.SplitN(path, "/", 2)[0]

		var body struct {
			Value string `json:"value"`
		}
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("Invalid request body"))
			return
		}
		slog.Info("[API] Update characteristics", "id", id, "value", body.Value)

		for _, r := range m.root.runners {
			if id == strconv.Itoa(r.id) {
				slog.Info("[API] Found runner", "id", r.id)
				r.runSetter(body.Value)
				break
			}
		}
		res.Write([]byte("{\"result\": \"OK\"}"))
	})
}

func (m *HMManager) getAllStat() stat.Stat {
	st := stat.Stat{
		Now:  time.Now(),
		Name: m.root.b.Name(),
	}

	csts := []*stat.CharacteristicsStat{}
	for _, r := range m.root.runners {
		cst := stat.CharacteristicsStat{
			Id:    r.id,
			Name:  r.name,
			Type:  r.config.Type,
			Value: characteristics.ConvertValueToCommandLine(r.getLastValue(), r.config.Type),
			Min:   numToString(r.c.MinVal),
			Max:   numToString(r.c.MaxVal),
			Step:  numToString(r.c.StepVal),
		}
		csts = append(csts, &cst)
	}
	st.Characteristics = csts

	asts := []*stat.AutomationStat{}
	for _, a := range m.automations {
		ast := stat.AutomationStat{
			Id:        a.id,
			Name:      a.config.Name,
			Cmd:       a.config.Cmd,
			Cron:      a.config.Cron,
			Tolerance: a.config.Tolerance,
			LastRun:   a.lastRun,
			LastError: errToString(a.lastError),
			NextRun:   a.nextRun,
		}
		asts = append(asts, &ast)
	}
	st.Automations = asts
	return st
}

func errToString(e error) string {
	if e != nil {
		return e.Error()
	}

	return ""
}

func numToString(n interface{}) string {
	switch v := n.(type) {
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', 2, 64)
	case nil:
		return ""
	default:
		slog.Error("[API] Unhandled type", "type", v, "n", n)
		return fmt.Sprintf("%v", v)
	}
}
