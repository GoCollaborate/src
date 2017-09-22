package remoteshared

import (
	"github.com/GoCollaborate/constants"
)

type Parameterizable interface {
	SerializeToJSON() string
	Validate() bool
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
	return Constraint{constants.ConstraintTypeMin, val}
}

func DefaultConstraintXMinimum(val float64) Constraint {
	return Constraint{constants.ConstraintTypeXMin, val}
}

func DefaultConstraintMaximum(val float64) Constraint {
	return Constraint{constants.ConstraintTypeMax, val}
}

func DefaultConstraintXMaximum(val float64) Constraint {
	return Constraint{constants.ConstraintTypeXMax, val}
}

func DefaultConstraintMaxLength(l uint64) Constraint {
	return Constraint{constants.ConstraintTypeMaxLength, l}
}

func DefaultConstraintMinLength(l uint64) Constraint {
	return Constraint{constants.ConstraintTypeMinLength, l}
}

func DefaultConstraintPattern(pt string) Constraint {
	return Constraint{constants.ConstraintTypePattern, pt}
}

func DefaultConstraintMaxItems(it uint64) Constraint {
	return Constraint{constants.ConstraintTypeMaxItems, it}
}

func DefaultConstraintMinItems(it uint64) Constraint {
	return Constraint{constants.ConstraintTypeMinItems, it}
}

func DefaultConstraintUniqueItems(uniq bool) Constraint {
	return Constraint{constants.ConstraintTypeUniqueItems, uniq}
}

func DefaultConstraintMaxProperties(mp uint64) Constraint {
	return Constraint{constants.ConstraintTypeMaxProperties, mp}
}

func DefaultConstraintMinProperties(mp uint64) Constraint {
	return Constraint{constants.ConstraintTypeMinProperties, mp}
}

func DefaultConstraintArrayAllOf(ins interface{}) Constraint {
	return Constraint{constants.ConstraintTypeAllOf, ins}
}

func DefaultConstraintArrayAnyOf(ins interface{}) Constraint {
	return Constraint{constants.ConstraintTypeAnyOf, ins}
}

func DefaultConstraintArrayOneOf(ins interface{}) Constraint {
	return Constraint{constants.ConstraintTypeOneOf, ins}
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
