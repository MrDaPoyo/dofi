package balena

type Return struct {
	Value interface{}
}

func NewReturn(value interface{}) Return {
	return Return{Value: value}
}

func (r Return) GetValue() interface{} {
	return r.Value
}
