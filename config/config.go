package config

import (
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

type NatsConfig struct {
	MQHost string `mapstructure:"MQ_HOST"`
	MQPort string `mapstructure:"MQ_PORT"`
}

type AppConfig struct {
	APPort string `mapstructure:"AP_PORT"`
}

// LoadNatsConfig
func LoadNatsConfig() (NatsConfig, error) {
	viper.AutomaticEnv()
	viper.BindEnv("MQ_HOST")
	viper.BindEnv("MQ_PORT")
	var config NatsConfig
	if err := viper.Unmarshal(&config); err != nil {
		return NatsConfig{}, err
	}
	log.Infof("[NATS] host: %v, port: %v", config.MQHost, config.MQPort)
	return config, nil
}

// LoadAppConfig
func LoadAppConfig() (AppConfig, error) {
	viper.AutomaticEnv()
	viper.BindEnv("AP_PORT")
	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		return AppConfig{}, err
	}
	log.Infof("[App] port: %v", config.APPort)
	return config, nil
}
