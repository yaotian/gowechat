package base

import "fmt"

const (
	ReturnCodeSuccess = "SUCCESS"
	ReturnCodeFail    = "FAIL"
)

const (
	ResultCodeSuccess = "SUCCESS"
	ResultCodeFail    = "FAIL"
)

//Error error
type Error struct {
	XMLName    struct{} `xml:"xml"                  json:"-"`
	ReturnCode string   `xml:"return_code"          json:"return_code"`
	ReturnMsg  string   `xml:"return_msg,omitempty" json:"return_msg,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("return_code: %q, return_msg: %q", e.ReturnCode, e.ReturnMsg)
}
