package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"io"
	"log"
	"os"
	"path/filepath"
	"silu-notifier/server/config"
	"silu-notifier/server/handler"
	"silu-notifier/server/middleware"
	"time"
)

func main() {
	svcConfig := &service.Config{
		Name:        "GeoNotifierServer", //服务显示名称
		DisplayName: "GeoNotifierServer", //服务名称
		Description: "两馆蓝牙定位通知服务",        //服务描述
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		logger.Error(err)
	}

	if err != nil {
		logger.Error(err)
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			if config.HomePath == "" {
				logger.Error("请设置环境变量 GEO_NOTIFY_SERVER_HOME")
				return
			}
			logger.Info("resource: ", config.ResourcePath, " home: ", config.HomePath)
			handler.MkDirAll(config.ResourcePath)

			s.Install()
			logger.Info("服务安装成功!")
			s.Start()
			logger.Info("服务启动成功!")
			break
		case "start":
			s.Start()
			logger.Info("服务启动成功!")
			break
		case "stop":
			s.Stop()
			logger.Info("服务关闭成功!")
			break
		case "restart":
			s.Stop()
			logger.Info("服务关闭成功!")
			s.Start()
			logger.Info("服务启动成功!")
			break
		case "remove":
			s.Stop()
			logger.Info("服务关闭成功!")
			s.Uninstall()
			logger.Info("服务卸载成功!")
			break
		}
		return
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}

func logInit() {
	gin.DisableConsoleColor()
	logdir := filepath.Join(config.HomePath, "./logs/")
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

var logger service.Logger = service.ConsoleLogger

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	// Do work here
	addr := flag.String("addr", ":8000", "listening IP:port address, default 0.0.0.0:8000")
	//logInit()
	flag.Parse()
	r := gin.Default()
	r.Use(middleware.LoggerToFile())
	r.Static("/resource", config.ResourcePath)
	r.POST("/sub", handler.Subscribe)
	r.GET("/devices", handler.Devices)
	r.POST("/pub", handler.Publishe)
	r.POST("/upload", handler.Upload)
	if err := r.Run(*addr); err != nil {
		log.Fatal(err)
	}

}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}
