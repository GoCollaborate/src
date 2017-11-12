package service

import (
	"github.com/GoCollaborate/src/artifacts/card"
)

type Registry struct {
	Cards []card.Card `json:"cards"`
}
