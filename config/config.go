package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Serv ServerConfig   `yaml:"server"`
	PG   PostgresConfig `yaml:"postgres"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

type ServerConfig struct {
	Port string `yaml:"port" env:"SERVER_PORT" env-default:"8080"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"POSTGRES_PORT" env-default:"5432"`
	User     string `yaml:"user" env:"POSTGRES_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-default:"postgres"`
	Name     string `yaml:"name" env:"POSTGRES_DB" env-default:"tasktracker"`
}
