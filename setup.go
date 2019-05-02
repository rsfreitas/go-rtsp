//
// Description:
// Author: Rodrigo Freitas
// Created at: Sun Apr 28 10:00:41 -03 2019
//
package rtsp

import (
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/rsfreitas/go-rtsp/internal/header"
	"github.com/rsfreitas/go-rtsp/internal/packet"
	"github.com/rsfreitas/go-rtsp/internal/rtp"
)

type setupMethod struct {
	ActiveSessions map[string]*rtp.Session
	ServerPortMin  int
	ServerPortMax  int
}

func (s *setupMethod) Verify(p *packet.Packet, handler interface{}) error {
	return nil
}

func (s *setupMethod) Handle(p *packet.Packet) {
	var (
		transport     *header.Transport
		clientSession string
		session       *rtp.Session
		err           error
	)

	// We look for a Session identification inside the Request Headers so
	// we can change it's parameters.
	if field, ok := p.Request.Headers["Session"]; ok {
		if s.ActiveSessions == nil {
			p.Response.StatusCode = StatusAggregateOperationNotAllowed
			p.Response.StatusText = StatusText(p.Response.StatusCode)
			return
		}

		clientSession = field[0]
		session, ok = s.ActiveSessions[clientSession]

		if !ok {
			p.Response.StatusCode = StatusSessionNotFound
			p.Response.StatusText = StatusText(p.Response.StatusCode)
			return
		}
	}

	// We also get the Transport request header parameter to gather more
	// information about what we're going to stream.
	if field, ok := p.Request.Headers["Transport"]; !ok {
		p.Response.StatusCode = StatusUnsupportedTransport
		p.Response.StatusText = StatusText(p.Response.StatusCode)
		return
	} else {
		transport, err = header.NewTransportFromString(field[0])

		if err != nil {
			p.Response.StatusCode = StatusUnsupportedTransport
			p.Response.StatusText = StatusText(p.Response.StatusCode)
			return
		}
	}

	if session == nil {
		// Creates a new RTP session to transfer data to client.
		u, err := uuid.NewV4()

		if err != nil {
			p.Response.StatusCode = http.StatusInternalServerError
			p.Response.StatusText = "Unable to create new session"
			return
		}

		var options rtp.Setup

		switch transport.LowerTransport {
		case "TCP":
			p.Response.StatusCode = StatusUnsupportedTransport
			p.Response.StatusText = StatusText(p.Response.StatusCode)
			return

		default:
			if !transport.HasParameter("client_port") {
				p.Response.StatusCode = StatusParameterNotUnderstood
				p.Response.StatusText = StatusText(p.Response.StatusCode)
				return
			}

			options = rtp.Setup{
				ServerPort:  s.ServerPortMin,
				ServerAddr:  "127.0.0.1",
				ClientAddr:  "127.0.0.1", // FIXME
				ClientPorts: transport.ClientPort,
			}
		}

		s.ActiveSessions[u.String()] = rtp.NewSession(options)
		p.Response.Headers.Add("Session", u.String())
	} else {
		// We still don't handle a session update coming from the client
		// side. So we return an error for that.
		p.Response.StatusCode = StatusAggregateOperationNotAllowed
		p.Response.StatusText = StatusText(p.Response.StatusCode)
		return
	}

	p.Response.Headers.Add("Transport", s.transportHeader(p, transport))
	p.Response.StatusCode = http.StatusOK
	p.Response.StatusText = http.StatusText(http.StatusOK)
}

func (s *setupMethod) transportHeader(p *packet.Packet, t *header.Transport) string {
	serverTransport := header.NewTransport()

	if t.HasParameter("client_port") {
		var d []interface{} = make([]interface{}, len(t.ClientPort))

		for i, n := range t.ClientPort {
			d[i] = n
		}

		serverTransport.AppendParameter("client_port", d...)
	}

	serverTransport.SetDelivery(header.TransportUnicast)
	serverTransport.SetTransport(header.TransportRTP)
	serverTransport.SetProfile()
	serverTransport.SetLowerTransport(header.TransportLowerUDP)
	serverTransport.AppendParameter("server_port", s.ServerPortMin, s.ServerPortMax)
	serverTransport.AppendParameter("mode", "PLAY")

	return serverTransport.String()
}

func (s *setupMethod) Type() methodType {
	return methodSetup
}
