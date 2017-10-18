package service

import (
	"github.com/GoCollaborate/artifacts/card"
)

type Heartbeat struct {
	Agent card.Card `json:"card"`
}
