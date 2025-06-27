package balena

type Return struct {
	value interface{}
}

func NewReturn(value interface{}) *Return {
	return &Return{value: value}
}

func (r *Return) Value() interface{} {
	return r.value
}