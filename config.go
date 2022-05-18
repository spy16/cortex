package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/chunked-app/cortex/pkg/config"
)

type appConfig struct {
	Port      int    `mapstructure:"port" default:"8080"`
	LogLevel  string `mapstructure:"log_level" default:"info"`
	LogPretty bool   `mapstructure:"log_pretty" default:"true"`
	Database  string `mapstructure:"database" default:":memory:"`
}

func loadConfig(cmd *cobra.Command) appConfig {
	opts := []config.Option{
		config.WithEnv(),
		config.WithName("cortex"),
	}

	var cfg appConfig
	if err := config.Load(&cfg, opts...); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	setupLogging(cfg)
	return cfg
}

func setupLogging(cfg appConfig) {
	lvl, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		lvl = log.WarnLevel
	}

	log.SetLevel(lvl)
	if cfg.LogPretty {
		log.SetFormatter(&log.TextFormatter{})
	} else {
		log.SetFormatter(&log.JSONFormatter{})
	}
}
