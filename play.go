//
// Description:
// Author: Rodrigo Freitas
// Created at: Sun Apr 28 10:05:05 -03 2019
//
package rtsp

import (
	"github.com/rsfreitas/go-rtsp/internal/packet"
	"github.com/rsfreitas/go-rtsp/internal/rtp"
)

type playMethod struct {
	clientHandler ClientPlay

	ActiveSessions map[string]*rtp.Session
}

func (p *playMethod) Verify(pkt *packet.Packet, handler interface{}) error {
	if m, ok := handler.(interface{ Play() }); ok {
		p.clientHandler = m
	} else {
		return &methodError{"client method not implemented", true}
	}

	return nil
}

func (p *playMethod) Handle(pkt *packet.Packet) {
	if p.clientHandler != nil {
		p.clientHandler.Play()
	}
}

func (p *playMethod) Type() methodType {
	return methodPlay
}
