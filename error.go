//
// Description: Internal errors.
// Author: Rodrigo Freitas
// Created at: Sun Apr 28 09:44:32 -03 2019
//
package rtsp

// parserError represents an error when parsing the received data from a
// client.
type parserError struct {
}

// methodError represents an error related to the method received from a
// client.
type methodError struct {
	err       string
	hasMethod bool
}

func (e *methodError) Error() string {
	return e.err
}

func (e *methodError) methodNotAllowed() bool {
	return e.hasMethod == false
}
