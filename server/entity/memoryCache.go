package entity

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var DeviceCache *cache.Cache
var EventCache *cache.Cache

func init() {
	DeviceCache = cache.New(time.Second*5, time.Second*10)
	EventCache=cache.New(0,0) //not expired
}
