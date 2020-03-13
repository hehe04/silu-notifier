package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"io"
	"log"
	"os"
	"path/filepath"
	"silu-notifier/server/handler"
	"silu-notifier/server/middleware"
	"time"
)

func main() {
	addr := flag.String("addr", ":8000", "listening IP:port address, default 0.0.0.0:8000")
	//logInit()
	r := gin.Default()
	r.Use(middleware.LoggerToFile())
	r.Static("/resource", handler.BasePath)
	r.POST("/sub", handler.Subscribe)
	r.GET("/devices", handler.Devices)
	r.POST("/pub", handler.Publishe)
	r.POST("/upload", handler.Upload)
	if err := r.Run(*addr); err != nil {
		log.Fatal(err)
	}
}

func logInit() {
	gin.DisableConsoleColor()
	logdir := "./logs/"
	handler.MkDirAll(logdir)
	logdir, err := filepath.Abs(logdir)
	if err != nil {
		log.Fatal(err)
	}

	logdir = filepath.Join(logdir, "access.log")
	logWriter, err := rotatelogs.New(
		logdir+".%Y-%m-%d-%H-%M",
		rotatelogs.WithLinkName(logdir),           // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		log.Fatal(err)
	}
	gin.DefaultWriter = io.MultiWriter(logWriter, os.Stdout)
}
