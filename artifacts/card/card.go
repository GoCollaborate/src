package card

import (
	"github.com/GoCollaborate/src/utils"
	"strconv"
)

// Card is the network config of server
func NewCard(ip string, port int32, alive bool, api string, seed bool) *Card {
	return &Card{
		IP:    ip,
		Port:  port,
		Alive: alive,
		API:   api,
		Seed:  seed,
	}
}

func (c *Card) IsInitialized() bool {
	if len(c.IP) > 0 {
		return true
	}
	return false
}

func (c *Card) GetFullIP() string {
	return c.IP + ":" + strconv.Itoa(int(c.Port))
}

func (c *Card) GetFullExposureAddress() string {
	return utils.MapToExposureAddress(c.IP) + ":" + strconv.Itoa(int(c.Port))
}

func (c *Card) GetFullExposureCard() Card {
	return Card{utils.MapToExposureAddress(c.IP), c.Port, c.Alive, c.API, c.Seed}
}

func (c *Card) GetFullEndPoint() string {
	return c.IP + ":" + strconv.Itoa(int(c.Port)) + "/" + c.API
}

func (c *Card) IsEqualTo(another *Card) bool {
	return c.GetFullIP() == another.GetFullIP() || c.GetFullExposureAddress() == another.GetFullExposureAddress()
}

func (c *Card) IsSeed() bool {
	return c.Seed
}

func (c *Card) ToSeed() {
	c.Seed = true
}

func (c *Card) SetAlive(alive bool) {
	c.Alive = alive
}

// current RPC port
func Default() *Card {
	return &Card{utils.GetLocalIP(), utils.GetPort(), true, "", true}
}
