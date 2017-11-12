package worker

import (
	"fmt"
	"github.com/GoCollaborate/src/artifacts/task"
	"github.com/GoCollaborate/src/logger"
	"github.com/GoCollaborate/src/store"
)

type Worker struct {
	ID          uint
	Alive       bool
	BaseTasks   chan *task.TaskFuture
	LowTasks    chan *task.TaskFuture
	MediumTasks chan *task.TaskFuture
	HighTasks   chan *task.TaskFuture
	UrgentTasks chan *task.TaskFuture
	Exit        chan bool
}

func (w *Worker) Start() {
	fs := store.GetInstance()
	go func() {
		for {
			select {
			case <-w.Exit:
				return
			case tkf := <-w.UrgentTasks:
				tk := tkf.Receive()
				logger.LogNormal(fmt.Sprintf("Worker%v:, Task Level:%v", w.ID, tk.Priority))
				logger.GetLoggerInstance().LogNormal(fmt.Sprintf("Worker%v:, Task Level:%v", w.ID, tk.Priority))
				tkf.Return(fs.Call((*tk).Consumable, &(*tk).Source, &(*tk).Result, (*tk).Context))
			default:
				select {
				case tkf := <-w.HighTasks:
					tk := tkf.Receive()
					logger.LogNormal(fmt.Sprintf("Worker%v:, Task Level:%v", w.ID, tk.Priority))
					logger.GetLoggerInstance().LogNormal(fmt.Sprintf("Worker%v:, Task Level:%v", w.ID, tk.Priority))
					tkf.Return(fs.Call((*tk).Consumable, &(*tk).Source, &(*tk).Result, (*tk).Context))
				default:
					select {
					case tkf := <-w.MediumTasks:
						tk := tkf.Receive()
						logger.LogNormal(fmt.Sprintf("Worker%v:, Task Level:%v", w.ID, tk.Priority))
						logger.GetLoggerInstance().LogNormal(fmt.Sprintf("Worker%v:, Task Level:%v", w.ID, tk.Priority))
						tkf.Return(fs.Call((*tk).Consumable, &(*tk).Source, &(*tk).Result, (*tk).Context))
					default:
						select {
						case tkf := <-w.LowTasks:
							tk := tkf.Receive()
							logger.LogNormal(fmt.Sprintf("Worker%v:, Task Level:%v", w.ID, tk.Priority))
							logger.GetLoggerInstance().LogNormal(fmt.Sprintf("Worker%v:, Task Level:%v", w.ID, tk.Priority))
							tkf.Return(fs.Call((*tk).Consumable, &(*tk).Source, &(*tk).Result, (*tk).Context))
						default:
							select {
							case tkf := <-w.BaseTasks:
								tk := tkf.Receive()
								logger.LogNormal(fmt.Sprintf("Worker%v:, Task Level:%v", w.ID, tk.Priority))
								logger.GetLoggerInstance().LogNormal(fmt.Sprintf("Worker%v:, Task Level:%v", w.ID, tk.Priority))
								tkf.Return(fs.Call((*tk).Consumable, &(*tk).Source, &(*tk).Result, (*tk).Context))
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

func (w *Worker) GetID() uint {
	return w.ID
}

func (w *Worker) Quit() {
	w.Exit <- true
}
