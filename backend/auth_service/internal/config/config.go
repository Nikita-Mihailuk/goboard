package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"sync"
	"time"
)

type Config struct {
	Env         string      `yaml:"env" env:"ENV"`
	GRPCServer  GRPC        `yaml:"grpc"`
	DB          DataBase    `yaml:"db"`
	Auth        Auth        `yaml:"auth"`
	UserService UserService `yaml:"userService"`
}

type GRPC struct {
	Port    string        `yaml:"port" env:"PORT"`
	Timeout time.Duration `yaml:"timeout" env:"GRPC_TIMEOUT"`
}

type DataBase struct {
	Host     string `yaml:"host" env:"DATABASE_HOST"`
	Port     string `yaml:"port" env:"DATABASE_PORT"`
	Password string `yaml:"password" env:"DATABASE_PASSWORD"`
	DbNumber int    `yaml:"db" env:"DATABASE_NUMBER"`
}

type Auth struct {
	AccessTokenTTL  time.Duration `yaml:"accessTokenTTL" env:"ACCESS_TOKEN_TTL"`
	RefreshTokenTTL time.Duration `yaml:"refreshTokenTTL" env:"REFRESH_TOKEN_TTL"`
	SecretKey       string        `yaml:"secretKey" env:"SECRET_KEY"`
}

type UserService struct {
	Host string `yaml:"host" env:"USER_SERVICE_HOST"`
	Port string `yaml:"port" env:"USER_SERVICE_PORT"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		path := fetchConfigPath()
		instance = LoadConfigByPath(path)
	})

	return instance
}

func LoadConfigByPath(path string) *Config {
	var cfg Config

	if path != "" {
		err := cleanenv.ReadConfig(path, &cfg)
		if err != nil {
			panic(err)
		}
	}

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "config file path")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
