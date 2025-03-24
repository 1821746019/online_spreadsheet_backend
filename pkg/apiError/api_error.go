package apiError

import "github.com/sztu/mutli-table/pkg/code"

type ApiError struct {
	Code code.RespCode `json:"code"`
	Msg  string        `json:"msg"`
}

func (e ApiError) Error() string {
	return e.Msg
}
