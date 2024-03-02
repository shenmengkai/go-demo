package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"github.com/shenmengkai/gogolook2024/pkg/app"
	"github.com/shenmengkai/gogolook2024/pkg/e"

	"github.com/shenmengkai/gogolook2024/internal/models"
	"github.com/shenmengkai/gogolook2024/internal/repo"
	"github.com/shenmengkai/gogolook2024/internal/service"
)

type TaskMiddleware interface {
	ListTasks(c *gin.Context)
	CreateTask(c *gin.Context)
	UpdateTask(c *gin.Context)
	DeleteTask(c *gin.Context)
}

type TaskMiddlewareImpl struct {
	TaskService task_service.TaskService
}

func GetTaskService() task_service.TaskService {
	var service task_service.TaskService = &task_service.TaskSerivceImpl{
		TaskRepo: &task_repo.TaskRepoImpl{},
	}
	return service
}

func (m *TaskMiddlewareImpl) ListTasks(c *gin.Context) {
	appG := app.Gin{C: c}

	list, err := m.TaskService.List()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_LIST_TASKS_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, list)
}

type CreateTaskForm struct {
	Text string `form:"text"`
}

func (m *TaskMiddlewareImpl) CreateTask(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		form   CreateTaskForm
		result models.Task
		err    error
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	result, err = m.TaskService.Create(form.Text)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CREATE_TASK_FAIL, nil)
		return
	}

	appG.Response(http.StatusCreated, e.SUCCESS, result)
}

type UpdateTaskForm struct {
	Text   *string `form:"text"`
	Status *int    `form:"status"`
}

func (m *TaskMiddlewareImpl) UpdateTask(c *gin.Context) {
	var (
		appG  = app.Gin{C: c}
		idStr = c.Param("id")
		form  = UpdateTaskForm{}
	)

	if len(idStr) == 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	taskId, err := com.StrTo(idStr).Int()
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	if taskId < 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	task, err := m.TaskService.Get(taskId)
	if err != nil {
		appG.Response(http.StatusNotFound, e.ERROR_EXIST_TASK_FAIL, nil)
		return
	}

	if _, errCode := app.BindAndValid(c, &form); errCode != e.SUCCESS {
		appG.Response(http.StatusBadRequest, errCode, nil)
		return
	}

	if form.Text != nil {
		task.Text = *form.Text
	}
	if form.Status != nil {
		task.Status = *form.Status
	}

	err = m.TaskService.Update(task)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_UPDATE_TASK_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, task)
}

func (m *TaskMiddlewareImpl) DeleteTask(c *gin.Context) {
	appG := app.Gin{C: c}
	idStr := c.Param("id")
	if len(idStr) == 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	taskId, err := com.StrTo(idStr).Int()
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	if taskId < 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	if _, err := m.TaskService.Get(taskId); err != nil {
		appG.Response(http.StatusNotFound, e.ERROR_EXIST_TASK_FAIL, nil)
		return
	}

	if err := m.TaskService.Delete(taskId); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_TASK_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
