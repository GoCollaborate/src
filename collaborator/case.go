package collaborator

import (
	"encoding/json"
	"github.com/GoCollaborate/artifacts/card"
	"github.com/GoCollaborate/artifacts/digest"
	"github.com/GoCollaborate/artifacts/iremote"
	"github.com/GoCollaborate/artifacts/message"
	"github.com/GoCollaborate/constants"
	"github.com/GoCollaborate/logger"
	"github.com/GoCollaborate/store"
	"github.com/GoCollaborate/wrappers/messageHelper"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var mu = sync.Mutex{}

type Case struct {
	CaseID string `json:"caseid,omitempty"`
	*Exposed
	*Reserved
}

type Exposed struct {
	Cards     map[string]card.Card `json:"cards,omitempty"`
	TimeStamp int64                `json:"timestamp,omitempty"`
}

type Reserved struct {
	// local is the Card of localhost
	Local       card.Card `json:"local,omitempty"`
	Coordinator card.Card `json:"coordinator,omitempty"`
}

func (c *Case) readStream() error {
	bytes, err := ioutil.ReadFile(constants.DefaultCasePath)
	if err != nil {
		panic(err)
	}
	// unmarshal, overwrite default if already existed in config file
	if err := json.Unmarshal(bytes, &c); err != nil {
		logger.LogError(err.Error())
		return err
	}
	return nil
}

func (c *Case) writeStream() error {
	mu.Lock()
	defer mu.Unlock()
	mal, err := json.Marshal(&c)
	err = ioutil.WriteFile(constants.DefaultCasePath, mal, os.ModeExclusive)
	return err
}

func (c *Case) Action() {
	go func() {
		for {
			select {
			case future := <-store.GetMsgChan():
				defer future.Close()
				out, err := c.HandleMessage(future.Receive())
				if err != nil {
					logger.LogError(err)
				}
				future.Return(out)
			default:
				continue
			}
		}
	}()
}

func (c *Case) Stamp() *Case {
	c.TimeStamp = time.Now().Unix()
	return c
}

func (c *Case) Cluster() string {
	return c.CaseID
}

func (c *Case) Digest() iremote.IDigest {
	return &digest.Digest{c.Cards, c.TimeStamp}
}

func (c *Case) Update(dgst iremote.IDigest) {
	c.Cards = dgst.Cards()
	c.TimeStamp = dgst.TimeStamp()
}

func (c *Case) Terminate(key string) *Case {
	mu.Lock()
	defer mu.Unlock()
	delete(c.Cards, key)
	return c
}

func (c *Case) ReturnByPos(pos int) card.Card {
	mu.Lock()
	defer mu.Unlock()
	if l := len(c.Cards); pos > l {
		pos = pos % l
	}
	counter := 0
	for _, a := range c.Cards {
		if counter == pos {
			return a
		}
		counter++
	}
	return card.Card{}
}

func (c *Case) HandleMessage(in *message.CardMessage) (*message.CardMessage, error) {
	// return if message is wrongly sent
	var (
		out *message.CardMessage = new(message.CardMessage)
		err error                = nil
	)

	if err = c.Validate(in, out); err != nil {
		return out, err
	}
	var (
		// local digest
		ldgst = c.Digest()
		// remote digest
		rdgst = in.Digest()
		// feedback digest
		fbdgst = ldgst
	)
	switch in.Type() {
	case iremote.MsgTypeSync:
		// msg has a more recent timestamp
		if messageHelper.Compare(ldgst, rdgst) {
			fbdgst = messageHelper.Merge(ldgst, rdgst)
			// update digest to local
			c.Update(fbdgst)
		}
		// update digest to feedback
		out.Update(fbdgst)
		// return ack message
		out.SetType(iremote.MsgTypeAck)
		out.SetStatus(constants.GossipHeaderOK)
	case iremote.MsgTypeAck:
		// msg has a more recent timestamp
		if messageHelper.Compare(ldgst, rdgst) {
			fbdgst = messageHelper.Merge(ldgst, rdgst)
			// update digest to local
			c.Update(fbdgst)
		}
		// return ack message
		out.SetType(iremote.MsgTypeAck2)
		out.SetStatus(constants.GossipHeaderOK)
	case iremote.MsgTypeAck2:
		// return ack message
		out.SetType(iremote.MsgTypeAck3)
		out.SetStatus(constants.GossipHeaderOK)
	case iremote.MsgTypeAck3:
		// do nothing
	default:
		out.SetStatus(constants.GossipHeaderUnknownMsgType)
		err = constants.ErrUnknownMsgType
	}
	out.SetTo(in.From())
	out.SetFrom(in.To())
	out.SetCluster(c.Cluster())
	return out, nil
}

func (c *Case) Validate(in *message.CardMessage, out *message.CardMessage) error {
	if c.Cluster() != in.Cluster() {
		out.SetStatus(constants.GossipHeaderCaseMismatch)
		return constants.ErrCaseMismatch
	}
	if to := in.To(); !c.Local.IsEqualTo(&to) {
		logger.LogError(c)
		out.SetStatus(constants.GossipHeaderCollaboratorMismatch)
		return constants.ErrCollaboratorMismatch
	}
	return nil
}
