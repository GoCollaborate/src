package service

import (
	"github.com/GoCollaborate/artifacts/card"
)

type Registry struct {
	Cards []card.Card `json:"cards"`
}
