package balena

import "github.com/mrdapoyo/dofi/balena/token"

type RuntimeError struct {
	Token   token.Token
	Message string
}

func (e *RuntimeError) Error() string {
	return e.Message
}