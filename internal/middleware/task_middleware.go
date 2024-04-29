package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"github.com/shenmengkai/go-demo/pkg/app"
	"github.com/shenmengkai/go-demo/pkg/e"

	"github.com/shenmengkai/go-demo/internal/models"
	task_repo "github.com/shenmengkai/go-demo/internal/repo"
	task_service "github.com/shenmengkai/go-demo/internal/service"
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

type ListTaskResponse struct {
	Result []models.Task `json:"result"`
}

type TaskResponse struct {
	Result models.Task `json:"result"`
}

// ListTasks godoc
// @Summary List all tasks
// @Description Retrieves a list of tasks.
// @Tags Task
// @Accept json
// @Produce json
// @Success 200 {object} ListTaskResponse "A list of tasks"
// @Router /tasks [get]
func (m *TaskMiddlewareImpl) ListTasks(c *gin.Context) {
	appG := app.Gin{C: c}

	list, err := m.TaskService.List()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_LIST_TASKS_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, ListTaskResponse{
		Result: list,
	})
}

type CreateTaskForm struct {
	Text string `form:"text"`
}

// CreateTask godoc
// @Summary Create a new task
// @Description Adds a new task with the provided text.
// @Tags Task
// @Accept json
// @Produce json
// @Param requestBody body CreateTaskForm true "Task creation request body"
// @Success 201 {object} TaskResponse "Successfully created task with the provided text, returns task which created."
// @Router /task [post]
func (m *TaskMiddlewareImpl) CreateTask(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		form   CreateTaskForm
		result models.Task
		err    error
	)

	if err := c.ShouldBindJSON(&form); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	if len(form.Text) == 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	result, err = m.TaskService.Create(form.Text)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CREATE_TASK_FAIL, nil)
		return
	}

	appG.Response(http.StatusCreated, e.SUCCESS, TaskResponse{
		Result: result,
	})
}

type UpdateTaskForm struct {
	ID     *int    `form:"id"`
	Text   *string `form:"text"`
	Status *int    `form:"status"`
}

// UpdateTask godoc
// @Summary Update a task
// @Description Update the task with the specified ID.
// @Tags Task
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param task body UpdateTaskForm true "Fields and values to update"
// @Success 200 {object} TaskResponse "Task updated successfully"
// @Failure 400 {string} string "Missing ID or incorrect fields"
// @Failure 404 {string} string "Task of ID not found"
// @Router /task/{id} [put]
func (m *TaskMiddlewareImpl) UpdateTask(c *gin.Context) {
	var (
		appG  = app.Gin{C: c}
		idStr = c.Param("id")
		form  = UpdateTaskForm{}
	)

	if len(idStr) == 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "missing id")
		return
	}

	taskId, err := com.StrTo(idStr).Int()
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, fmt.Sprintf("incorrect id(\"%s\")", idStr))
		return
	}

	if taskId < 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, fmt.Sprintf("incorrect id(%d)", taskId))
		return
	}

	task, err := m.TaskService.Get(taskId)
	if err != nil {
		appG.Response(http.StatusNotFound, e.ERROR_EXIST_TASK_FAIL, "task not exists")
		return
	}

	if err := c.ShouldBindJSON(&form); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, "incorrect format")
		return
	}

	haveSet := false
	if form.Text != nil {
		if len(*form.Text) == 0 {
			appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
			return
		} else {
			task.Text = *form.Text
			haveSet = true
		}
	}

	if form.Status != nil {
		if *form.Status < 0 || *form.Status > 1 {
			appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
			return
		} else {
			task.Status = *form.Status
			haveSet = true
		}
	}

	if !haveSet {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	err = m.TaskService.Update(task)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_UPDATE_TASK_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, TaskResponse{
		Result: task,
	})
}

// DeleteTask godoc
// @Summary Delete a task
// @Description Deletes a task by its ID.
// @Tags Task
// @Accept json
// @Produce json
// @Param id path int true "ID of the task to delete"
// @Success 200 "Task deleted successfully"
// @Failure 400 "Missing ID"
// @Failure 403 "Task of ID not found"
// @Router /task/{id} [delete]
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
