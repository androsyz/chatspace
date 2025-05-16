package config

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func NewViper() (*Config, error) {
	cfg := &Config{}

	if err := loadConfiguration(cfg); err != nil {
		log.Err(err).Msg("Failed to load configuration")
		return nil, err
	}

	log.Info().Msg("Configuration loaded successfully")
	return cfg, nil
}

func loadConfiguration(cfg *Config) error {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
		log.Warn().Msg("Config file not found. Using environment variables or defaults.")
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()

	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}
