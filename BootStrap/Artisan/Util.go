package Artisan

import (
	"xwork/Extend/IpAddress"
	"fmt"
	"github.com/spf13/viper"
	"reflect"
	"strings"
)

// 根据不同环境添加缓存前缀
// @param key string 缓存key
// @return string 新缓存key
// @author wf
// @time 2020-11-24
func GetEnvCacheKey(key string) string {
	env := viper.GetString("Env")
	if env != "prod" {
		key = fmt.Sprintf("dev_%s", key)
	}
	return key
}

// 查询IP地区
// @param ipAddr string IP地址
// @return string 国家
// @return string 省/直辖市
// @return string 市/区
// @author wf
// @time 2020-11-24
func IpAddr(ipAddr string) (string, string, string) {
	ipAddress := IpAddress.Find(ipAddr)
	addr := strings.Split(ipAddress, "\t")
	lenAddr := len(addr)
	if lenAddr > 2 {
		return addr[0], addr[1], addr[2]
	} else if lenAddr > 1 {
		return addr[0], addr[1], ""
	} else if lenAddr > 0 {
		return addr[0], "", ""
	} else {
		return "", "", ""
	}
}

// 获取查询struct字段
// gorm:column 解析要查询的字段名
// table:"*" 解析表前缀
// @param model interface{}
// @param table string
// @return []string e.g. (a.id,b.id)
func GetFiledName(model interface{}) []string {
	t := reflect.TypeOf(model)
	ss := t.Elem()
	var fileds []string
	for i := 0; i < ss.NumField(); i++ {
		tags := strings.Split(ss.Field(i).Tag.Get("gorm"), ";")
		table := ss.Field(i).Tag.Get("table")
		for _, v := range tags {
			if !strings.Contains(v, "column") {
				continue
			}
			column := strings.Split(v, ":")
			filed := column[1]
			if filed == "-" {
				continue
			}
			if table != "" {
				filed = fmt.Sprintf("%s.%s", table, filed)
			}
			fileds = append(fileds, filed)
		}
	}
	return fileds
}
