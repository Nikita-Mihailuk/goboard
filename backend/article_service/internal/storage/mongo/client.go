package mongo

import (
	"context"
	"fmt"
	"github.com/Nikita-Mihailuk/goboard/backend/article_service/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Storage struct {
	db *mongo.Client
}

func NewStorage(cfg *config.Config) *Storage {
	connStr := fmt.Sprintf("mongodb://%s:%s/",
		cfg.DB.Host,
		cfg.DB.Port)

	opts := options.Client().ApplyURI(connStr)
	if cfg.DB.Username != "" && cfg.DB.Password != "" {
		opts = opts.SetAuth(options.Credential{
			Username: cfg.DB.Username,
			Password: cfg.DB.Password,
		})
	}

	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	return &Storage{
		db: client,
	}
}
