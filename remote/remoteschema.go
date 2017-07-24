package remote

const (
	ArgTypeInteger          = "integer"
	ArgTypeNumber           = "number"
	ArgTypeObject           = "object"
	ArgTypeBoolean          = "boolean"
	ArgTypeNull             = "null"
	ArgTypeArray            = "array"
	ConstraintTypeMax       = "maximum"
	ConstraintTypeMin       = "minimum"
	ConstraintTypeMaxLength = "maxLength"
	ConstraintTypeMinLength = "minLength"
	ConstraintTypePattern   = "pattern"
	ConstraintTypeMaxItems  = "maxItems"
	ConstraintTypeMinItems  = "minItems"
	ConstraintTypeEnum      = "enum" // value of interface{} should accept a slice
)

type Argument struct {
	DataType    string
	Description string
	Constraints []Constraint
	Required    bool
}

type Constraint struct {
	Key   string
	Value interface{}
}

func (arg *Argument) SerializeToJSON() {

}
