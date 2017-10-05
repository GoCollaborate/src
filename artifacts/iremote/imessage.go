package iremote

// From() and To() are not mandatory for Case
type ICardMessage interface {
	Cluster() string
	Type() MsgType
	IDigest
	Digest() IDigest
	Update(IDigest)
}

type MsgType string

var (
	MsgTypeSync MsgType = MsgType("SYNC")
	MsgTypeAck  MsgType = MsgType("ACK")
	MsgTypeAck2 MsgType = MsgType("ACK2")
	MsgTypeAck3 MsgType = MsgType("ACK3")
)
