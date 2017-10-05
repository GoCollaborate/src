package message

import (
	"github.com/GoCollaborate/artifacts/card"
	"github.com/GoCollaborate/artifacts/digest"
	"github.com/GoCollaborate/artifacts/iremote"
	"github.com/GoCollaborate/constants"
	"time"
)

type CardMessage struct {
	Cluster_ string           `json:"cluster,omitempty"`
	From_    card.Card        `json:"from"`
	To_      card.Card        `json:"to"`
	Status_  constants.Header `json:"status"`
	Type_    iremote.MsgType  `json:"msgtype"`
	Digest_  digest.Digest    `json:"digest"`
}

func NewCardMessage() *CardMessage {
	return new(CardMessage)
}

func NewCardMessageWithOptions(cluster string, from card.Card, to card.Card, cards map[string]card.Card, timestamp int64, msgType iremote.MsgType) *CardMessage {
	message := CardMessage{}
	message.SetCluster(cluster).SetFrom(from).SetTo(to).SetCards(cards).SetTimeStamp(timestamp).SetStatus(constants.GossipHeaderUnknownError).SetType(msgType)
	return &message
}

func (cm *CardMessage) SetCluster(cluster string) *CardMessage {
	cm.Cluster_ = cluster
	return cm
}

func (cm *CardMessage) SetFrom(from card.Card) *CardMessage {
	cm.From_ = from
	return cm
}

func (cm *CardMessage) SetTo(to card.Card) *CardMessage {
	cm.To_ = to
	return cm
}

func (cm *CardMessage) SetCards(cards map[string]card.Card) *CardMessage {
	cm.Digest_.Cards_ = cards
	return cm
}

func (cm *CardMessage) SetTimeStamp(timestamp int64) *CardMessage {
	cm.Digest_.Ts_ = timestamp
	return cm
}

func (cm *CardMessage) SetStatus(status constants.Header) *CardMessage {
	cm.Status_ = status
	return cm
}

func (cm *CardMessage) SetType(msgType iremote.MsgType) *CardMessage {
	cm.Type_ = msgType
	return cm
}

func (cm *CardMessage) Type() iremote.MsgType {
	return cm.Type_
}

func (cm *CardMessage) Cards() map[string]card.Card {
	return cm.Digest_.Cards()
}

func (cm *CardMessage) Cluster() string {
	return cm.Cluster_
}

func (cm *CardMessage) From() card.Card {
	return cm.From_
}

func (cm *CardMessage) To() card.Card {
	return cm.To_
}

func (cm *CardMessage) TimeStamp() int64 {
	return cm.Digest_.TimeStamp()
}

func (cm *CardMessage) Status() constants.Header {
	return cm.Status_
}

func (cm *CardMessage) Terminate(key string) {
	delete(cm.Digest_.Cards_, key)
}

func (cm *CardMessage) Stamp() {
	cm.SetTimeStamp(time.Now().Unix())
}

func (cm *CardMessage) Digest() iremote.IDigest {
	return &cm.Digest_
}

func (cm *CardMessage) Update(dgst iremote.IDigest) {
	cm.SetCards(dgst.Cards())
	cm.SetTimeStamp(dgst.TimeStamp())
}
