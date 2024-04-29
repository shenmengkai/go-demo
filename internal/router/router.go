package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	_ "github.com/shenmengkai/go-demo/docs"
	"github.com/shenmengkai/go-demo/internal/middleware"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/tasks", mdw().ListTasks)
	r.POST("/task", mdw().CreateTask)
	r.PUT("/task/:id", mdw().UpdateTask)
	r.DELETE("/task/:id", mdw().DeleteTask)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}

func mdw() middleware.TaskMiddleware {
	var handler middleware.TaskMiddleware = &middleware.TaskMiddlewareImpl{
		TaskService: middleware.GetTaskService(),
	}
	return handler
}
