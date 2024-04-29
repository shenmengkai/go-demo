package task_service

import (
	"github.com/shenmengkai/go-demo/pkg/e"

	"github.com/shenmengkai/go-demo/internal/models"
	"github.com/shenmengkai/go-demo/internal/repo"
)

type TaskService interface {
	List() ([]models.Task, error)
	Get(id int) (models.Task, error)
	Create(text string) (models.Task, error)
	Update(task models.Task) error
	Delete(id int) error
}

type TaskSerivceImpl struct {
	TaskRepo task_repo.TaskRepo
}

func (s *TaskSerivceImpl) List() ([]models.Task, error) {
	var (
		list []models.Task
	)

	list, err := s.TaskRepo.ListTask()
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *TaskSerivceImpl) Get(id int) (models.Task, error) {
	list, err := s.List()
	if err != nil {
		return models.Task{}, err
	}
	for _, t := range list {
		if t.ID == id {
			return t, nil
		}
	}
	return models.Task{}, e.Error{Code: e.ERROR_NOT_EXIST_TASK}
}

func (s *TaskSerivceImpl) Create(text string) (models.Task, error) {
	return s.TaskRepo.CreateTask(text)
}

func (s *TaskSerivceImpl) Update(task models.Task) error {
	return s.TaskRepo.UpdateTask(task.ID, task)
}

func (s *TaskSerivceImpl) Delete(id int) error {
	return s.TaskRepo.DeleteTask(id)
}
