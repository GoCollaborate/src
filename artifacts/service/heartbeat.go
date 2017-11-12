package service

import (
	"github.com/GoCollaborate/src/artifacts/card"
)

type Heartbeat struct {
	Agent card.Card `json:"card"`
}
