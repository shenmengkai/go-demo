package app

import (
	"github.com/gin-gonic/gin"

	"github.com/shenmengkai/gogolook2024/pkg/e"
)

type Gin struct {
	C *gin.Context
}

type ErrorResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Response setting gin.JSON
func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	if errCode == e.SUCCESS {
		if data != nil {
			g.C.JSON(httpCode, data)
		} else {
			g.C.Status(httpCode)
		}
		return
	}
	g.C.JSON(httpCode, ErrorResponse{
		Code: errCode,
		Msg:  e.GetMsg(errCode),
		Data: data,
	})
}
