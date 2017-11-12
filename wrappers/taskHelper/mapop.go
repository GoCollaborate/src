package taskHelper

import (
	"github.com/GoCollaborate/src/artifacts/task"
)

type mapop struct {
	entries [][]int
	inmaps  map[int]*task.Task
}

func Map(ins map[int]*task.Task, ens ...[]int) *mapop {
	return &mapop{
		entries: ens,
		inmaps:  ins,
	}
}

// reduce the result from input set to the source of output set
func (mop *mapop) ReduceTo(exs []int) map[int]*task.Task {
	var (
		length  int                = len(exs)
		outmaps map[int]*task.Task = map[int]*task.Task{}
	)

	if lens := len(mop.entries); length > lens {
		length = lens
	}

	for i := 0; i < length; i++ {
		ens := mop.entries[i]
		outmaps[exs[i]] = mop.inmaps[exs[i]]

		for _, e := range ens {
			outmaps[exs[i]].Source = append(outmaps[exs[i]].Source, mop.inmaps[e].Source)
		}
	}

	return outmaps
}
