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
	MkDirAll(BasePath)
}

// Subscribe the device on server
func Subscribe(c *gin.Context) {
	var dev entity.Device
	/*	body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Printf("%s", err)
			c.JSON(500, &entity.Response{
				Status:  http.StatusText(500),
				Message: err.Error(),
			})
			return
		}
		log.Printf("Request %s\n %s", "Subscribe", body)*/

	if err := c.ShouldBindJSON(&dev); err != nil {
		log.Println(err)
		c.JSON(400, &entity.Response{
			Status:  http.StatusText(400),
			Message: err.Error(),
		})
		return
	}
	dev.CreatedOn = time.Now()
	log.Printf("/sub | %+v", dev)
	entity.DeviceCache.Set(dev.DevUUID, &dev, 0)
	resp := entity.ResponseEvent{Data: []string{}}
	if img, ok := entity.EventCache.Get(dev.DevUUID); ok {
		e, ok := img.(*entity.Event)
		if !ok {
			panic(fmt.Sprintf(
				"Type convert failed, dev.Data of %s is not []string", dev.DevUUID))
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
		obj.Describe = fmt.Sprintf("expired: %s",
			time.Unix(0, v.Expiration).Format("2006-01-02T15:04:05Z07:00"))

		devs = append(devs, obj)
	}

	c.JSON(200, &entity.ResponseDevices{
		Response: entity.Response{
			Status: http.StatusText(200),
		},
		Count:   len(entity.DeviceCache.Items()),
		Devices: devs,
	})
}

func Publishe(c *gin.Context) {
	var event entity.Event
	err := c.ShouldBindJSON(&event)
	if err != nil {
		c.JSON(400, &entity.Response{
			Status:  http.StatusText(400),
			Message: err.Error(),
		})
		c.Error(err)
		return
	}
	//og.Printf("/pub | %+v",event)
	if _, ok := entity.DeviceCache.Get(event.To); !ok {
		c.JSON(200, &entity.Response{
			Status:  http.StatusText(200),
			Message: "Device not found or be offline",
		})
		return
	}

	err = entity.EventCache.Add(event.To, &event, 0)
	if err != nil {
		c.JSON(http.StatusConflict, &entity.Response{
			Status:  "Conflict",
			Message: "Consumer has a exception, message not be received",
		})
		c.Error(err)
		return
	}
	//log.Printf("Device %s send to %s: %s\n", event.From, event.To, event.Data)
	c.JSON(200, &entity.Response{
		Status:  http.StatusText(200),
		Message: "",
	})
}

func Upload(c *gin.Context) {
	var resp entity.ResponseEvent
	// single file
	file, _ := c.FormFile("file")
	//log.Println(file.Filename)
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
		//log.Println("save file error", err)
		resp.Status = http.StatusText(500)
		resp.Message = err.Error()
		c.JSON(500, &resp)
		c.Error(err)
		return
	}
	err = os.Chmod(p, 0666)
	if err != nil {
		c.Error(err)
	}
	resp.Status = http.StatusText(200)
	resp.Message = "Uploaded"
	resp.Data = []string{p}
	c.JSON(http.StatusOK, &resp)
}
