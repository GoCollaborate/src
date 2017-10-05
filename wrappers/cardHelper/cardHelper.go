package cardHelper

import (
	"github.com/GoCollaborate/artifacts/card"
	"github.com/GoCollaborate/logger"
)

func UnmarshalCards(original []interface{}) []card.Card {
	var cards []card.Card
	for _, o := range original {
		oo := o.(map[string]interface{})
		if oo["api"] != nil {
			cards = append(cards, card.Card{oo["ip"].(string), int(oo["port"].(float64)), oo["alive"].(bool), oo["api"].(string)})
			continue
		}
		cards = append(cards, card.Card{oo["ip"].(string), int(oo["port"].(float64)), oo["alive"].(bool), ""})
	}
	return cards
}

func RangePrint(cards map[string]card.Card) {
	logger.LogNormal("Cards:")
	for _, c := range cards {
		var alive string
		if c.Alive {
			alive = "Alive"
		} else {
			alive = "Terminated"
		}
		logger.LogListPoint(c.GetFullIP(), alive)
	}
}
