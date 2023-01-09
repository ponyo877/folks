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

type RedisConfig struct {
	KVSPassword string `mapstructure:"KVS_PASSWORD"`
	KVSHost     string `mapstructure:"KVS_HOST"`
	KVSDatabase int    `mapstructure:"KVS_DATABASE"`
	KVSPort     string `mapstructure:"KVS_PORT"`
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

// LoadRedisConfig
func LoadRedisConfig() (RedisConfig, error) {
	viper.AutomaticEnv()
	viper.BindEnv("KVS_PASSWORD")
	viper.BindEnv("KVS_HOST")
	viper.BindEnv("KVS_DATABASE")
	viper.BindEnv("KVS_PORT")
	var config RedisConfig
	if err := viper.Unmarshal(&config); err != nil {
		return RedisConfig{}, err
	}
	log.Infof("[Redis] pass: %v, host: %v, db: %v, port: %v", config.KVSPassword, config.KVSHost, config.KVSDatabase, config.KVSPort)
	return config, nil

}
