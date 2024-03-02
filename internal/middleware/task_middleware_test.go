package middleware

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/shenmengkai/gogolook2024/internal/models"
)

func TestTaskMdw_ListTasks(t *testing.T) {
	svcMock := new(TaskServiceImplMock)
	svcMock.On("List").
		Return([]models.Task{}, nil)

	ctx, resp := mockContext(t, http.MethodGet, "/tasks", "")
	testMdw(svcMock).ListTasks(ctx)

	require.Equal(t,
		http.StatusOK,
		resp.Code,
		"TaskMiddleware.ListTasks() not OK")
	require.Equal(t,
		`{"result":[]}`,
		resp.Body.String(),
		"TaskMiddleware.ListTasks() body is incorrect")
}

func TestTaskMdw_CreateTask(t *testing.T) {
	text := "TEST_CREATE_TASK"
	svcMock := new(TaskServiceImplMock)
	svcMock.On("Create", text).
		Return(models.Task{
			ID:     1,
			Text:   text,
			Status: 0,
		}, nil)

	ctx, resp := mockContext(t, http.MethodPost, "/task", fmt.Sprintf(`{ "text": "%s" }`, text))
	testMdw(svcMock).CreateTask(ctx)

	require.Equal(t,
		http.StatusCreated,
		resp.Code,
		"TaskMiddleware.CreateTask() not OK")
	require.Equal(t,
		fmt.Sprintf(`{"result":{"id":1,"text":"%s","status":0}}`, text),
		resp.Body.String(),
		"TaskMiddleware.CreateTask() body is incorrect")
}

func TestTaskMdw_UpdateTask(t *testing.T) {
	text := "TEST_UPDATE_TASK"
	oldTask := models.Task{
		ID:     90,
		Text:   "OLD_TASK",
		Status: 0,
	}
	newTask := models.Task{
		ID:     90,
		Text:   text,
		Status: 0,
	}
	svcMock := new(TaskServiceImplMock)
	svcMock.On("Get", 90).
		Return(oldTask, nil)
	svcMock.On("Update", newTask).
		Return(nil)

	ctx, resp := mockContext(t, http.MethodPut, "/task/90", fmt.Sprintf(`{ "text": "%s" }`, text))
	testMdw(svcMock).UpdateTask(ctx)

	require.Equal(t,
		http.StatusOK,
		resp.Code,
		"TaskMiddleware.UpdateTask() not OK")
	require.Equal(t,
		fmt.Sprintf(`{"result":{"id":90,"text":"%s","status":0}}`, text),
		resp.Body.String(),
		"TaskMiddleware.UpdateTask() body is incorrect")
}

// func TestTaskMdw_UpdateTaskBadForm(t *testing.T) {
// 	text := "TEST_UPDATE_TASK"
// 	id := 95
// 	oldTask := models.Task{
// 		ID:     id,
// 		Text:   "OLD_TASK",
// 		Status: 0,
// 	}
// 	svcMock := new(TaskServiceImplMock)
// 	svcMock.On("Get", id).
// 		Return(oldTask, nil)

// 	ctx, resp := mockContext(t, http.MethodPut,
// 		fmt.Sprintf(`/task/%d`, id),
// 		fmt.Sprintf(`{ "badtext": "%s" }`, text))
// 	testMdw(svcMock).UpdateTask(ctx)
// 	svcMock.AssertNotCalled(t, "Update")

// 	require.Equal(t,
// 		http.StatusBadRequest,
// 		resp.Code,
// 		"TaskMiddleware.UpdateTask() not OK")
// }

func TestTaskMdw_UpdateTaskWithoutId(t *testing.T) {
	svcMock := new(TaskServiceImplMock)
	ctx, resp := mockContext(t, http.MethodPut, "/task", "")

	testMdw(svcMock).UpdateTask(ctx)
	svcMock.AssertNotCalled(t, "Get")
	svcMock.AssertNotCalled(t, "Update")

	require.Equal(t,
		http.StatusBadRequest,
		resp.Code,
		"TaskMiddleware.UpdateTask() should be NOT_FOUND when missing ID")
}

func TestTaskMdw_UpdateTaskNotExist(t *testing.T) {
	id := 94
	text := "TEST_UPDATE_TASK"
	svcMock := new(TaskServiceImplMock)
	svcMock.On("Get", id).
		Return(models.Task{}, errors.New("not exist"))

	ctx, resp := mockContext(t,
		http.MethodPut,
		fmt.Sprintf(`/task/%d`, id),
		fmt.Sprintf(`{ "text": "%s" }`, text))
	testMdw(svcMock).UpdateTask(ctx)
	svcMock.AssertNotCalled(t, "Update")

	require.Equal(t,
		http.StatusNotFound,
		resp.Code,
		"TaskMiddleware.UpdateTask() should be NOT_FOUND")
}

func TestTaskMdw_DeleteTask(t *testing.T) {
	id := 119
	svcMock := new(TaskServiceImplMock)
	svcMock.On("Get", id).
		Return(models.Task{}, nil)
	svcMock.On("Delete", id).
		Return(nil)

	ctx, resp := mockContext(t, http.MethodDelete, fmt.Sprintf(`/task/%d`, id), "")
	testMdw(svcMock).DeleteTask(ctx)

	require.Equal(t,
		http.StatusOK,
		resp.Code,
		"TaskMiddleware.DeleteTask() not OK")
	require.Equal(t,
		"",
		resp.Body.String(),
		"TaskMiddleware.DeleteTask() body is incorrect")
}

func TestTaskMdw_DeleteTaskWithoutId(t *testing.T) {
	svcMock := new(TaskServiceImplMock)
	ctx, resp := mockContext(t, http.MethodDelete, "/task", "")

	testMdw(svcMock).DeleteTask(ctx)
	svcMock.AssertNotCalled(t, "Get")
	svcMock.AssertNotCalled(t, "Delete")

	require.Equal(t,
		http.StatusBadRequest,
		resp.Code,
		"TaskMiddleware.DeleteTask() should be BAD_REQUEST when missing ID")
}
func TestTaskMdw_DeleteTaskNotExist(t *testing.T) {
	id := 140
	svcMock := new(TaskServiceImplMock)
	svcMock.On("Get", id).
		Return(models.Task{}, errors.New("not exist"))

	ctx, resp := mockContext(t, http.MethodDelete, fmt.Sprintf(`/task/%d`, id), "")
	testMdw(svcMock).DeleteTask(ctx)
	svcMock.AssertNotCalled(t, "Delete")

	require.Equal(t,
		http.StatusNotFound,
		resp.Code,
		"TaskMiddleware.DeleteTask() should be NOT_FOUND when not exists")
}

func testMdw(svcMock *TaskServiceImplMock) *TaskMiddlewareImpl {
	return &TaskMiddlewareImpl{
		TaskService: svcMock,
	}
}

func mockContext(t *testing.T, method string, path string, payload string) (*gin.Context, *httptest.ResponseRecorder) {
	setParamIdIfPresent := func(ctx *gin.Context) {
		tokens := strings.Split(path, "/")
		if len(tokens) < 1 {
			return
		}
		ctx.Params = []gin.Param{
			{Key: "id", Value: tokens[len(tokens)-1]},
		}
	}

	newRequest := func() *http.Request {
		var buf io.Reader = nil
		if len(payload) > 0 {
			buf = bytes.NewBuffer([]byte(payload))
		}
		req := httptest.NewRequest(method, path, buf)
		if len(payload) > 0 {
			req.Header.Set("Content-Type", "application/json")
		}
		return req
	}

	resp := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(resp)
	setParamIdIfPresent(ctx)
	ctx.Request = newRequest()
	return ctx, resp
}

type TaskServiceImplMock struct {
	mock.Mock
}

func (m *TaskServiceImplMock) List() ([]models.Task, error) {
	args := m.Called()
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *TaskServiceImplMock) Get(id int) (models.Task, error) {
	args := m.Called(id)
	return args.Get(0).(models.Task), args.Error(1)

}

func (m *TaskServiceImplMock) Create(text string) (models.Task, error) {
	args := m.Called(text)
	return args.Get(0).(models.Task), args.Error(1)
}

func (m *TaskServiceImplMock) Update(task models.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *TaskServiceImplMock) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
