package iremote

import (
	"github.com/GoCollaborate/src/artifacts/card"
)

type IDigest interface {
	GetCards() map[string]*card.Card
	GetTimeStamp() int64
	SetCards(map[string]*card.Card)
	SetTimeStamp(int64)
}
