package util

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

// 定义JSON操作
var (
	json              = jsoniter.ConfigCompatibleWithStandardLibrary
	JSONMarshal       = json.Marshal
	JSONUnmarshal     = json.Unmarshal
	JSONMarshalIndent = json.MarshalIndent
	JSONNewDecoder    = json.NewDecoder
	JSONNewEncoder    = json.NewEncoder
)

// JSONMarshalToString JSON编码为字符串
func JSONMarshalToString(v interface{}) string {
	s, err := jsoniter.MarshalToString(v)
	if err != nil {
		return ""
	}
	return s
}

func JsonUnmarshalFromString(str string, v interface{}) interface{}  {
	e 	:= jsoniter.UnmarshalFromString(str, v)
	if e != nil {
		return ""
	}
	return v
}

func JsonRespond(code int, message string, data interface{}, c *gin.Context) {
	c.JSON(code, gin.H{
		"code"	: code,
		"message": message,
		"data"   : data,
	})
}
