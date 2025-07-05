package homekit

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/waynezhang/homekit-proxy/internal/homekit/characteristics"
	"github.com/waynezhang/homekit-proxy/internal/homekit/stat"
	"github.com/waynezhang/homekit-proxy/internal/utils"
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
	m.server.ServeMux().HandleFunc("/s/action_logs", func(res http.ResponseWriter, req *http.Request) {
		entityID := req.URL.Query().Get("entity_id")
		characteristicType := req.URL.Query().Get("characteristic_type")

		var logs interface{}
		var err error

		if entityID != "" && characteristicType != "" {
			logs, err = m.actionLog.GetLogsByEntityAndCharacteristicType(entityID, characteristicType, 1000)
		} else if entityID != "" {
			logs, err = m.actionLog.GetLogsByEntity(entityID, 100)
		} else {
			logs, err = m.actionLog.GetRecentLogs(100)
		}

		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(err.Error()))
			return
		}
		j, _ := json.MarshalIndent(logs, "", "  ")
		res.Write(j)
	})
	m.server.ServeMux().HandleFunc("/s/action_logs/entities", func(res http.ResponseWriter, req *http.Request) {
		entities, err := m.actionLog.GetDistinctEntityIDs()
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(err.Error()))
			return
		}
		j, _ := json.MarshalIndent(entities, "", "  ")
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
			if id == strconv.Itoa(r.Id) {
				slog.Info("[API] Found runner", "id", r.Id)
				r.RunSetter(body.Value)
				break
			}
		}
		res.Write([]byte("{\"result\": \"OK\"}"))
	})

	m.server.ServeMux().HandleFunc("/s/a/{id}", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("Invalid method"))
			return
		}

		path := strings.TrimPrefix(req.URL.Path, "/s/a/")
		id, err := strconv.Atoi(strings.SplitN(path, "/", 2)[0])
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("Invalid id"))
			return
		}

		var body struct {
			Enabled string `json:"value"`
		}
		err = json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("Invalid request body"))
			return
		}

		enabled, err := strconv.ParseBool(body.Enabled)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("Invalid 'enabled' value"))
			return
		}

		slog.Info("[API] Update automation", "id", id, "enabled", body.Enabled)

		m.config.SetAutomationEnabled(id, enabled)
		res.Write([]byte("{\"result\": \"OK\"}"))
	})
    m.server.ServeMux().HandleFunc("/s/a/{id}/run", func(res http.ResponseWriter, req *http.Request) {
        if req.Method != "POST" {
            res.WriteHeader(http.StatusBadRequest)
            res.Write([]byte("Invalid method"))
            return
        }

        path := strings.TrimPrefix(req.URL.Path, "/s/a/")
        id, err := strconv.Atoi(strings.SplitN(path, "/", 2)[0])
        if err != nil {
            res.WriteHeader(http.StatusBadRequest)
            res.Write([]byte("Invalid id"))
            return
        }

        slog.Info("[API] Run automation", "id", id)

        err = m.runAutomation(id)
        if err != nil {
            res.WriteHeader(http.StatusInternalServerError)
            res.Write([]byte(err.Error()))
            return
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
			Id:    r.Id,
			Name:  r.Name,
			Type:  r.Config.Type,
			Value: characteristics.ConvertValueToCommandLine(r.GetLastValue(), r.Config.Type),
			Min:   utils.NumberToString(r.C.MinVal),
			Max:   utils.NumberToString(r.C.MaxVal),
			Step:  utils.NumberToString(r.C.StepVal),
			Area:  r.Area,
			Icon:  r.Icon,
		}
		csts = append(csts, &cst)
	}
	st.Characteristics = csts

	asts := []*stat.AutomationStat{}
	for _, a := range m.automations {
		ast := stat.AutomationStat{
			Id:        a.Config.Id,
			Name:      a.Config.Name,
			Cmd:       a.Config.Cmd,
			Cron:      a.Config.Cron,
			Tolerance: a.Config.Offset,
			LastRun:   a.LastRun,
			LastError: utils.ErrStringOrEmpty(a.LastError),
			NextRun:   time.Time{},
			Enabled:   a.Config.Enabled,
			Group:     a.Config.Group,
		}
		if a.Config.Enabled {
			ast.NextRun = a.NextRun
		}
		asts = append(asts, &ast)
	}
	st.Automations = asts
	return st
}
