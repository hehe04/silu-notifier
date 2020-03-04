package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()
	r.Static("/resource", "./resource")
	if err := r.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}
