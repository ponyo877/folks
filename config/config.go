package config

import (
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

type NatsConfig struct {
	MQPassword string `mapstructure:"MQ_PASSWORD"`
	MQHost     string `mapstructure:"MQ_HOST"`
	MQPort     string `mapstructure:"MQ_PORT"`
}

type AppConfig struct {
	APRoot string `mapstructure:"AP_ROOT"`
	APPort string `mapstructure:"AP_PORT"`
}

// LoadNatsConfig
func LoadNatsConfig() (NatsConfig, error) {
	viper.AutomaticEnv()
	viper.BindEnv("MQ_PASSWORD")
	viper.BindEnv("MQ_HOST")
	viper.BindEnv("MQ_PORT")
	var config NatsConfig
	if err := viper.Unmarshal(&config); err != nil {
		return NatsConfig{}, err
	}
	log.Infof("[NATS] pass: %v, host: %v, port: %v", config.MQPassword, config.MQHost, config.MQPort)
	return config, nil
}

// LoadAppConfig
func LoadAppConfig() (AppConfig, error) {
	viper.AutomaticEnv()
	viper.BindEnv("AP_ROOT")
	viper.BindEnv("AP_PORT")
	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		return AppConfig{}, err
	}
	log.Infof("[App] root: %v, port: %v", config.APRoot, config.APPort)
	return config, nil
}
