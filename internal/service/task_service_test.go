package task_service

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/shenmengkai/gogolook2024/internal/models"
)

func TestTaskService_List(t *testing.T) {
	repoMock := new(TaskRepoImplMock)
	repoMock.On("ListTask").
		Return([]models.Task{
			{ID: 0, Text: "Text1", Status: 0},
		}, nil)
	taskService := TaskSerivceImpl{
		TaskRepo: repoMock,
	}
	list, err := taskService.List()

	repoMock.AssertCalled(t, "ListTask")
	require.Equal(t, err, nil, "TaskService.List() has error")
	require.NotEqual(t, list, nil, "TaskService.List() got empty")
	require.Equal(t, len(list), 1, "TaskService.List() got wrong length on List(), actual: %d", len(list))
}

func TestTaskService_Create(t *testing.T) {
	text := "TEST_CREATE"
	repoMock := new(TaskRepoImplMock)
	repoMock.On("CreateTask", text).
		Return(models.Task{
			ID:     0,
			Text:   text,
			Status: 0,
		}, nil)
	taskService := TaskSerivceImpl{
		TaskRepo: repoMock,
	}
	result, err := taskService.Create(text)

	repoMock.AssertCalled(t, "CreateTask", text)
	require.Equal(t, err, nil, "TaskService.Create() has error")
	require.NotEqual(t, result, nil, "TaskService.Create() got empty")
	require.Equal(t, result.Text, "TEST_CREATE", "TaskService.Create() got wrong length on List(), actual: %s", result.Text)
}

func TestTaskService_UpdateOnText(t *testing.T) {
	repoMock := new(TaskRepoImplMock)
	id := 0
	newText := "TEST_NEW_TEXT"
	mockTask := models.Task{
		Text: newText,
	}
	repoMock.On("UpdateTask", id, mockTask).
		Return(nil)
	taskService := TaskSerivceImpl{
		TaskRepo: repoMock,
	}
	err := taskService.Update(mockTask)

	repoMock.AssertCalled(t, "UpdateTask", id, mockTask)
	require.Equal(t, err, nil, "TaskService.Update() on Text has error")
}

func TestTaskService_UpdateOnStatus(t *testing.T) {
	repoMock := new(TaskRepoImplMock)
	id := 0
	newStatus := 1
	mockTask := models.Task{
		Status: newStatus,
	}
	repoMock.On("UpdateTask", id, mockTask).
		Return(nil)
	taskService := TaskSerivceImpl{
		TaskRepo: repoMock,
	}
	err := taskService.Update(mockTask)

	repoMock.AssertCalled(t, "UpdateTask", id, mockTask)
	require.Equal(t, err, nil, "TaskService.Update() on Status has error")
}

func TestTaskService_Delete(t *testing.T) {
	repoMock := new(TaskRepoImplMock)
	id := 0
	repoMock.On("DeleteTask", id).
		Return(nil)
	taskService := TaskSerivceImpl{
		TaskRepo: repoMock,
	}
	err := taskService.Delete(id)

	repoMock.AssertCalled(t, "DeleteTask", id)
	require.Equal(t, err, nil, "TaskService.Delete() has error")
}

type TaskRepoImplMock struct {
	mock.Mock
}

func (m *TaskRepoImplMock) ListTask() ([]models.Task, error) {
	args := m.Called()
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *TaskRepoImplMock) CreateTask(text string) (models.Task, error) {
	args := m.Called(text)
	return args.Get(0).(models.Task), args.Error(1)
}

func (m *TaskRepoImplMock) UpdateTask(id int, newTask models.Task) error {
	args := m.Called(id, newTask)
	return args.Error(0)
}

func (m *TaskRepoImplMock) DeleteTask(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
