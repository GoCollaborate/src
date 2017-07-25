package remote

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
