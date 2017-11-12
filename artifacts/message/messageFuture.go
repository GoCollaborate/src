package message

import (
	"github.com/GoCollaborate/src/constants"
)

type CardMessageFuture struct {
	in  *CardMessage
	out chan *CardMessage
	err error
}

func NewCardMessageFuture(in *CardMessage) *CardMessageFuture {
	return &CardMessageFuture{in, make(chan *CardMessage), nil}
}

func (cf *CardMessageFuture) Receive() *CardMessage {
	return cf.in
}

func (cf *CardMessageFuture) Return(out *CardMessage) {
	cf.out <- out
}

func (cf *CardMessageFuture) Close() {
	if len(cf.out) > 0 {
		cf.err = constants.ErrMessageChannelDirty
	}
	close(cf.out)
}

func (cf *CardMessageFuture) Done() *CardMessage {
	return <-cf.out
}

func (cf *CardMessageFuture) Error() error {
	return cf.err
}
