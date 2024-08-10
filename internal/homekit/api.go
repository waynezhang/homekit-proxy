package homekit

import (
	"encoding/json"
	"net/http"
	"time"
)

func (m *HMManager) startHealthCheckHandler() {
	m.server.ServeMux().HandleFunc("/health", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("OK"))
	})
}

func (m *HMManager) startAPIHandler() {
	handleGetAll(m)
}

func handleGetAll(m *HMManager) {
	type automationStat struct {
		Name      string    `json:"name"`
		Cmd       string    `json:"cmd"`
		Cron      string    `json:"cron"`
		Margin    int       `json:"margin"`
		LastRun   time.Time `json:"last_run"`
		LastError error     `json:"last_error"`
		NextRun   time.Time `json:"next_run"`
	}
	type characteristicsStat struct {
		Name  string `json:"name"`
		Value any    `json:"value"`
	}
	type stat struct {
		Now             time.Time              `json:"now"`
		Characteristics []*characteristicsStat `json:"characteristics"`
		Automations     []*automationStat      `json:"automations"`
	}

	m.server.ServeMux().HandleFunc("/s/all", func(res http.ResponseWriter, req *http.Request) {
		st := stat{
			Now: time.Now(),
		}

		csts := []*characteristicsStat{}
		for _, a := range m.root.runners {
			cst := characteristicsStat{
				Name:  a.name,
				Value: a.lastValue,
			}
			csts = append(csts, &cst)
		}
		st.Characteristics = csts

		asts := []*automationStat{}
		for _, a := range m.automations {
			ast := automationStat{
				Name:      a.config.Name,
				Cmd:       a.config.Cmd,
				Cron:      a.config.Cron,
				Margin:    a.config.Margin,
				LastRun:   a.lastRun,
				LastError: a.lastError,
				NextRun:   a.nextRun,
			}
			asts = append(asts, &ast)
		}
		st.Automations = asts

		j, _ := json.MarshalIndent(st, "", "  ")
		res.Write(j)
	})
}
