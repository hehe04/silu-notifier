package entity

import (
	"time"
)

// Device Type
const (
	DeviceTypePC     = 1
	DeviceTypeMobile = 2
)

// Device mobile or PC
type Device struct {
	DevUUID   string    `json:"dev_uuid" binding:"required"`
	Describe  string    `json:"describe"`
	//NotifyURL string    `json:"notify_url" binding:"required"`
	BluetoothID string `json:"bluetooth_id" binding: "required"`
	CreatedOn time.Time `json:"created_on"`
	Expired   time.Time
	Type      uint `json:"type"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ResponseDevices struct {
	Response
	Count int `json:"count"`
	Devices []*Device `json:"devices"`
}

type ResponseEvent struct {
	Response
	Data []string `json:"data"`
}

type Event struct {
	From string `json:"from"`
	To string `json:"to"`
	Data []string `json:"data"`
}