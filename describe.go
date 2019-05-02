//
// Description:
// Author: Rodrigo Freitas
// Created at: Sun Apr 28 09:59:01 -03 2019
//
package rtsp

import (
	"fmt"
	"net/http"

	"github.com/rsfreitas/go-rtsp/internal/packet"
	"github.com/rsfreitas/go-rtsp/internal/sdp"
)

type describeMethod struct {
	Session *sdp.Session
}

func (d *describeMethod) Verify(p *packet.Packet, handler interface{}) error {
	return nil
}

func (d *describeMethod) isAcceptable(p *packet.Packet) bool {
	accept, ok := p.Request.Headers["Accept"]

	if !ok {
		// Assumes the default, which is 'application/sdp' and continue
		return true
	}

	if accept[0] == "application/sdp" {
		return true
	}

	return false
}

func (d *describeMethod) Handle(p *packet.Packet) {
	if !d.isAcceptable(p) {
		p.Response.StatusCode = http.StatusNotAcceptable
		p.Response.StatusText = http.StatusText(http.StatusNotAcceptable)
		return
	}

	// We send the SDP representation of options used when creating the
	// server
	p.Response.Body = d.Session.Bytes()
	//	var b []byte
	//	p.Response.Body = d.Session.AppendTo(b)

	p.Response.StatusCode = http.StatusOK
	p.Response.StatusText = http.StatusText(http.StatusOK)

	p.Response.Headers.Add("Content-Length",
		fmt.Sprintf("%d", len(string(p.Response.Body))))
}

func (d *describeMethod) Type() methodType {
	return methodDescribe
}
