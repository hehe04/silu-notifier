package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"silu-notifier/server/entity"
	"strings"
	"time"
)

var contentype = "application/json"
var baseUrl = "http://139.159.184.174:8000/"
var dst string = "/tmp"

func sub(e *entity.Device) {
	postData, err := json.Marshal(e)
	if err != nil {
		log.Fatal(err)
	}
	subUrl := baseUrl + "sub"
	for {
		respData, err := callSub(subUrl, postData)
		if err != nil {
			log.Println(err)
			continue
		}

		var msg entity.ResponseEvent
		json.Unmarshal(respData, &msg)
		if len(msg.Data) > 0 {
			for _, u := range msg.Data {
				log.Println("starting download ", u)
				names := strings.Split(u, "/")
				go func() {
					d := dst + "/" + names[len(names)-1]
					err := download(u, d)
					if err != nil {
						log.Println(err)
						return
					}
					showImg(d)
				}()
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func callSub(subUrl string, postData []byte) ([]byte, error) {
	resp, err := http.Post(subUrl, contentype, bytes.NewBuffer(postData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("%s/n", respData)
		return nil, errors.New(resp.Status)
	}
	return respData, nil
}

func download(url, dst string) error {
	url = baseUrl + url
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	m := float32(resp.ContentLength) / 1024 / 1024
	log.Printf("Download success, name: %s, lenth: %.2fMi", url, m)
	f, err := os.Create(dst)
	if err != nil {
		log.Println(err)
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err

}

func showImg(p string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("recover from a panic", r)
		}
	}()
	cmdO := exec.Command("open", p)
	err := cmdO.Start()
	if err != nil {
		log.Println(err)
	}
}

func main() {
	e := &entity.Device{
		DevUUID:     uuid.New().String(),
		Type:        2,
		BluetoothID: "1",
	}
	sub(e)
}
