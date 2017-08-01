package remote

import (
	"encoding/json"
	"github.com/GoCollaborate/constants"
)

type Parameterizable interface {
	SerializeToJSON() string
	Raw() interface{}
}

type Parameter struct {
	Type        string       `json:"type"`
	Description string       `json:"description"`
	Constraints []Constraint `json:"constraints"`
	Required    bool         `json:"required"`
}

type Constraint struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func (arg *Parameter) SerializeToJSON() string {
	mal, err := json.Marshal(*arg)
	if err != nil {
		// log here
	}
	return string(mal)
}

func (arg *Parameter) Raw() interface{} {
	return nil
}

func DefaultArgString(val interface{}) *Parameter {
	return &Parameter{constants.ArgTypeString, "Remote Argument Type: String", []Constraint{}, false}
}

func DefaultArgNumber(val interface{}) *Parameter {
	return &Parameter{constants.ArgTypeNumber, "Remote Argument Type: Number", []Constraint{}, false}
}

func DefaultArgInteger(val interface{}) *Parameter {
	return &Parameter{constants.ArgTypeInteger, "Remote Argument Type: Integer", []Constraint{}, false}
}

func DefaultArgObject(val interface{}) *Parameter {
	return &Parameter{constants.ArgTypeObject, "Remote Argument Type: Object", []Constraint{}, false}
}

func DefaultArgBoolean(val interface{}) *Parameter {
	return &Parameter{constants.ArgTypeBoolean, "Remote Argument Type: Boolean", []Constraint{}, false}
}

func DefaultArgArray(val interface{}) *Parameter {
	return &Parameter{constants.ArgTypeArray, "Remote Argument Type: Array", []Constraint{}, false}
}

func DefaultArgNull(val interface{}) *Parameter {
	return &Parameter{constants.ArgTypeNull, "Remote Argument Type: Null", []Constraint{}, false}
}

func DefaultConstraintMinimum(val float64) Constraint {
	return Constraint{"minimum", val}
}

func DefaultConstraintXMinimum(val float64) Constraint {
	return Constraint{"exclusiveMinimum", val}
}

func DefaultConstraintMaximum(val float64) Constraint {
	return Constraint{"maximum", val}
}

func DefaultConstraintXMaximum(val float64) Constraint {
	return Constraint{"exclusiveMaximum", val}
}

func DefaultConstraintMaxLength(l uint64) Constraint {
	return Constraint{"maxLength", l}
}

func DefaultConstraintMinLength(l uint64) Constraint {
	return Constraint{"minLength", l}
}

func DefaultConstraintPattern(pt string) Constraint {
	return Constraint{"pattern", pt}
}

func DefaultConstraintMaxItems(it uint64) Constraint {
	return Constraint{"maxItems", it}
}

func DefaultConstraintMinItems(it uint64) Constraint {
	return Constraint{"minItems", it}
}

func DefaultConstraintUniqueItems(uniq bool) Constraint {
	return Constraint{"uniqueItems", uniq}
}

func DefaultConstraintMaxProperties(mp uint64) Constraint {
	return Constraint{"maxProperties", mp}
}

func DefaultConstraintMinProperties(mp uint64) Constraint {
	return Constraint{"minProperties", mp}
}

func DefaultConstraintArrayAllOf(ins interface{}) Constraint {
	return Constraint{"allOf", ins}
}

func DefaultConstraintArrayAnyOf(ins interface{}) Constraint {
	return Constraint{"anyOf", ins}
}

func DefaultConstraintArrayOneOf(ins interface{}) Constraint {
	return Constraint{"oneOf", ins}
}

func UnmarshalParameters(original []interface{}) []Parameter {
	var parameters []Parameter
	for _, o := range original {
		oo := o.(map[string]interface{})
		parameters = append(parameters, Parameter{oo["type"].(string), oo["description"].(string), UnmarshalConstraints(oo["constraints"].([]interface{})), oo["required"].(bool)})
	}
	return parameters
}

func UnmarshalConstraints(original []interface{}) []Constraint {
	var constraints []Constraint
	for _, o := range original {
		oo := o.(map[string]interface{})
		constraints = append(constraints, Constraint{oo["key"].(string), oo["value"]})
	}
	return constraints
}

func UnmarshalStringArray(original []interface{}) []string {
	var strings []string
	for _, o := range original {
		oo := o.(string)
		strings = append(strings, oo)
	}
	return strings
}
