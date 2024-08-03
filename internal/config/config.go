package config

import (
	"github.com/spf13/viper"
	"github.com/waynezhang/homekit-proxy/internal/utils"
)

type Config struct {
	Bridge      BridgeConfig
	Accessories []AccessoriesConfig
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
}

func Parse(file string) Config {
	config := Config{}

	v := viper.New()
	v.SetConfigFile(file)

	err := v.ReadInConfig()
	utils.CheckFatalError(err, "Failed to parse config file %s", file)

	err = v.Unmarshal(&config)
	utils.CheckFatalError(err, "Failed to parse config file %s", file)

	return config
}
