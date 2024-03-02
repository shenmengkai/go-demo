package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/shenmengkai/gogolook2024/internal/middleware"
)

func mdw() middleware.TaskMiddleware {
	var handler middleware.TaskMiddleware = &middleware.TaskMiddlewareImpl{
		TaskService: middleware.GetTaskService(),
	}
	return handler
}

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/tasks", mdw().ListTasks)
	r.POST("/task", mdw().CreateTask)
	r.PUT("/task/:id", mdw().UpdateTask)
	r.DELETE("/task/:id", mdw().DeleteTask)

	return r
}
