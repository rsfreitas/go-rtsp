//
// Description:
// Author: Rodrigo Freitas
// Created at: Mon Apr 22 14:04:27 -03 2019
//
package rtsp

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/rsfreitas/go-rtsp/internal/adt"
	"github.com/rsfreitas/go-rtsp/internal/packet"
	"github.com/rsfreitas/go-rtsp/internal/rtp"
	"github.com/rsfreitas/go-rtsp/internal/sdp"
)

type MediaSetup struct {
	Port       int
	ClientHost string
}

// ServerSetup holds all available options to create a Server object.
type ServerSetup struct {
	Port       int
	Username   string
	Password   string
	AuthType   AuthorizationType
	UDPPortMin uint32
	UDPPortMax uint32

	// MediaSetup must contain all video spec that will be available to
	// clients through the DESCRIBE request.
	*MediaSetup
}

// Server is the server object.
type Server struct {
	ServerSetup

	rtspListener   *net.TCPListener
	handler        interface{}
	shutdown       chan bool
	clientTimeout  int
	activeSessions map[string]*rtp.Session
	availablePorts *adt.RangeBox
	session        *sdp.Session
}

const (
	defaultRequestBufferSize = 10240
	osReceiveBufferSize      = 51200
)

// Close releases all internal server resources.
func (s *Server) Close() {
	s.rtspListener.Close()
	close(s.shutdown)
}

// Start puts the server to receive incoming connections.
func (s *Server) Start() {
	for {
		c, err := s.rtspListener.AcceptTCP()

		if err != nil {
			if _, ok := err.(net.Error); ok && strings.HasSuffix(err.Error(), ": use of closed network connection") {
				break
			}

			fmt.Println(err)
			continue
		}

		c.SetReadBuffer(osReceiveBufferSize)

		// Handle the new connection
		go s.handleConnection(c)
	}
}

// handleConnection handles a client connection, receiving and handling a method.
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, defaultRequestBufferSize)
	timeoutDuration := time.Duration(s.clientTimeout) * time.Millisecond

	for {
		p := packet.NewPacket()
		conn.SetReadDeadline(time.Now().Add(timeoutDuration))
		n, err := conn.Read(buffer)

		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				continue
			}

			if err == io.EOF {
				break
			}

			return
		}

		// Does the package still have data to be read?
		for p.StillNeedsRead(buffer, n) == true {
			conn.SetReadDeadline(time.Now().Add(timeoutDuration))
			m, err := conn.Read(buffer[n:])

			if err != nil {
				if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
					continue
				}

				if err == io.EOF {
					break
				}

				return
			}

			n += m
		}

		var r []byte

		if err := p.UnmarshalRequest(buffer, n); err != nil {
			r = p.MarshalResponseError(err)
		} else {
			fmt.Println(p.Request)
			s.handleRequestOption(p)
			r, err = p.MarshalResponse()

			if err != nil {
				return
			}

			fmt.Println(string(r))
		}

		conn.Write(r)
	}

	return
}

// handleRequestOption calls the received request option callback filling in
// packet with response. The callback to be called will be handled if the
// internal handler supports it or a default will be used.
func (s *Server) handleRequestOption(p *packet.Packet) {
	var m method

	switch p.Request.Method {
	case "OPTIONS":
		m = &optionsMethod{}

	case "DESCRIBE":
		m = &describeMethod{
			Session: s.session,
		}

	case "SETUP":
		port, err := s.availablePorts.Request()

		if err != nil {
			s.internalServerError(p, "server doesn't have available ports to transfer")
			return
		}

		m = &setupMethod{
			ActiveSessions: s.activeSessions,
			ServerPortMin:  int(port),
			ServerPortMax:  int(port + 1),
		}

	case "PLAY":
		m = &playMethod{
			ActiveSessions: s.activeSessions,
		}

	case "PAUSE":
		m = &pauseMethod{
			ActiveSessions: s.activeSessions,
		}

	case "TEARDOWN":
		m = &teardownMethod{
			ActiveSessions: s.activeSessions,
			AvailablePorts: s.availablePorts,
		}

	case "RECORD":
		m = &recordMethod{}

	case "ANNOUNCE":
		m = &announceMethod{}

	case "GET_PARAMETER":
		m = &getParameterMethod{}

	case "SET_PARAMETER":
		m = &setParameterMethod{}
	}

	if m == nil {
		s.unsupportedMethod(p)
		return
	}

	if err := m.Verify(p, s.handler); err != nil {
		if err, ok := err.(*methodError); ok {
			if err.methodNotAllowed() {
				s.methodNotAllowed(p)
			}
		}
	} else {
		m.Handle(p)
	}
}

// unsupportedMethod is a method to set the server response as HTTP 501.
func (s *Server) unsupportedMethod(p *packet.Packet) {
	p.Response.StatusCode = http.StatusNotImplemented
	p.Response.StatusText = http.StatusText(p.Response.StatusCode)
}

// methodNotAllowed is a method to set the server response as HTTP 405.
func (s *Server) methodNotAllowed(p *packet.Packet) {
	p.Response.StatusCode = http.StatusMethodNotAllowed
	p.Response.StatusText = http.StatusText(p.Response.StatusCode)
}

// internalServerError is a method to set the response as HTTP 500 with a
// custom message.
func (s *Server) internalServerError(p *packet.Packet, msg string) {
	p.Response.StatusCode = http.StatusInternalServerError
	p.Response.StatusText = msg
}

// NewServer creates a new server handler to listen for incoming requests.
func NewServer(options ServerSetup, handler interface{}) (*Server, error) {
	if options.MediaSetup == nil {
		return nil, errors.New("no MediaSetup was found")
	}

	ports, err := adt.NewRangeBox(options.UDPPortMin, options.UDPPortMax)

	if err != nil {
		return nil, err
	}

	addr := fmt.Sprintf("0.0.0.0:%d", options.Port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)

	if err != nil {
		return nil, err
	}

	l, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		return nil, err
	}

	return &Server{
		ServerSetup:    options,
		rtspListener:   l,
		handler:        handler,
		shutdown:       make(chan bool),
		clientTimeout:  500,
		activeSessions: make(map[string]*rtp.Session),
		availablePorts: ports,
		session: sdp.NewSession(sdp.Setup{
			ClientHost: options.MediaSetup.ClientHost,
			Video: sdp.MediaSetup{
				Port: options.Port,
			},
		}),
	}, nil
}
