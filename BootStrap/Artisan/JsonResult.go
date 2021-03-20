package Artisan

import jsoniter "github.com/json-iterator/go"

//成功,站位其他错误不允许使用
const SUCCESS = 200

//未登陆,站位其他错误不允许使用
const NOTLOGIN = 401

//未找到页面,站位其他错误不允许使用
const NOTFOUND = 404

//服务器错误,站位其他错误不允许使用
const SERVICEFAIL = 500

type JsonResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func JsonData(code int, message string, data interface{}) JsonResult {
	return JsonResult{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func JsonSuccess(message string, data interface{}) JsonResult {
	return JsonResult{
		Code:    SUCCESS,
		Message: message,
		Data:    data,
	}
}

func JsonError(code int, message string, data interface{}) JsonResult {
	return JsonResult{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

//json解析
func JsonDecode(from string, to interface{}) error {
	var iJson = jsoniter.ConfigCompatibleWithStandardLibrary
	err := iJson.Unmarshal([]byte(from), &to)
	if err != nil {
		return err
	}
	return nil
}

//转json
func JsonEncode(from interface{}) ([]byte, error) {
	var iJson = jsoniter.ConfigCompatibleWithStandardLibrary
	return iJson.Marshal(&from)
}
