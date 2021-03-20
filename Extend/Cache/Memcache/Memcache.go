package Memcache

import (
	"xwork/BootStrap/Artisan"
	cache "xwork/Extend/Cache"
	"xwork/Extend/Gophp"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/dgryski/go-fastlz"
	"github.com/gookit/validate"
	"github.com/spf13/viper"
)

// Cache Memcache adapter.
type Cache struct {
	conn     *memcache.Client
	server   string
	conninfo []string
}

// NewMemCache create new memcache adapter.
func NewMemCache() cache.Cache {
	return &Cache{}
}

// Get get value from memcache.
func (rc *Cache) Get(key string) interface{} {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return nil
		}
	}
	if item, err := rc.conn.Get(key); err == nil {
		item := rc.uncompress(item)
		return rc.unserialize(item)
	}
	return nil
}

// GetPic get value from memcache.
// 作品缓存使用
func (rc *Cache) GetPic(key string) interface{} {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return nil
		}
	}
	if item, err := rc.conn.Get(key); err == nil {
		item := rc.uncompress(item)
		return rc.unserializePic(item)
	}
	return nil
}

// GetMulti get value from memcache.
func (rc *Cache) GetMulti(keys []string) map[string]interface{} {
	var rv = make(map[string]interface{})
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return nil
		}
	}
	mv, err := rc.conn.GetMulti(keys)
	if err == nil && len(mv) > 0 {
		for k, v := range mv {
			v = rc.uncompress(v)
			d := rc.unserialize(v)
			rv[k] = d
		}
		return rv
	} else {
		return nil
	}
}

// GetPicMulti get value from memcache.
// 作品缓存使用
func (rc *Cache) GetPicMulti(keys []string) map[string]interface{} {
	var rv = make(map[string]interface{})
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return nil
		}
	}
	mv, err := rc.conn.GetMulti(keys)
	if err == nil && len(mv) > 0 {
		for k, v := range mv {
			v := rc.uncompress(v)
			d := rc.unserializePic(v)
			rv[k] = d
		}
		return rv
	} else {
		return nil
	}
}

// Put put value to memcache.
func (rc *Cache) Put(key string, val interface{}, timeout int32) error {
	var err error
	if rc.conn == nil {
		if err = rc.connectInit(); err != nil {
			return err
		}
	}
	item := &memcache.Item{Key: key, Expiration: timeout}
	item, err = rc.serialize(item, val)
	if err != nil {
		return err
	}
	item = rc.compress(item)
	return rc.conn.Set(item)
}

// PutPic put value to memcache.
// 作品缓存使用
func (rc *Cache) PutPic(key string, val interface{}, timeout int32) error {
	var err error
	if rc.conn == nil {
		if err = rc.connectInit(); err != nil {
			return err
		}
	}
	item := &memcache.Item{Key: key, Expiration: timeout}
	item, err = rc.serializePic(item, val)
	if err != nil {
		return err
	}
	item = rc.compress(item)
	return rc.conn.Set(item)
}

// PutMulti put value to memcache.
func (rc *Cache) PutMulti(kvs map[string]interface{}, timeout int32) error {
	var err error
	if rc.conn == nil {
		if err = rc.connectInit(); err != nil {
			return err
		}
	}
	for k, v := range kvs {
		item := &memcache.Item{Key: k, Expiration: timeout}
		item, err = rc.serialize(item, v)
		if err != nil {
			return err
		}
		item = rc.compress(item)
		rc.conn.Set(item)
	}

	return nil
}

// Delete delete value in memcache.
func (rc *Cache) Delete(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	return rc.conn.Delete(key)
}

// Incr increase counter.
func (rc *Cache) Incr(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Increment(key, 1)
	return err
}

// Decr decrease counter.
func (rc *Cache) Decr(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Decrement(key, 1)
	return err
}

// IsExist check value exists in memcache.
func (rc *Cache) IsExist(key string) bool {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return false
		}
	}
	_, err := rc.conn.Get(key)
	return err == nil
}

// ClearAll clear all cached in memcache.
func (rc *Cache) ClearAll() error {
	return nil
}

// StartAndGC start memcache adapter.
// config string is like {"conn":"connection info"}.
// if connecting error, return.
func (rc *Cache) StartAndGC(server string) error {
	rc.server = server
	rc.conninfo = strings.Split(viper.GetString(fmt.Sprintf("cache.%s.host", server)), ";")
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	return nil
}

// connect to memcache and keep the connection.
func (rc *Cache) connectInit() error {
	timeout := viper.GetInt64(fmt.Sprintf("cache.%s.timeout", rc.server))
	initConns := viper.GetInt(fmt.Sprintf("cache.%s.initConns", rc.server))
	rc.conn = memcache.New(rc.conninfo...)
	rc.conn.Timeout = time.Duration(timeout) * time.Second
	rc.conn.MaxIdleConns = initConns
	return nil
}

func (rc *Cache) uncompress(item *memcache.Item) *memcache.Item {
	if item.Flags != 84 && item.Flags != 80 {
		return item
	}
	var b []byte
	b, err := fastlz.Decode(b, item.Value)
	if err != nil {
		return item
	}
	item.Value = b
	return item
}

func (rc *Cache) compress(item *memcache.Item) *memcache.Item {
	if len(item.Value) < 3000 {
		return item
	}
	var b []byte
	var err error
	switch item.Flags {
	case 0:
		b, err = fastlz.Encode(b, item.Value)
		if err != nil {
			return item
		}
		item.Flags = 80
		item.Value = b
	case 4:
		b, err = fastlz.Encode(b, item.Value)
		if err != nil {
			return item
		}
		item.Flags = 84
		item.Value = b
	}
	return item
}

//json序列化，主要提供正常数据
func (rc *Cache) serialize(item *memcache.Item, val interface{}) (*memcache.Item, error) {
	if validate.IsString(val) {
		item.Flags = 0
		item.Value = []byte(fmt.Sprintf("%v", val))
	} else if validate.IsNumeric(val) {
		item.Flags = 1
		item.Value = []byte(fmt.Sprintf("%v", val))
	} else if validate.IsFloat(val) {
		item.Flags = 2
		item.Value = []byte(fmt.Sprintf("%v", val))
	} else if validate.IsBool(val) {
		var boolInt int8
		switch val.(bool) {
		case true:
			boolInt = 1
		default:
			boolInt = 0
		}
		item.Flags = 3
		item.Value = []byte(fmt.Sprintf("%v", boolInt))
	} else if validate.IsArray(val) || validate.IsSlice(val) || validate.IsMap(val) || Artisan.IsStruct(val) {
		b, err := Artisan.JsonEncode(val)
		if err != nil {
			return nil, err
		}
		item.Value = b
		item.Flags = 4
	} else {
		return nil, errors.New("val only support string and []byte")
	}
	return item, nil
}

//序列化数据，兼容PHP主要提供作品缓存
func (rc *Cache) serializePic(item *memcache.Item, val interface{}) (*memcache.Item, error) {
	if validate.IsString(val) {
		item.Flags = 0
		item.Value = []byte(fmt.Sprintf("%v", val))
	} else if validate.IsNumeric(val) {
		item.Flags = 1
		item.Value = []byte(fmt.Sprintf("%v", val))
	} else if validate.IsFloat(val) {
		item.Flags = 2
		item.Value = []byte(fmt.Sprintf("%v", val))
	} else if validate.IsBool(val) {
		var boolInt int8
		switch val.(bool) {
		case true:
			boolInt = 1
		default:
			boolInt = 0
		}
		item.Flags = 3
		item.Value = []byte(fmt.Sprintf("%v", boolInt))
	} else if validate.IsArray(val) || validate.IsSlice(val) || validate.IsMap(val) {
		b, err := Gophp.Serialize(val)
		if err != nil {
			return nil, err
		}
		item.Value = b
		item.Flags = 4
	} else if Artisan.IsStruct(val) {
		return nil, errors.New("Can not is Struct!")
	} else {
		return nil, errors.New("val only support string and []byte")
	}
	return item, nil
}

//0 is_string
//1 is_number
//2 is_float
//3 is_bool
//4,84 is_serialize
//80 is_compress
//反序列化数据，主要提供正常缓存
func (rc *Cache) unserialize(item *memcache.Item) interface{} {
	var data interface{}
	var err error
	switch item.Flags {
	case 0:
		data = fmt.Sprintf("%s", item.Value)
	case 1:
		data, err = strconv.ParseInt(fmt.Sprintf("%s", item.Value), 10, 64)
	case 2:
		data, err = strconv.ParseFloat(fmt.Sprintf("%s", item.Value), 64)
	case 3:
		data, err = strconv.ParseBool(fmt.Sprintf("%s", item.Value))
	case 4, 84:
		err = Artisan.JsonDecode(string(item.Value), &data)
	default:
		data = fmt.Sprintf("%s", item.Value)
	}
	if err != nil {
		return fmt.Sprintf("%s", item.Value)
	} else {
		return data
	}
}

//0 is_string
//1 is_number
//2 is_float
//3 is_bool
//4,84 is_serialize
//80 is_compress
//反序列化数据，主要提供给作品缓存
func (rc *Cache) unserializePic(item *memcache.Item) interface{} {
	var data interface{}
	var err error
	switch item.Flags {
	case 0:
		data = fmt.Sprintf("%s", item.Value)
	case 1:
		data, err = strconv.ParseInt(fmt.Sprintf("%s", item.Value), 10, 64)
	case 2:
		data, err = strconv.ParseFloat(fmt.Sprintf("%s", item.Value), 64)
	case 3:
		data, err = strconv.ParseBool(fmt.Sprintf("%s", item.Value))
	case 4, 84:
		data, err = Gophp.Unserialize(item.Value)
	default:
		data = fmt.Sprintf("%s", item.Value)
	}
	if err != nil {
		return fmt.Sprintf("%s", item.Value)
	} else {
		return data
	}
}

func init() {
	cache.Register("memcache", NewMemCache)
}
