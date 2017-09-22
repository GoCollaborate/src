package remoteshared

import (
	"github.com/GoCollaborate/constants"
	"time"
)

// From() and To() are not mandatory for Case
type ICardMessage interface {
	Cards() map[string]Card
	Cluster() string
	TimeStamp() int64
}

type CardMessage struct {
	Cluster_ string           `json:"cluster,omitempty"`
	From_    Card             `json:"from"`
	To_      Card             `json:"to"`
	Cards_   map[string]Card  `json:"cards"`
	Ts_      int64            `json:"timestamp"`
	Status_  constants.Header `json:"status"`
}

func NewCardMessage() *CardMessage {
	return new(CardMessage)
}

func NewCardMessageWithOptions(cluster string, from Card, to Card, cards map[string]Card, timestamp int64) *CardMessage {
	message := CardMessage{}
	message.SetCluster(cluster).SetFrom(from).SetTo(to).SetCards(cards).SetTimeStamp(timestamp).SetStatus(constants.GossipHeaderUnknownError)
	return &message
}

func (cm *CardMessage) SetCluster(cluster string) *CardMessage {
	cm.Cluster_ = cluster
	return cm
}

func (cm *CardMessage) SetFrom(from Card) *CardMessage {
	cm.From_ = from
	return cm
}

func (cm *CardMessage) SetTo(to Card) *CardMessage {
	cm.To_ = to
	return cm
}

func (cm *CardMessage) SetCards(cards map[string]Card) *CardMessage {
	cm.Cards_ = cards
	return cm
}

func (cm *CardMessage) SetTimeStamp(timestamp int64) *CardMessage {
	cm.Ts_ = timestamp
	return cm
}

func (cm *CardMessage) SetStatus(status constants.Header) *CardMessage {
	cm.Status_ = status
	return cm
}

func (cm *CardMessage) Cards() map[string]Card {
	return cm.Cards_
}

func (cm *CardMessage) Cluster() string {
	return cm.Cluster_
}

func (cm *CardMessage) From() Card {
	return cm.From_
}

func (cm *CardMessage) To() Card {
	return cm.To_
}

func (cm *CardMessage) TimeStamp() int64 {
	return cm.Ts_
}

func (cm *CardMessage) Status() constants.Header {
	return cm.Status_
}

func (cm *CardMessage) Terminate(key string) {
	delete(cm.Cards_, key)
}

func (cm *CardMessage) Stamp() {
	cm.SetTimeStamp(time.Now().Unix())
}
