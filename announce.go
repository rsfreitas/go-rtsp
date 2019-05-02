//
// Description:
// Author: Rodrigo Freitas
// Created at: Sun Apr 28 10:05:05 -03 2019
//
package rtsp

import (
	"github.com/rsfreitas/go-rtsp/internal/packet"
)

type announceMethod struct {
	clientHandler ClientAnnounce
}

func (a *announceMethod) Verify(p *packet.Packet, handler interface{}) error {
	if m, ok := handler.(interface{ Announce() }); ok {
		a.clientHandler = m
	} else {
		return &methodError{"client method not implemented", true}
	}

	return nil
}

func (a *announceMethod) Handle(p *packet.Packet) {
	if a.clientHandler != nil {
		a.clientHandler.Announce()
	}
}

func (a *announceMethod) Type() methodType {
	return methodAnnounce
}
