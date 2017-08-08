package core

import (
	"fmt"
	"github.com/GoCollaborate/server/task"
)

func TaskA() task.Task {
	return task.Task{task.PERMANENT, task.BASE, func() bool {
		fmt.Println("Task A Executed...")
		return true
	}}
}
func TaskB() task.Task {
	return task.Task{task.LONG, task.MEDIUM, func() bool {
		fmt.Println("Task B Executed...")
		return true
	}}
}
func TaskC() task.Task {
	return task.Task{task.SHORT, task.URGENT, func() bool {
		fmt.Println("Task C Executed...")
		return true
	}}
}
