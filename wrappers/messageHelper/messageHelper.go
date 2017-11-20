package messageHelper

import (
	"github.com/GoCollaborate/src/artifacts/iremote"
	"github.com/GoCollaborate/src/artifacts/message"
	"github.com/GoCollaborate/src/store"
)

func Exchange(in *message.CardMessage) (out *message.CardMessage, err error) {
	future := message.NewCardMessageFuture(in)
	// push future into message chan
	store.GetMsgChan() <- future
	// wait until future is consumed
	out = future.Done()
	return out, future.Error()
}

func Compare(a iremote.IDigest, b iremote.IDigest) bool {
	if a.GetTimeStamp() < b.GetTimeStamp() {
		return true
	}
	return false
}

func Merge(a iremote.IDigest, b iremote.IDigest) (c iremote.IDigest) {
	if a.GetTimeStamp() < b.GetTimeStamp() {
		a.SetCards(b.GetCards())
		a.SetTimeStamp(b.GetTimeStamp())
		return a
	}
	b.SetCards(a.GetCards())
	b.SetTimeStamp(a.GetTimeStamp())
	return b
}
