package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"sync"
)

type Config struct {
	Env        string   `yaml:"env" env:"ENV"`
	GRPCServer GRPC     `yaml:"grpc"`
	DB         DataBase `yaml:"db"`
	Kafka      Kafka    `yaml:"kafka"`
}

type GRPC struct {
	Port string `yaml:"port" env:"PORT"`
}

type DataBase struct {
	Host     string `yaml:"host" env:"DATABASE_HOST"`
	Port     string `yaml:"port" env:"DATABASE_PORT"`
	Username string `yaml:"userName" env:"DATABASE_USERNAME"`
	Password string `yaml:"password" env:"DATABASE_PASSWORD"`
}

type Kafka struct {
	Address          string `yaml:"address" env:"KAFKA_ADDRESS"`
	UserServiceTopic string `yaml:"userServiceTopic" env:"KAFKA_USER_SERVICE_TOPIC"`
	ConsumerGroup    string `yaml:"consumerGroup" env:"KAFKA_CONSUMER_GROUP"`
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
