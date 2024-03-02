package e

const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400

	ERROR_EXIST_TASK       = 10001
	ERROR_EXIST_TASK_FAIL  = 10002
	ERROR_NOT_EXIST_TASK   = 10003
	ERROR_LIST_TASKS_FAIL  = 10004
	ERROR_CREATE_TASK_FAIL = 10006
	ERROR_UPDATE_TASK_FAIL = 10007
	ERROR_DELETE_TASK_FAIL = 10008
)

type Error struct {
	Code int
}

func (err Error) Error() string {
	return GetMsg(err.Code)
}
