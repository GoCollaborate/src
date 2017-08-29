package core

import (
	"fmt"
	//"github.com/GoCollaborate/server/mapper"
	//"github.com/GoCollaborate/server/reducer"
	"github.com/GoCollaborate/server/task"
	"net/http"
)

// it is recommended that the task get registered at
// RegCenter
func TaskAHandler() (interface{}, []string) {
	return func(w http.ResponseWriter, r *http.Request) task.Task {
		return task.Task{task.PERMANENT, task.BASE, func(source []task.Countable, result []task.Countable, context *task.TaskContext) bool {
			// deal with passed in request
			fmt.Println("Task A Executed...")

			return true
		}, []task.Countable{}, []task.Countable{}, task.NewTaskContext(struct {
			A string
			B int
			C []bool
		}{})}
	}, []string{"GET", "POST"}
}
func TaskBHandler() (interface{}, []string) {
	return func(w http.ResponseWriter, r *http.Request) task.Task {
		return task.Task{task.PERMANENT, task.BASE, func(source []task.Countable, result []task.Countable, context *task.TaskContext) bool {
			// deal with passed in request
			fmt.Println("Task B Executed...")
			return true
		}, []task.Countable{}, []task.Countable{}, task.NewTaskContext(struct {
			A string
			B int
			C []bool
		}{})}
	}, []string{"GET"}
}
func TaskCHandler() (interface{}, []string) {
	return func(w http.ResponseWriter, r *http.Request) task.Task {
		return task.Task{task.PERMANENT, task.BASE, func(source []task.Countable, result []task.Countable, context *task.TaskContext) bool {
			// deal with passed in request
			fmt.Println("Task C Executed...")
			return true
		}, []task.Countable{}, []task.Countable{}, task.NewTaskContext(struct {
			A string
			B int
			C []bool
		}{})}
	}, []string{"POST"}
}

type SimpleMapper int

func (m *SimpleMapper) Map(t *task.Task) (map[int64]task.Task, error) {
	maps := make(map[int64]task.Task)
	for i, r := range t.Source {
		maps[int64(i)] = task.Task{t.Type, t.Priority, func([]task.Countable, []task.Countable, *task.TaskContext) bool { return false }, []task.Countable{r}, []task.Countable{}, nil}
	}
	return maps, nil
}

type SimpleReducer int

func (r *SimpleReducer) Reduce(sources map[int64]task.Task, result *task.Task) error {
	rs := *result
	for _, s := range sources {
		rs.Result = append(rs.Result, s.Result...)
	}
	return nil
}
