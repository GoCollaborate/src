package taskHelper

import (
	"github.com/GoCollaborate/artifacts/task"
)

// slice the data source of the map into N separate segments
func Slice(inmaps map[int]*task.Task, n int) map[int]*task.Task {
	if n < 2 {
		return inmaps
	}

	var (
		gap     = len(inmaps)
		outmaps = make(map[int]*task.Task)
	)
	for k, t := range inmaps {
		var (
			sgap = len(t.Source)
			i    = 0
		)

		for ; i < n-1; i++ {
			source := t.Source[i*sgap/n : (i+1)*sgap/n]
			outmaps[(k+1)*gap+i] = &task.Task{t.Type, t.Priority, t.Consumable, source, t.Result, t.Context, t.Stage}
		}

		source := t.Source[i*sgap/n:]
		outmaps[(k+1)*gap+i] = &task.Task{t.Type, t.Priority, t.Consumable, source, t.Result, t.Context, t.Stage}
	}

	return outmaps
}
