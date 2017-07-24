package server

import (
	"fmt"
)

type taskType int

const (
	SHORT taskType = iota
	LONG
	ROUTINE
	PERMANENT
)

type taskPriority int

const (
	BASE taskPriority = iota
	LOW
	MEDIUM
	HIGH
	URGENT
)

type TaskType interface {
	GetType() taskType
}

type TaskPriority interface {
	GetPriority() taskPriority
}

func (t *taskType) GetType() taskType {
	return *t
}

func (t *taskPriority) GetPriority() taskPriority {
	return *t
}

type Task struct {
	Type       taskType
	Priority   taskPriority
	Consumable func() bool
}

func TaskA() Task {
	return Task{PERMANENT, BASE, func() bool {
		fmt.Println("Task A Executed...")
		return true
	}}
}
func TaskB() Task {
	return Task{LONG, MEDIUM, func() bool {
		fmt.Println("Task B Executed...")
		return true
	}}
}
func TaskC() Task {
	return Task{SHORT, URGENT, func() bool {
		fmt.Println("Task C Executed...")
		return true
	}}
}
