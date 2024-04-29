package task_repo

import (
	"encoding/json"

	"github.com/shenmengkai/go-demo/internal/models"
	"github.com/shenmengkai/go-demo/pkg/gredis"
)

type TaskRepo interface {
	ListTask() ([]models.Task, error)
	CreateTask(name string) (models.Task, error)
	UpdateTask(id int, newTag models.Task) error
	DeleteTask(id int) error
}

func Setup() {
	_, err := fetch()
	if err != nil {
		gredis.Delete("list")
		gredis.Set("list", []models.Task{})
	}
}

type TaskRepoImpl struct{}

func (repo *TaskRepoImpl) ListTask() ([]models.Task, error) {
	return fetch()
}

func (repo *TaskRepoImpl) CreateTask(name string) (models.Task, error) {
	newTask := models.Task{
		Text: name,
	}
	list, err := repo.ListTask()
	if err != nil {
		return newTask, err
	}
	newId, err := nextId()
	if err != nil {
		return newTask, err
	}
	newTask.ID = newId
	list = append(list, newTask)
	save(list)

	return newTask, nil
}

func (repo *TaskRepoImpl) UpdateTask(id int, newTag models.Task) error {
	list, err := repo.ListTask()
	if err != nil {
		return err
	}

	for i, tag := range list {
		if tag.ID != id {
			continue
		}
		list[i] = newTag
		save(list)
		return nil
	}

	return nil
}

func (repo *TaskRepoImpl) DeleteTask(id int) error {
	list, err := repo.ListTask()
	if err != nil {
		return err
	}

	newList := []models.Task{}
	for _, tag := range list {
		if tag.ID != id {
			newList = append(newList, tag)
		}
	}
	if len(newList) != len(list) {
		save(newList)
	}

	return nil
}

func fetch() ([]models.Task, error) {
	var (
		list []models.Task
	)
	data, err := gredis.Get("list")
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &list)

	return list, nil
}

func save(list []models.Task) {
	gredis.Set("list", list)
}

func nextId() (int, error) {
	return gredis.NextInt("lastId")
}
