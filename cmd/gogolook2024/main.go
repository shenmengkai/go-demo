package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/shenmengkai/gogolook2024/pkg/gredis"
	"github.com/shenmengkai/gogolook2024/pkg/logging"
	"github.com/shenmengkai/gogolook2024/pkg/setting"

	"github.com/shenmengkai/gogolook2024/internal/repo"
	"github.com/shenmengkai/gogolook2024/internal/router"
)

func init() {
	setting.Setup()
	logging.Setup()
	gredis.Setup()
	task_repo.Setup()
}

func main() {
	gin.SetMode(setting.ServerSetting.RunMode)

	routersInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info] start http server listening %s", endPoint)

	server.ListenAndServe()
}
