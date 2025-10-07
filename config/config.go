package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	App AppConfig `toml:"app"`
	DB  DBConfig  `toml:"db"`
}

type AppConfig struct {
	WebsiteName string `toml:"website_name"`
	Env         string `toml:"env"`
	Port        int    `toml:"port"`
	Hostname    string `toml:"hostname"`
	AutoMigrate bool   `toml:"auto_migrate"`
}

type DBConfig struct {
	Type     string `toml:"type"`
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Name     string `toml:"name"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
