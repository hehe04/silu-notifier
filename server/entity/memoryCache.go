package entity

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var DeviceCache *cache.Cache
var EventCache *cache.Cache

func init() {
	DeviceCache = cache.New(time.Second*5, time.Hour*2)
	EventCache = cache.New(0, time.Hour*2) //not expired
}
