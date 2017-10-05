package parameter

import (
	"encoding/json"
	"github.com/GoCollaborate/artifacts/restful"
	"github.com/GoCollaborate/constants"
	"github.com/GoCollaborate/logger"
)

type Parameter struct {
	Type        string       `json:"type"`
	Description string       `json:"description"`
	Constraints []Constraint `json:"constraints"`
	Required    bool         `json:"required"`
}

func (par *Parameter) SerializeToJSON() string {
	mal, err := json.Marshal(*par)
	if err != nil {
		logger.LogWarning("Internal error during serialization: " + err.Error())
	}
	return string(mal)
}

func (par *Parameter) Validate(dat *restful.Resource) bool {
	if par.Required && dat == nil {
		return false
	}
	if par.Type != dat.Type {
		return false
	}
	switch par.Type {
	case constants.ArgTypeInteger, constants.ArgTypeNumber:
		for _, cons := range par.Constraints {
			switch cons.Key {
			case constants.ConstraintTypeMax:
			case constants.ConstraintTypeMin:
			case constants.ConstraintTypeEnum:
			default:
				return false
			}
		}
	case constants.ArgTypeString:
		for _, cons := range par.Constraints {
			switch cons.Key {
			case constants.ConstraintTypeMaxLength:
			case constants.ConstraintTypeMinLength:
			case constants.ConstraintTypePattern:
			case constants.ConstraintTypeEnum:
			default:
				return false
			}
		}
	case constants.ArgTypeArray:
		for _, cons := range par.Constraints {
			switch cons.Key {
			case constants.ConstraintTypeMaxItems:
			case constants.ConstraintTypeMinItems:
			case constants.ConstraintTypeUniqueItems:
			case constants.ConstraintTypeAllOf:
			case constants.ConstraintTypeOneOf:
			case constants.ConstraintTypeAnyOf:
			default:
				return false
			}
		}
	case constants.ArgTypeObject:
		for _, cons := range par.Constraints {
			switch cons.Key {
			case constants.ConstraintTypeMaxProperties:
			case constants.ConstraintTypeMinProperties:
			default:
				return false
			}
		}
	default:
		return false
	}
	return true
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
