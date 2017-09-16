package servershared

import (
	"fmt"
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/server/task"
	"github.com/GoCollaborate/store"
)

type Worker struct {
	ID          uint64
	Alive       bool
	BaseTasks   chan *task.Task
	LowTasks    chan *task.Task
	MediumTasks chan *task.Task
	HighTasks   chan *task.Task
	UrgentTasks chan *task.Task
	Exit        chan bool
}

func (w *Worker) Start() {
	fs := store.GetInstance()
	go func() {
		for {
			select {
			case <-w.Exit:
				return
			case tk := <-w.UrgentTasks:
				logger.LogNormal(fmt.Sprintf("Worker%v:, Task Level:%v", w.ID, tk.Priority))
				fs.Call((*tk).Consumable, &(*tk).Source, &(*tk).Result, (*tk).Context)
			default:
				select {
				case tk := <-w.HighTasks:
					logger.LogNormal(fmt.Sprintf("Worker%v:, Task Level:%v", w.ID, tk.Priority))
					fs.Call((*tk).Consumable, &(*tk).Source, &(*tk).Result, (*tk).Context)
				default:
					select {
					case tk := <-w.MediumTasks:
						logger.LogNormal(fmt.Sprintf("Worker%v:, Task Level:%v", w.ID, tk.Priority))
						fs.Call((*tk).Consumable, &(*tk).Source, &(*tk).Result, (*tk).Context)
					default:
						select {
						case tk := <-w.LowTasks:
							logger.LogNormal(fmt.Sprintf("Worker%v:, Task Level:%v", w.ID, tk.Priority))
							fs.Call((*tk).Consumable, &(*tk).Source, &(*tk).Result, (*tk).Context)
						default:
							select {
							case tk := <-w.BaseTasks:
								logger.LogNormal(fmt.Sprintf("Worker%v:, Task Level:%v", w.ID, tk.Priority))
								fs.Call((*tk).Consumable, &(*tk).Source, &(*tk).Result, (*tk).Context)
							default:
								continue
							}
						}
					}
				}
			}
		}
	}()
}

func (w *Worker) GetID() uint64 {
	return w.ID
}

func (w *Worker) Quit() {
	w.Exit <- true
}
