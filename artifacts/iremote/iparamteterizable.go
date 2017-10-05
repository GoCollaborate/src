package iremote

type IParameterizable interface {
	SerializeToJSON() string
	Validate() bool
}
