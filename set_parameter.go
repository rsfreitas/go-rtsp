//
// Description:
// Author: Rodrigo Freitas
// Created at: Sun Apr 28 10:05:05 -03 2019
//
package rtsp

import (
	"github.com/rsfreitas/go-rtsp/internal/packet"
)

type setParameterMethod struct {
	clientHandler ClientSetParameter
}

func (s *setParameterMethod) Verify(p *packet.Packet, handler interface{}) error {
	if m, ok := handler.(interface{ SetParameter() }); ok {
		s.clientHandler = m
	} else {
		return &methodError{"client method not implemented", true}
	}

	return nil
}

func (s *setParameterMethod) Handle(p *packet.Packet) {
	if s.clientHandler != nil {
		s.clientHandler.SetParameter()
	}
}

func (s *setParameterMethod) Type() methodType {
	return methodSetParameter
}
