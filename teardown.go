//
// Description:
// Author: Rodrigo Freitas
// Created at: Sun Apr 28 10:05:05 -03 2019
//
package rtsp

import (
	"net/http"

	"github.com/rsfreitas/go-rtsp/internal/adt"
	"github.com/rsfreitas/go-rtsp/internal/packet"
	"github.com/rsfreitas/go-rtsp/internal/rtp"
)

type teardownMethod struct {
	clientHandler ClientTeardown

	ActiveSessions map[string]*rtp.Session
	AvailablePorts *adt.RangeBox
}

func (t *teardownMethod) Verify(p *packet.Packet, handler interface{}) error {
	if m, ok := handler.(interface{ Teardown() }); ok {
		t.clientHandler = m
		//	} else {
		//		return &methodError{"client method not implemented", true}
	}

	return nil
}

func (t *teardownMethod) Handle(p *packet.Packet) {
	var (
		clientSession string
		session       *rtp.Session
	)

	if t.clientHandler != nil {
		t.clientHandler.Teardown()
	}

	if field, ok := p.Request.Headers["Session"]; ok {
		if t.ActiveSessions == nil {
			p.Response.StatusCode = StatusAggregateOperationNotAllowed
			p.Response.StatusText = StatusText(p.Response.StatusCode)
			return
		}

		clientSession = field[0]
		session, ok = t.ActiveSessions[clientSession]

		if !ok {
			p.Response.StatusCode = StatusSessionNotFound
			p.Response.StatusText = StatusText(p.Response.StatusCode)
			return
		}
	}

	port := session.Port()
	t.AvailablePorts.Release(uint32(port))
	session.Close()

	p.Response.StatusCode = http.StatusOK
	p.Response.StatusText = http.StatusText(http.StatusOK)
}

func (t *teardownMethod) Type() methodType {
	return methodTeardown
}
