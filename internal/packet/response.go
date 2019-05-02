//
// Description:
// Author: Rodrigo Freitas
// Created at: Tue Apr 23 15:29:24 -03 2019
//
package packet

import (
	"net/textproto"
)

type Response struct {
	StatusCode int
	StatusText string
	Headers    textproto.MIMEHeader
	Body       []byte
}

func newResponse() *Response {
	return &Response{
		Headers: make(textproto.MIMEHeader),
	}
}
