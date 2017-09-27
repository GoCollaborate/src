package taskutils

import (
	"github.com/GoCollaborate/server/task"
	"sort"
)

func Keys(maps map[int]*task.Task) []int {
	keys := make([]int, 0)
	for key := range maps {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	return keys
}

func KeysReverseOrder(maps map[int]*task.Task) []int {
	keys := make([]int, 0)
	for key := range maps {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})
	return keys
}
