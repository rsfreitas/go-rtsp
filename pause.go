//
// Description:
// Author: Rodrigo Freitas
// Created at: Sun Apr 28 10:05:05 -03 2019
//
package rtsp

import (
	"net/http"

	"github.com/rsfreitas/go-rtsp/internal/packet"
	"github.com/rsfreitas/go-rtsp/internal/rtp"
)

type pauseMethod struct {
	clientHandler ClientPause

	ActiveSessions map[string]*rtp.Session
}

func (p *pauseMethod) Verify(pkt *packet.Packet, handler interface{}) error {
	if m, ok := handler.(interface{ Pause() }); ok {
		p.clientHandler = m
		//	} else {
		//		return &methodError{"client method not implemented", true}
	}

	return nil
}

func (p *pauseMethod) Handle(pkt *packet.Packet) {
	var (
		clientSession string
		session       *rtp.Session
	)

	if p.clientHandler != nil {
		p.clientHandler.Pause()
	}

	if field, ok := pkt.Request.Headers["Session"]; ok {
		if p.ActiveSessions == nil {
			pkt.Response.StatusCode = StatusAggregateOperationNotAllowed
			pkt.Response.StatusText = StatusText(pkt.Response.StatusCode)
			return
		}

		clientSession = field[0]
		session, ok = p.ActiveSessions[clientSession]

		if !ok {
			pkt.Response.StatusCode = StatusSessionNotFound
			pkt.Response.StatusText = StatusText(pkt.Response.StatusCode)
			return
		}
	}

	session.Pause()
	pkt.Response.StatusCode = http.StatusOK
	pkt.Response.StatusText = http.StatusText(http.StatusOK)
}

func (p *pauseMethod) Type() methodType {
	return methodPause
}
