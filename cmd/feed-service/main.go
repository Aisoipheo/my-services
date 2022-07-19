package main

import (
	"github.com/gin-gonic/gin"

	"my-service/internal/models"
	"my-service/internal/middleware"
	"my-service/pkg/db/postgres"
)

func main() {
	var cfg models.Config

	cfg.PostgresUser.GetEnv("POSTGRES_USER")
	cfg.PostgresPassword.GetEnv("POSTGRES_PASSWORD")
	cfg.PostgresDBName.GetEnv("POSTGRES_DBNAME")
	cfg.PostgresHost.GetEnv("POSTGRES_HOST")
	cfg.PostgresPort.GetEnv("POSTGRES_PORT")
	cfg.RouterHost.GetEnv("ROUTER_HOST")
	cfg.RouterPort.GetEnv("ROUTER_PORT")

	postgreSQLConfig := postgres.PostgreSQLConfig{
		User	: cfg.PostgresUser.String(),
		Password: cfg.PostgresPassword.String(),
		DBName	: cfg.PostgresDBName.String(),
		Host	: cfg.PostgresHost.String(),
		Port	: cfg.PostgresPort.String(),
	}

	conn, err := postgres.NewPostgresDB(&postgreSQLConfig);
	if err != nil {
		panic("PostgreSQL connection failed")
	}

	ctrl := middleware.Controller {
		&cfg,
		conn,
		"0.0.1-alpha",
	}

	router := gin.Default()
	router.POST("/likes", ctrl.PostLike)
	router.POST("/dislikes", ctrl.PostDislike)
	router.POST("/new-post", ctrl.PostNewPost)
	router.GET("/posts", ctrl.GetPosts)
	router.GET("/healthz", ctrl.GetHealthz)

	router.Run(cfg.RouterHost.String() + ":" + cfg.RouterPort.String())
}
