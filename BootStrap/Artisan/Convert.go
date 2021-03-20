package Artisan

import (
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

// interface数组转int数组
func InterfacesToInts(from []interface{}) (to []int) {
	for _, v := range from {
		to = append(to, v.(int))
	}
	return to
}

// interface数组转string数组
func InterfacesToStrings(from []interface{}) (to []string) {
	for _, v := range from {
		to = append(to, v.(string))
	}
	return to
}

// 数字转int64，统一转成int64防止转类型错误
// 无法确定传入的数字类型，主要缓存写入读取类型可能会有变化
func NumberToInt64(val interface{}) (valInt64 int64) {
	switch val.(type) {
	case string:
		valInt64, _ = strconv.ParseInt(val.(string), 10, 64)
	case int64:
		valInt64 = val.(int64)
	case int8:
		valInt64 = int64(val.(int8))
	case int16:
		valInt64 = int64(val.(int16))
	case int:
		valInt64 = int64(val.(int))
	case int32:
		valInt64 = int64(val.(int32))
	case uint8:
		valInt64 = int64(val.(uint8))
	case uint16:
		valInt64 = int64(val.(uint16))
	case uint:
		valInt64 = int64(val.(uint))
	case uint32:
		valInt64 = int64(val.(uint32))
	case uint64:
		valInt64 = int64(val.(uint64))
	default:
		valInt64 = 0
	}
	return valInt64
}

// map转struct
func MapToStruct(from map[string]interface{}, to interface{}) error {
	var iJson = jsoniter.ConfigCompatibleWithStandardLibrary
	fromJson, err := iJson.Marshal(&from)
	if err != nil {
		return err
	}
	err = iJson.Unmarshal(fromJson, to)
	if err != nil {
		return err
	}
	return nil
}
