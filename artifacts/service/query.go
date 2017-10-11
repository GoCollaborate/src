package service

import (
	"github.com/GoCollaborate/artifacts/card"
)

type Query struct {
	Agent card.Card `json:"card"`
}
