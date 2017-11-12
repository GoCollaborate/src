package service

import (
	"github.com/GoCollaborate/src/artifacts/card"
)

type Query struct {
	Agent card.Card `json:"card"`
}
