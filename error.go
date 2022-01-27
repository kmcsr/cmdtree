
package cmdtree

import (
	// errors "errors"
	fmt "fmt"
)

type ErrUnknownCommand struct{
	cmd string
}

func (e *ErrUnknownCommand)Error()(string){
	return fmt.Sprintf("Unknown command: '%s'", e.cmd)
}

type ErrUnknownArg struct{
	cmd string
}

func (e *ErrUnknownArg)Error()(string){
	if len(e.cmd) == 0 {
		return "Unknown argument"
	}
	return fmt.Sprintf("Unknown argument: '%s'", trimRightStr(e.cmd))
}

type ErrArgOutOfRange struct{
	min, max interface{}
}

func (e *ErrArgOutOfRange)Error()(string){
	return fmt.Sprintf("Argument out of range [%v, %v]", e.min, e.max)
}

type ErrArgRequest struct{
	Err error
}

func (e *ErrArgRequest)Error()(string){
	return "Argument error: " + e.Err.Error()
}
