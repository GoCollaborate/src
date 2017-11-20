package message

import (
	"github.com/GoCollaborate/src/artifacts/card"
	"github.com/GoCollaborate/src/artifacts/digest"
	"github.com/GoCollaborate/src/artifacts/iremote"
	"github.com/GoCollaborate/src/constants"
	"time"
)

func NewCardMessage() *CardMessage {
	cm := new(CardMessage)
	cm.Digest = new(digest.Digest)
	return cm
}

func NewCardMessageWithOptions(cluster string, from *card.Card, to *card.Card, cards map[string]*card.Card, timestamp int64, msgType CardMessage_Type) *CardMessage {
	message := new(CardMessage)
	message.Digest = new(digest.Digest)
	message.SetCluster(cluster).SetFrom(from).SetTo(to).SetCards(cards).SetTimeStamp(timestamp).SetStatus(constants.GossipHeaderUnknownError).SetType(msgType)
	return message
}

func (cm *CardMessage) SetCluster(cluster string) *CardMessage {
	cm.Cluster = cluster
	return cm
}

func (cm *CardMessage) SetFrom(from *card.Card) *CardMessage {
	f := *from
	cm.From = &f
	return cm
}

func (cm *CardMessage) SetTo(to *card.Card) *CardMessage {
	t := *to
	cm.To = &t
	return cm
}

func (cm *CardMessage) SetCards(cards map[string]*card.Card) *CardMessage {
	cm.Digest.Cards = cards
	return cm
}

func (cm *CardMessage) SetTimeStamp(timestamp int64) *CardMessage {
	cm.Digest.Ts = timestamp
	return cm
}

func (cm *CardMessage) SetStatus(status constants.Header) *CardMessage {
	cm.Status = &Status{status.Key, status.Value}
	return cm
}

func (cm *CardMessage) SetType(msgType CardMessage_Type) *CardMessage {
	cm.Type = msgType
	return cm
}

func (cm *CardMessage) DeleteDigestCard(key string) {
	delete(cm.Digest.GetCards(), key)
}

func (cm *CardMessage) Update(dgst iremote.IDigest) {
	cm.SetCards(dgst.GetCards())
	cm.SetTimeStamp(dgst.GetTimeStamp())
}

func (cm *CardMessage) Stamp() {
	cm.SetTimeStamp(time.Now().Unix())
}
