package config

import (
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Database Database `yaml:"database"`
	Redis    Redis    `yaml:"redis"`
	Server   Server   `yaml:"http_server"`
}

type Redis struct {
	Addr     string `yaml:"addr" env:"REDIS_ADDR" env-default:"redis:6379"`
	Password string `yaml:"password" env:"REDIS_PASSWORD"`
	DB       int    `yaml:"db" env:"REDIS_DB" env-default:"0"`
}

type Database struct {
	Host         string        `yaml:"host" env:"DB_HOST"`
	Port         string        `yaml:"port" env:"DB_PORT"`
	User         string        `yaml:"user" env:"DB_USER"`
	Password     string        `yaml:"password" env:"DB_PASSWORD"`
	Name         string        `yaml:"name" env:"DB_NAME"`
	MaxPoolSize  int           `yaml:"max_pool_size" env:"DB_MAX_POOL_SIZE" env-default:"5"`
	ConnAttempts int           `yaml:"conn_attempts" env:"DB_CONN_ATTEMPTS" env-default:"10"`
	Timeout      time.Duration `yaml:"timeout" env:"DB_TIMEOUT" env-default:"1s"`
}

type Server struct {
	JwtSecret string        `yaml:"jwt-secret" env:"JWT_SECRET"`
	Timeout   time.Duration `yaml:"timeout" env:"SERVER_TIMEOUT" env-default:"5s"`
	Port      string        `yaml:"port" env:"SERVER_PORT" env-default:"8090"`
}

var configPath *string

func init() {
	configPath = flag.String("config_path", "", "Path to config file")
}

func NewConfig() Config {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	if *configPath == "" {
		*configPath = os.Getenv("CONFIG_PATH")
	}
	if *configPath == "" {
		logger.Warn("config path is not specified, using default \"config/config.yml\"")
		*configPath = "config/config.yml"
	}
	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		logger.Error("config file does not exist: " + *configPath)
		panic("config file does not exist: " + *configPath)
	}
	cfg := Config{}
	if err := cleanenv.ReadConfig(*configPath, &cfg); err != nil {
		logger.Error("cannot read yaml: " + err.Error())
		panic("cannot read yaml: " + err.Error())
	}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		logger.Error("cannot read env: " + err.Error())
		panic("cannot read env: " + err.Error())
	}
	return cfg
}
