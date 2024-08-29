package config

import (
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/waynezhang/homekit-proxy/internal/utils"
)

type Config struct {
	Bridge      BridgeConfig
	Accessories []*AccessoriesConfig
	Automations []*AutomationConfig
	kv          *kv
}

type BridgeConfig struct {
	Name         string
	Manufacturer string
	Model        string
	Firmware     string
	PinCode      string
}

type AccessoriesConfig struct {
	Id           int
	Name         string
	Manufacturer string
	Model        string
	Firmware     string
	TypeByte     int
	Services     []ServicesConfig
}

type ServicesConfig struct {
	TypeString      string
	Characteristics []CharacteristicsConfig
}

type CharacteristicsConfig struct {
	Type string
	Poll int
	Set  string
	Get  string
	Min  int
	Max  int
	Step int
}

type AutomationConfig struct {
	Name      string
	Cron      string
	Cmd       string
	Tolerance int
	Enabled   bool
	Id        int
}

func Parse(file string, directory string) Config {
	config := Config{}

	v := viper.New()
	v.SetConfigFile(file)

	err := v.ReadInConfig()
	utils.CheckFatalError(err, "Failed to parse config file %s", file)

	err = v.Unmarshal(&config)
	utils.CheckFatalError(err, "Failed to parse config file %s", file)

	config.kv = newKV(filepath.Join(directory, "automation-config.json"))
	for _, a := range config.Automations {
		a.Enabled = config.kv.getBool(a.Id, true)
	}

	return config
}
