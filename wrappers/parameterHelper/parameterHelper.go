package parameterHelper

import (
	"github.com/GoCollaborate/artifacts/parameter"
)

func UnmarshalParameters(original []interface{}) []parameter.Parameter {
	var parameters []parameter.Parameter
	for _, o := range original {
		oo := o.(map[string]interface{})
		parameters = append(parameters, parameter.Parameter{oo["type"].(string), oo["description"].(string), UnmarshalConstraints(oo["constraints"].([]interface{})), oo["required"].(bool)})
	}
	return parameters
}

func UnmarshalConstraints(original []interface{}) []parameter.Constraint {
	var constraints []parameter.Constraint
	for _, o := range original {
		oo := o.(map[string]interface{})
		constraints = append(constraints, parameter.Constraint{oo["key"].(string), oo["value"]})
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
