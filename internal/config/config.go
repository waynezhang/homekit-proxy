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
	Name    string
	Cron    string
	Cmd     string
	Offset  int
	Enabled bool
	Id      int
}

func Parse(configDir string, directory string) Config {
	config := Config{}

	// Parse device.yaml
	deviceFile := filepath.Join(configDir, "device.yaml")
	vDevice := viper.New()
	vDevice.SetConfigFile(deviceFile)

	err := vDevice.ReadInConfig()
	utils.CheckFatalError(err, "Failed to parse device config file %s", deviceFile)

	var deviceConfig struct {
		Bridge      BridgeConfig
		Accessories []*AccessoriesConfig
	}
	err = vDevice.Unmarshal(&deviceConfig)
	utils.CheckFatalError(err, "Failed to parse device config file %s", deviceFile)

	config.Bridge = deviceConfig.Bridge
	config.Accessories = deviceConfig.Accessories

	// Parse automation.yaml
	automationFile := filepath.Join(configDir, "automation.yaml")
	vAutomation := viper.New()
	vAutomation.SetConfigFile(automationFile)

	err = vAutomation.ReadInConfig()
	utils.CheckFatalError(err, "Failed to parse automation config file %s", automationFile)

	var automationConfig struct {
		Automations []*AutomationConfig
	}
	err = vAutomation.Unmarshal(&automationConfig)
	utils.CheckFatalError(err, "Failed to parse automation config file %s", automationFile)

	config.Automations = automationConfig.Automations

	config.kv = newKV(filepath.Join(directory, "automation-config.json"))
	for _, a := range config.Automations {
		a.Enabled = config.kv.getBool(a.Id, true)
	}

	return config
}

func (cfg *Config) SetAutomationEnabled(id int, enabled bool) {
	for _, a := range cfg.Automations {
		if a.Id == id {
			a.Enabled = enabled
			cfg.kv.setBool(id, enabled)
			break
		}
	}
}
