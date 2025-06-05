package errorx

import (
	"fmt"
	"strings"
)

type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type ErrInterface interface {
	Code() int
	String() string
}

func NewError(err ErrInterface, msg ...string) error {
	return &CodeError{Code: err.Code(), Msg: fmt.Sprint(err.String(), strings.Join(msg, " "))}
}

func (e *CodeError) Error() string {
	return e.Msg
}
