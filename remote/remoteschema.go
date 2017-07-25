package remote

import (
	"github.com/GoCollaborate/constants"
)

type Parameterizable interface {
	SerializeToJSON() string
	Raw() interface{}
}

type Argument struct {
	DataType    string
	Description string
	Constraints []Constraint
	Required    bool
	Value       interface{}
}

type Constraint struct {
	Key   string
	Value interface{}
}

func (arg *Argument) SerializeToJSON() string {
	return ""
}

func (arg *Argument) Raw() interface{} {
	return nil
}

func DefaultArgString(val interface{}) *Argument {
	return &Argument{constants.ArgTypeString, "Remote Argument Type: String", []Constraint{}, false, val}
}

func DefaultArgNumber(val interface{}) *Argument {
	return &Argument{constants.ArgTypeNumber, "Remote Argument Type: Number", []Constraint{}, false, val}
}

func DefaultArgInteger(val interface{}) *Argument {
	return &Argument{constants.ArgTypeInteger, "Remote Argument Type: Integer", []Constraint{}, false, val}
}

func DefaultArgObject(val interface{}) *Argument {
	return &Argument{constants.ArgTypeObject, "Remote Argument Type: Object", []Constraint{}, false, val}
}

func DefaultArgBoolean(val interface{}) *Argument {
	return &Argument{constants.ArgTypeBoolean, "Remote Argument Type: Boolean", []Constraint{}, false, val}
}

func DefaultArgArray(val interface{}) *Argument {
	return &Argument{constants.ArgTypeArray, "Remote Argument Type: Array", []Constraint{}, false, val}
}

func DefaultArgNull(val interface{}) *Argument {
	return &Argument{constants.ArgTypeNull, "Remote Argument Type: Null", []Constraint{}, false, val}
}
