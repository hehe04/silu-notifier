package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"silu-notifier/server/handler"
)

func main() {
	r := gin.Default()
	r.Static("/resource", handler.BasePath)
	r.POST("/sub", handler.Subscribe)
	r.GET("/devices", handler.Devices)
	r.POST("/pub", handler.Publishe)
	r.POST("/upload", handler.Upload)
	if err := r.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}
