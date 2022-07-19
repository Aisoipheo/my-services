package main

import (
	"github.com/gin-gonic/gin"

	"my-service/internal/entity"
	"my-service/internal/posts"
	"my-service/internal/healthz"
	"my-service/pkg/db/postgres"
)

func main() {
	var cfg Config

	cfg.PostgresUser.GetEnv("POSTGRES_USER")
	cfg.PostgresPassword.GetEnv("POSTGRES_PASSWORD")
	cfg.PostgresDBName.GetEnv("POSTGRES_DBNAME")
	cfg.PostgresHost.GetEnv("POSTGRES_HOST")
	cfg.PostgresPort.GetEnv("POSTGRES_PORT")
	cfg.RouterHost.GetEnv("ROUTER_HOST")
	cfg.RouterPort.GetEnv("ROUTER_PORT")

	postgreSQLConfig := PostgreSQLConfig{
		User	: cfg.PostgresUser.String(),
		Password: cfg.PostgresPassword.String(),
		DBName	: cfg.PostgresDBName.String(),
		Host	: cfg.PostgresHost.String(),
		Port	: cfg.PostgresPort.String()
	}

	conn, err := NewPostgresDB(postgreSQLConfig);
	if err != nil {
		panic("PostgreSQL connection failed")
	}

	ctrl := Controller {
		&cfg,
		&postgreSQLConfig,
		"0.0.1-alpha"
	}

	router := gin.Default()
	router.POST("/likes", ctrl.postLike)
	router.POST("/dislikes", ctrl.postDislike)
	router.POST("/new-post", ctrl.postNewPost)
	router.GET("/posts", ctrl.getPosts)
	router.GET("/healthz", ctrl.getHealthz)

	router.Run(cfg.RouterHost + ":" + cfg.RouterPort)
}
