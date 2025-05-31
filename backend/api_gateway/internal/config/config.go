package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"sync"
	"time"
)

type Config struct {
	HTTPServer     HTTPServer     `yaml:"server"`
	UserService    UserService    `yaml:"userService"`
	ArticleService ArticleService `yaml:"articleService"`
}

type HTTPServer struct {
	Port    string        `yaml:"port" env:"HTTP_PORT"`
	Timeout time.Duration `yaml:"timeout" env:"GRPC_TIMEOUT"`
}

type UserService struct {
	Host string `yaml:"host" env:"USER_SERVICE_HOST"`
	Port string `yaml:"port" env:"USER_SERVICE_PORT"`
}

type ArticleService struct {
	Host string `yaml:"host" env:"ARTICLE_SERVICE_HOST"`
	Port string `yaml:"port" env:"ARTICLE_SERVICE_PORT"`
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
