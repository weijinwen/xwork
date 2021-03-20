package Artisan

import (
	"reflect"
	"regexp"
)

//判断手机号
func IsMobile(str string) bool {
	regular := "^1[3456789]\\d{9}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(str)
}

//判断是否是结构体
func IsStruct(val interface{}) bool {
	if val == nil {
		return false
	}

	rv := reflect.ValueOf(val)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	return rv.Kind() == reflect.Struct
}
