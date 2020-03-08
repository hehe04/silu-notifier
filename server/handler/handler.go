package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"silu-notifier/server/entity"
	"time"
)

var (
	BasePath string = "./resource"
)

func init() {
	if !Exists(BasePath) {
		err := os.MkdirAll(BasePath, 0766) // default will be effected by umask
		if err != nil {
			log.Fatal(err)
			return
		}
		if err := os.Chmod(BasePath, 0766); err != nil {
			log.Fatal(err)
		}
	}
}

// Subscribe the device on server
func Subscribe(c *gin.Context) {
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
	entity.DeviceCache.Set(dev.DevUUID, &dev, time.Second*15)
	resp := entity.ResponseEvent{Data: []string{}}
	if img, ok := entity.EventCache.Get(dev.DevUUID); ok {
		e, ok := img.(*entity.Event)
		if !ok {
			log.Printf("Type convert failed, dev.Data of %s is not []string", dev.DevUUID)
		} else {
			resp.Data = e.Data
		}
	}
	resp.Status = http.StatusText(200)
	c.JSON(200, &resp)
	entity.EventCache.Delete(dev.DevUUID)
}

// List all of devices
func Devices(c *gin.Context) {
	var devs []*entity.Device
	for _, v := range entity.DeviceCache.Items() {
		obj := v.Object.(*entity.Device)
		obj.Describe = fmt.Sprintf("expired: %v", v.Expired())
		devs = append(devs, obj)
	}

	c.JSON(200, &entity.ResponseDevices{
		Response: entity.Response{
			Status: http.StatusText(200),
		},
		Count:   entity.DeviceCache.ItemCount(),
		Devices: devs,
	})
}

func Publishe(c *gin.Context) {
	var event entity.Event
	err := c.ShouldBind(&event)
	if err != nil {
		log.Println(err)
		c.JSON(400, &entity.Response{
			Status:  http.StatusText(400),
			Message: err.Error(),
		})
		return
	}
	if _, ok := entity.DeviceCache.Get(event.To); !ok {
		c.JSON(200, &entity.Response{
			Status:  http.StatusText(200),
			Message: "Device not found or be offline",
		})
		return
	}

	err = entity.EventCache.Add(event.To, &event, 0)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusConflict, &entity.Response{
			Status:  "Conflict",
			Message: "Consumer has a exception, message not be received",
		})
		return
	}
	log.Printf("Device %s send to %s: %s\n", event.From, event.To, event.Data)
	c.JSON(200, &entity.Response{
		Status:  http.StatusText(200),
		Message: "",
	})
}

func Upload(c *gin.Context) {
	var resp entity.ResponseEvent
	// single file
	file, _ := c.FormFile("file")
	log.Println(file.Filename)
	/*p, err := filepath.Abs(BasePath)
	if err != nil {
		log.Println(err)
		resp.Status = http.StatusText(500)
		resp.Message = err.Error()
		c.JSON(500, &resp)
		return
	}*/
	p := filepath.Join(BasePath, file.Filename)
	// Upload the file to specific dst.
	err := c.SaveUploadedFile(file, p)
	if err != nil {
		log.Println("save file error", err)
		resp.Status = http.StatusText(500)
		resp.Message = err.Error()
		c.JSON(500, &resp)
		return
	}
	resp.Status = http.StatusText(200)
	resp.Message = "Uploaded"
	resp.Data = []string{p}
	c.JSON(http.StatusOK, &resp)
}
