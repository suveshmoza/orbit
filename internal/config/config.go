package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type ConfigFile struct {
	Config     AppConfig `json:"config"`
	DNSServers []Server  `json:"servers"`
}

type AppConfig struct {
	Samples     int      `json:"samples"`
	TimeoutMs   int      `json:"timeout_ms"`
	TestDomains []string `json:"test_domains"`
}

type Server struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Port    int    `json:"port"`
}

func LoadConfig(path string) (ConfigFile, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return ConfigFile{}, fmt.Errorf("read config: %w", err)
	}

	var cfg ConfigFile
	if err := json.Unmarshal(content, &cfg); err != nil {
		return ConfigFile{}, fmt.Errorf("parse config: %w", err)
	}

	return cfg, nil
}
