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
	AuthService    AuthService    `yaml:"authService"`
	CommentService CommentService `yaml:"commentService"`
	Auth           Auth           `yaml:"auth"`
	Redis          Redis          `yaml:"redis"`
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

type AuthService struct {
	Host string `yaml:"host" env:"AUTH_SERVICE_HOST"`
	Port string `yaml:"port" env:"AUTH_SERVICE_PORT"`
}

type CommentService struct {
	Host string `yaml:"host" env:"COMMENT_SERVICE_HOST"`
	Port string `yaml:"port" env:"COMMENT_SERVICE_PORT"`
}

type Auth struct {
	SecretKey string `yaml:"secretKey" env:"SECRET_KEY"`
}

type Redis struct {
	Host     string        `yaml:"host" env:"REDIS_HOST"`
	Port     string        `yaml:"port" env:"REDIS_PORT"`
	Password string        `yaml:"password" env:"REDIS_PASSWORD"`
	DbNumber int           `yaml:"db" env:"REDIS_DB_NUMBER"`
	CacheTTL time.Duration `yaml:"cacheTTL" env:"REDIS_CACHE_TTL"`
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
