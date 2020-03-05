package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"silu-notifier/server/entity"
	"time"
)

// Register the device on server
func Register(c *gin.Context) {
	var dev entity.Device
	if err := c.ShouldBind(&dev); err != nil {
		log.Println(err)
		c.JSON(400, &entity.Response{
			Status:  http.StatusText(400),
			Message: err.Error(),
		})
		return
	}
	dev.CreatedOn = time.Now()
	entity.Cache.Set(dev.DevUUID, &dev, time.Second*10)
	c.JSON(200, &entity.Response{
		Status:  http.StatusText(200),
		Message: "",
	})
}

// List all of devices
func Devices(c *gin.Context) {
	var devs []*entity.Device
	for _, v := range entity.Cache.Items() {
		obj := v.Object.(*entity.Device)
		obj.Describe = fmt.Sprintf("expired: %v", v.Expired())
		devs = append(devs, obj)
	}

	c.JSON(200, &entity.ResponseDevices{
		Response: entity.Response{
			Status: http.StatusText(200),
		},
		Count:   entity.Cache.ItemCount(),
		Devices: devs,
	})
}
