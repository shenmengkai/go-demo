package e

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "fail",
	INVALID_PARAMS: "invalid parameters",

	ERROR_EXIST_TASK:       "task already exists",
	ERROR_EXIST_TASK_FAIL:  "failed to get existing task",
	ERROR_NOT_EXIST_TASK:   "tag does not exist",
	ERROR_LIST_TASKS_FAIL:  "Failed to get all tasks",
	ERROR_CREATE_TASK_FAIL: "failed to create task",
	ERROR_UPDATE_TASK_FAIL: "failed to update task",
	ERROR_DELETE_TASK_FAIL: "failed to delete task",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
