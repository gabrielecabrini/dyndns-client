package config

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Daemon    bool
	IpUrl     string `yaml:"ip-url"`
	Interval  time.Duration
	Providers map[string]map[string]interface{}
}

func GetConfig(path string) (*Config, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c Config) GetIP() (string, error) {
	resp, _ := http.Get(c.IpUrl)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	ip := string(data)
	log.Printf("Retrieved IP: %s", ip)
	return ip, nil
}
