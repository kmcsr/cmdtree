
package cmdtree

import (
	errors "errors"
	fmt "fmt"
)

var (
	ErrUnknownArg = errors.New("Unknown argument")
)

type ErrArgOutOfRange struct{
	min, max interface{}
}

func (e *ErrArgOutOfRange)Error()(string){
	return fmt.Sprintf("Argument out of range [%v, %v]", e.min, e.max)
}
