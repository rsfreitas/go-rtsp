//
// Description:
// Author: Rodrigo Freitas
// Created at: Sun Apr 28 10:05:05 -03 2019
//
package rtsp

import (
	"github.com/rsfreitas/go-rtsp/internal/packet"
)

type recordMethod struct {
	clientHandler ClientRecord
}

func (r *recordMethod) Verify(p *packet.Packet, handler interface{}) error {
	if m, ok := handler.(interface{ Record() }); ok {
		r.clientHandler = m
	} else {
		return &methodError{"client method not implemented", true}
	}

	return nil
}

func (r *recordMethod) Handle(p *packet.Packet) {
	if r.clientHandler != nil {
		r.clientHandler.Record()
	}
}

func (r *recordMethod) Type() methodType {
	return methodRecord
}
