package balena

import (
	"fmt"

	"github.com/mrdapoyo/dofi/balena/token"
)

type RuntimeError struct {
	Token   token.Token
	Message string
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("[line %d] %s", e.Token.Line, e.Message)
}
