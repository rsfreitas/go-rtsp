//
// Description:
// Author: Rodrigo Freitas
// Created at: Sun Apr 28 10:05:05 -03 2019
//
package rtsp

import (
	"github.com/rsfreitas/go-rtsp/internal/packet"
)

type getParameterMethod struct {
	clientHandler ClientGetParameter
}

func (g *getParameterMethod) Verify(p *packet.Packet, handler interface{}) error {
	if m, ok := handler.(interface{ GetParameter() }); ok {
		g.clientHandler = m
	} else {
		return &methodError{"client method not implemented", true}
	}

	return nil
}

func (g *getParameterMethod) Handle(p *packet.Packet) {
	if g.clientHandler != nil {
		g.clientHandler.GetParameter()
	}
}

func (g *getParameterMethod) Type() methodType {
	return methodGetParameter
}
