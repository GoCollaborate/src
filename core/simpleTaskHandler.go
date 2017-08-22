package core

import (
	"fmt"
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
