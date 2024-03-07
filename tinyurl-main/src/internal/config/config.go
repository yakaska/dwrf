package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string `yaml:"env" env:"ENV" env-required:"true"`
	Host     string `yaml:"host" env:"HOST" env-required:"true"`
	Port     int    `yaml:"port" env:"PORT" env-required:"true"`
	DbConfig `yaml:"database"`
}

type DbConfig struct {
	DbHost       string `yaml:"host" env:"DB_HOST" env-required:"true"`
	DbPort       int    `yaml:"port" env:"DB_PORT" env-required:"true"`
	DbUser       string `yaml:"user" env:"DB_USER" env-required:"true"`
	DbPassword   string `yaml:"password" env:"DB_USER_PASSWORD" env-required:"true"`
	DbName       string `yaml:"name" env:"DB_NAME" env-required:"true"`
	DbConTimeout int    `yaml:"timeout" env:"DB_CON_TIMEOUT" env-required:"true"`
}

func Load() *Config {
	var config Config

	if err := cleanenv.ReadEnv(&config); err != nil {
		panic("Can't read config from environment variables: " + err.Error())
	}
	return &config
}

func (c Config) Dsn() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.DbHost,
		c.DbPort,
		c.DbUser,
		c.DbPassword,
		c.DbName,
	)
}

func (c Config) BaseUrl() string {
	return fmt.Sprintf("http://%s:%d", c.Host, c.Port)
}
func (c Config) ListenAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
