package utils

import (
	"sort"
)

func SortArrayInt(origin []int) {
	sort.SliceStable(origin, func(i, j int) bool {
		return origin[i] < origin[j]
	})
}

func SortArrayIntReverse(origin []int) {
	sort.SliceStable(origin, func(i, j int) bool {
		return origin[i] > origin[j]
	})
}
