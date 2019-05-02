//
// Description:
// Author: Rodrigo Freitas
// Created at: Mon Apr 29 18:50:13 -03 2019
//
package header_test

import (
	"fmt"
	"testing"

	"github.com/rsfreitas/go-rtsp/internal/header"
	"github.com/stretchr/testify/assert"
)

func TestNewTransport(t *testing.T) {
	assert := assert.New(t)

	r, err := header.NewTransportFromString("RTP/AVP/UDP;unicast;destination=192.168.10.95;client_port=8000-8001;server_port=39000-35968;ssrc=46a81ad7;mode=play")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(r)
	assert.Equal(1, 1, "Should be equal")
}
