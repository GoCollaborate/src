package taskHelper

import (
	"github.com/GoCollaborate/artifacts/task"
)

func Filter(inmaps map[int]*task.Task, f func(int, *task.Task) bool) map[int]*task.Task {
	var (
		outmaps map[int]*task.Task
	)

	for key, val := range inmaps {
		if f(key, val) {
			outmaps[key] = val
		}
	}

	return outmaps
}
