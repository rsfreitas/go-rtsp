//
// Description:
// Author: Rodrigo Freitas
// Created at: Sun Apr 28 09:50:44 -03 2019
//
package rtsp

import (
	"net/http"
	"strings"

	"github.com/rsfreitas/go-rtsp/internal/packet"
)

type optionsMethod struct {
	clientHandler interface{}
}

func (o *optionsMethod) Verify(p *packet.Packet, handler interface{}) error {
	o.clientHandler = handler

	return nil
}

func (o *optionsMethod) Handle(p *packet.Packet) {
	var options strings.Builder
	options.WriteString("OPTIONS, DESCRIBE, SETUP")

	if _, ok := o.clientHandler.(interface{ Play() }); ok {
		options.WriteString(", PLAY")
	}

	if _, ok := o.clientHandler.(interface{ Pause() }); ok {
		options.WriteString(", PAUSE")
	}

	if _, ok := o.clientHandler.(interface{ Teardown() }); ok {
		options.WriteString(", TEARDOWN")
	}

	if _, ok := o.clientHandler.(interface{ Record() }); ok {
		options.WriteString(", RECORD")
	}

	if _, ok := o.clientHandler.(interface{ Announce() }); ok {
		options.WriteString(", ANNOUNCE")
	}

	if _, ok := o.clientHandler.(interface{ GetParameter() }); ok {
		options.WriteString(", GET_PARAMETER")
	}

	if _, ok := o.clientHandler.(interface{ SetParameter() }); ok {
		options.WriteString(", SET_PARAMETER")
	}

	p.Response.StatusCode = http.StatusOK
	p.Response.StatusText = http.StatusText(http.StatusOK)

	p.Response.Headers.Add("Content-Length", "0")
	p.Response.Headers.Add("Public", options.String())
}

func (o *optionsMethod) Type() methodType {
	return methodOptions
}
