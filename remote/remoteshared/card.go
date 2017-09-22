package remoteshared

import (
	"github.com/GoCollaborate/utils"
	"strconv"
)

// Card is the network config of server
type Card struct {
	IP    string `json:"ip"`
	Port  int    `json:"port"`
	Alive bool   `json:"alive"`
	API   string `json:"api,omitempty"`
}

func UnmarshalCards(original []interface{}) []Card {
	var Cards []Card
	for _, o := range original {
		oo := o.(map[string]interface{})
		if oo["api"] != nil {
			Cards = append(Cards, Card{oo["ip"].(string), int(oo["port"].(float64)), oo["alive"].(bool), oo["api"].(string)})
			continue
		}
		Cards = append(Cards, Card{oo["ip"].(string), int(oo["port"].(float64)), oo["alive"].(bool), ""})
	}
	return Cards
}

func (c *Card) GetFullIP() string {
	return c.IP + ":" + strconv.Itoa(c.Port)
}

func (c *Card) GetFullExposureAddress() string {
	return utils.MapToExposureAddress(c.IP) + ":" + strconv.Itoa(c.Port)
}

func (c *Card) GetFullExposureCard() Card {
	return Card{utils.MapToExposureAddress(c.IP), c.Port, c.Alive, c.API}
}

func (c *Card) GetFullEndPoint() string {
	return c.IP + ":" + strconv.Itoa(c.Port) + "/" + c.API
}

func (c *Card) IsEqualTo(another *Card) bool {
	return c.GetFullIP() == another.GetFullIP() || c.GetFullExposureAddress() == another.GetFullExposureAddress()
}

// current RPC port
func Default() *Card {
	return &Card{utils.GetLocalIP(), utils.GetPort(), true, ""}
}
