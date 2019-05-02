//
// Description:
// Author: Rodrigo Freitas
// Created at: Thu Apr 25 16:34:13 -03 2019
//
package rtsp

import (
	"github.com/rsfreitas/go-rtsp/internal/packet"
)

// methodType is a method identification.
type methodType int

const (
	methodOptions methodType = iota + 1
	methodDescribe
	methodSetup
	methodPlay
	methodPause
	methodTeardown
	methodRecord
	methodAnnounce
	methodGetParameter
	methodSetParameter
)

// method defines a method set of functions it must have to be internally
// supported.
type method interface {
	// Verify must check if request contains every required information for
	// the specific method in its headers or body. It returns
	Verify(*packet.Packet, interface{}) error

	// Handle should call the method handler function and it also should be
	// responsible for preparing the Response for the client.
	Handle(*packet.Packet)

	// Type must return the method identification
	Type() methodType
}
