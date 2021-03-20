package Cache

import (
	"fmt"

	"github.com/spf13/viper"
)

type Cache interface {
	// Get cached value by key.
	Get(key string) interface{}
	// GetPic cached value by key.
	GetPic(key string) interface{}
	// GetMulti is a batch version of Get.
	GetMulti(keys []string) map[string]interface{}
	// GetPicMulti is a batch version of Get.
	GetPicMulti(keys []string) map[string]interface{}
	// set cached value with key and expire time.
	Put(key string, val interface{}, timeout int32) error
	// set cached value with key and expire time.
	PutMulti(kvs map[string]interface{}, timeout int32) error
	// set cached value with key and expire time.
	PutPic(key string, val interface{}, timeout int32) error
	// delete cached value by key.
	Delete(key string) error
	// increase cached int value by key, as a counter.
	Incr(key string) error
	// decrease cached int value by key, as a counter.
	Decr(key string) error
	// check if cached value exists or not.
	IsExist(key string) bool
	// clear all cache.
	ClearAll() error
	// start gc routine based on config string settings.
	StartAndGC(server string) error
}

// Instance is a function create a new Cache Instance
type Instance func() Cache

var adapters = make(map[string]Instance)

var servers = make(map[string]Cache)

// Register makes a cache adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, adapter Instance) {
	if adapter == nil {
		panic("cache: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("cache: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

// NewCache Create a new cache driver by adapter name and config string.
// config need to be correct JSON as string: {"interval":360}.
// it will start gc automatically.
func NewCache(server string) (adapter Cache) {
	config := viper.GetStringMapString(fmt.Sprintf("cache.%s", server))
	adaptersName, ok := config["server"]
	if !ok {
		panic(fmt.Sprintf("cache: server is null name %q", server))
	}
	if adapter, ok := servers[server]; ok && adapter != nil {
		return adapter
	}
	instanceFunc, ok := adapters[adaptersName]
	if !ok {
		panic(fmt.Sprintf("cache: unknown adapter name %q (forgot to import?)", config["server"]))
	}
	_, ok = config["host"]
	if !ok {
		panic(fmt.Sprintf("cache: unknown host name %q", server))
	}
	adapter = instanceFunc()
	err := adapter.StartAndGC(server)
	if err != nil {
		panic(fmt.Sprintf("cache: server connection fail name %q", server))
	}
	servers[server] = adapter
	return
}
