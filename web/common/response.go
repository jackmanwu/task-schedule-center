package common

const (
	Success int8 = iota
)

type Result struct {
	Code int8        `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func NewResult(code int8, msg string, data interface{}) *Result {
	return &Result{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func NewSuccess() *Result {
	return NewResult(Success, "succ", nil)
}

func NewSuccessWithData(data interface{}) *Result {
	return NewResult(Success, "succ", data)
}
