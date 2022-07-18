package main

import (

)

func main() {
	router := gin.Default()
	router.POST("/likes", postLike)
	router.POST("/dislikes", postDislike)
	router.POST("/new-post"), postNewPost)
	router.GET("/posts", getPosts)
	router.GET("/healthz", getHealthz)

	router.Run() // TODO add ip:port from env SERVICE_ENGINE_IP SERVICE_ENGINE_PORT
}
