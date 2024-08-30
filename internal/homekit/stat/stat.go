package stat

import "time"

type AutomationStat struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Cmd       string    `json:"cmd"`
	Cron      string    `json:"cron"`
	Tolerance int       `json:"tolerance"`
	LastRun   time.Time `json:"last_run"`
	LastError string    `json:"last_error"`
	NextRun   time.Time `json:"next_run"`
	Enabled   bool      `json:"enabled"`
}

type CharacteristicsStat struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
	Min   string `json:"min"`
	Max   string `json:"max"`
	Step  string `json:"step"`
}

type Stat struct {
	Now             time.Time              `json:"now"`
	Name            string                 `json:"name"`
	Characteristics []*CharacteristicsStat `json:"characteristics"`
	Automations     []*AutomationStat      `json:"automations"`
}
