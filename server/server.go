package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"silu-notifier/server/handler"
)

func main() {
	r := gin.Default()
	r.Static("/resource", "./resource")
	r.POST("/register", handler.Register)
	r.GET("/devices", handler.Devices)
	if err := r.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}
