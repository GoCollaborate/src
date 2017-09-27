package collaborator

import (
	"encoding/json"
	"github.com/GoCollaborate/constants"
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/remote/remoteshared"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var mu = sync.Mutex{}

type Case struct {
	CaseID string `json:"caseid,omitempty"`
	Exposed
	Reserved
}

type Exposed struct {
	Cards     map[string]remoteshared.Card `json:"cards,omitempty"`
	TimeStamp int64                        `json:"timestamp,omitempty"`
}

type Reserved struct {
	// local is the local Card representation
	Local       remoteshared.Card `json:"local,omitempty"`
	Coordinator remoteshared.Card `json:"coordinator,omitempty"`
}

func populate(cb *Case) error {
	bytes, err := ioutil.ReadFile(constants.DefaultCasePath)
	if err != nil {
		panic(err)
	}
	// unmarshal, overwrite default if already existed in config file
	if err := json.Unmarshal(bytes, &cb); err != nil {
		logger.LogError(err.Error())
		return err
	}
	return nil
}

func (c *Case) Disconnect() {
	for _, e := range c.Cards() {
		if e.Alive && (!e.IsEqualTo(&c.Local)) {
			client, err := launchClient(e.IP, e.Port)
			if err != nil {
				logger.LogWarning("Connection failed while disconnecting")
				continue
			}

			from := c.Local.GetFullExposureCard()
			to := e
			var in *remoteshared.CardMessage = remoteshared.NewCardMessageWithOptions(c.Cluster(), from, to, c.Cards(), c.TimeStamp())
			var out *remoteshared.CardMessage = remoteshared.NewCardMessage()
			err = client.Disconnect(in, out)

			if err != nil {
				logger.LogWarning("Calling method failed while disconnecting")
				continue
			}
			if Compare(c, out) {
				cs, _ := Merge(c, out)
				cs.writeStream()
			}
		}
	}
}

func RemoteLoad(in *remoteshared.CardMessage, out *remoteshared.CardMessage) error {
	var (
		local *Case = new(Case)
		err   error
	)

	err = populate(local)

	if err != nil {
		return err
	}

	err = local.Validate(in, out)
	if err != nil {
		return err
	}

	var (
		update bool              = false
		from   remoteshared.Card = in.From()
	)

	update = Compare(local, in)

	cp := local.Cards()

	if h, ok := cp[from.GetFullIP()]; !ok || !h.IsEqualTo(&from) {
		local.Exposed.Cards[from.GetFullIP()] = from
		local.Stamp()
		update = true
	}

	if update {
		local.writeStream()
	}

	// update local config to remote call
	out.SetCards(local.Cards())
	out.SetTimeStamp(local.TimeStamp())
	out.SetStatus(constants.GossipHeaderOK)

	return nil
}

func RemoteDisconnect(in *remoteshared.CardMessage, out *remoteshared.CardMessage) error {
	var (
		local *Case = new(Case)
		err   error
	)

	err = populate(local)
	if err != nil {
		return err
	}

	err = local.Validate(in, out)
	if err != nil {
		return err
	}

	var (
		from   remoteshared.Card = in.From()
		key                      = from.GetFullIP()
		update bool              = false
	)

	update = Compare(local, in)

	cp := local.Cards()

	if h, ok := cp[key]; ok && h.IsEqualTo(&from) {
		local.Terminate(key)
		local.Stamp()
		Merge(local, in)
		update = true
	}

	if update {
		local.writeStream()
	}

	// update local config to remote call
	out.SetCards(local.Cards())
	out.SetTimeStamp(local.TimeStamp())
	out.SetStatus(constants.GossipHeaderOK)

	return nil
}

func (c *Case) writeStream() {
	mu.Lock()
	defer mu.Unlock()
	mal, err := json.Marshal(&c)
	err = ioutil.WriteFile(constants.DefaultCasePath, mal, os.ModeExclusive)
	if err != nil {
		logger.LogError(err)
	}
}

func (c *Case) Stamp() *Case {
	c.Exposed.TimeStamp = time.Now().Unix()
	return c
}

func (c *Case) Cards() map[string]remoteshared.Card {
	return c.Exposed.Cards
}

func (c *Case) Cluster() string {
	return c.CaseID
}

func (c *Case) TimeStamp() int64 {
	return c.Exposed.TimeStamp
}

func (c *Case) Terminate(key string) *Case {
	mu.Lock()
	defer mu.Unlock()
	delete(c.Exposed.Cards, key)
	return c
}

func (c *Case) ReturnByPos(pos int) remoteshared.Card {
	mu.Lock()
	defer mu.Unlock()
	if l := len(c.Exposed.Cards); pos > l {
		pos = pos % l
	}
	counter := 0
	for _, a := range c.Exposed.Cards {
		if counter == pos {
			return a
		}
		counter++
	}
	return remoteshared.Card{}
}

func Compare(a remoteshared.ICardMessage, b remoteshared.ICardMessage) bool {
	if a.TimeStamp() < b.TimeStamp() {
		return true
	}
	return false
}

func Merge(local *Case, remote *remoteshared.CardMessage) (*Case, *remoteshared.CardMessage) {
	if local.TimeStamp() < remote.TimeStamp() {
		local.Exposed.Cards = remote.Cards()
		local.Exposed.TimeStamp = remote.TimeStamp()
		return local, remote
	}
	remote.SetCards(local.Exposed.Cards)
	remote.SetTimeStamp(local.Exposed.TimeStamp)
	return local, remote
}

func (c *Case) Validate(in *remoteshared.CardMessage, out *remoteshared.CardMessage) error {
	if c.Cluster() != in.Cluster() {
		out.SetStatus(constants.GossipHeaderCaseMismatch)
		return constants.ErrCaseMismatch
	}
	if to := in.To(); !c.Reserved.Local.IsEqualTo(&to) {
		logger.LogError(c)
		out.SetStatus(constants.GossipHeaderCollaboratorMismatch)
		return constants.ErrCollaboratorMismatch
	}
	return nil
}
