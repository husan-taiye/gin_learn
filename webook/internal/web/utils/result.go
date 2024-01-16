package utils

type Result struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Data    any    `json:"data"`
	Success any    `json:"success"`
}
