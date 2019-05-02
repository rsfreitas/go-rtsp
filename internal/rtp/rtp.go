//
// Description:
// Author: Rodrigo Freitas
// Created at: Mon Apr 29 18:34:47 -03 2019
//
package rtp

import (
	"fmt"
	"net"

	"github.com/wernerd/gortp/src/net/rtp"
)

// Setup holds options to initialize a RTP session between server and client.
type Setup struct {
	ServerPort  int
	ServerAddr  string
	ClientPorts []int
	ClientAddr  string
}

// Session holds a RTP session, to transfer data to the client.
type Session struct {
	s            *rtp.Session
	stopReceiver chan bool
	stopSender   chan bool
	port         int
}

func (r *Session) Close() {
	// closes goroutines (send/recv)
	r.stopReceiver <- true
	fmt.Println("Closing RTP session")
}

func (r *Session) Pause() {
}

func (r *Session) Port() int {
	return r.port
}

func rtpReceiver(r *Session) {
	dataReceiver := r.s.CreateDataReceiveChan()

loop:
	for {
		select {
		case rp := <-dataReceiver:
			rp.FreePacket()

		case <-r.stopReceiver:
			break loop
		}
	}

	fmt.Println("Closing receiver")
}

func rtpSender(r *Session) {
}

func NewSession(options Setup) *Session {
	serverAddr, _ := net.ResolveIPAddr("ip", options.ServerAddr)
	clientAddr, _ := net.ResolveIPAddr("ip", options.ClientAddr)

	tpLocal, _ := rtp.NewTransportUDP(serverAddr, options.ServerPort, "")
	session := rtp.NewSession(tpLocal, tpLocal)
	portB := options.ClientPorts[0] + 1

	if len(options.ClientPorts) > 1 {
		portB = options.ClientPorts[1]
	}

	session.AddRemote(&rtp.Address{clientAddr.IP, options.ClientPorts[0], portB, ""})

	r := &Session{
		port:         options.ServerPort,
		s:            session,
		stopReceiver: make(chan bool),
		stopSender:   make(chan bool),
	}

	go rtpReceiver(r)
	session.ListenOnTransports()
	go rtpSender(r)

	fmt.Println("New RTP session")
	return r
}
