package task

import (
	"github.com/GoCollaborate/src/constants"
)

type TaskFuture struct {
	in  *Task
	out chan bool
	err error
}

func NewTaskFuture(in *Task) *TaskFuture {
	return &TaskFuture{in, make(chan bool), nil}
}

func (tf *TaskFuture) Receive() *Task {
	return tf.in
}

func (tf *TaskFuture) Return(out bool) {
	tf.out <- out
}

func (tf *TaskFuture) Close() {
	if len(tf.out) > 0 {
		tf.err = constants.ErrTaskChannelDirty
	}
	close(tf.out)
}

func (tf *TaskFuture) IsDone() chan bool {
	return tf.out
}

func (tf *TaskFuture) Done() bool {
	return <-tf.out
}

func (tf *TaskFuture) Error() error {
	return tf.err
}
