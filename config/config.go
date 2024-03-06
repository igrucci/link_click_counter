package config

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

func NewConfig() (*Config, error) {

	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		logrus.Fatalf("error: %v", err)
	}
	cfg := &Config{}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		logrus.Fatalf("error: %v", err)
	}
	return cfg, nil
}
