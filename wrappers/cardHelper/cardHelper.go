package cardHelper

import (
	"github.com/GoCollaborate/src/artifacts/card"
	"github.com/GoCollaborate/src/logger"
)

func UnmarshalCards(original []interface{}) []card.Card {
	var cards []card.Card
	for _, o := range original {
		oo := o.(map[string]interface{})

		var (
			api  string = ""
			seed bool   = false
		)

		if oo["api"] != nil {
			api = oo["api"].(string)
		}

		if oo["seed"] != nil {
			seed = oo["seed"].(bool)
		}

		cards = append(cards, card.Card{oo["ip"].(string), int(oo["port"].(float64)), oo["alive"].(bool), api, seed})
	}
	return cards
}

func RangePrint(cards map[string]card.Card) {
	for _, c := range cards {
		var (
			alive string
			seed  string
		)
		if c.Alive {
			alive = "Alive"
		} else {
			alive = "Terminated"
		}

		if c.IsSeed() {
			seed = "Seed"
		} else {
			seed = "Non-Seed"
		}
		logger.LogListPoint(c.GetFullIP(), alive, seed)
		logger.GetLoggerInstance().LogListPoint(c.GetFullIP(), alive, seed)
	}
}
