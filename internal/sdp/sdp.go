//
// Description:
// Author: Rodrigo Freitas
// Created at: Mon Apr 29 13:43:19 -03 2019
//
package sdp

import (
	"net"

	"github.com/gortc/sdp"
)

type MediaSetup struct {
	Port int
}

type Setup struct {
	Video      MediaSetup
	ClientHost string

	message *sdp.Message
}

type Session struct {
	Setup

	session sdp.Session
}

func (s *Session) Bytes() []byte {
	var b []byte

	return s.session.AppendTo(b)
}

func NewSession(options Setup) *Session {
	// video
	video := sdp.Media{
		Description: sdp.MediaDescription{
			Type:     "video",
			Port:     options.Video.Port,
			Formats:  []string{"99"},
			Protocol: "RTP/AVP",
		},
	}

	video.AddAttribute("rtpmap", "99", "h263-1998/90000")

	// message
	message := &sdp.Message{
		Origin: sdp.Origin{
			Address: "127.0.0.1",
		},
		Connection: sdp.ConnectionData{
			IP: net.ParseIP(options.ClientHost),
		},
		Name:   "video forwarding",
		Medias: []sdp.Media{video},
	}

	// session
	var ss sdp.Session
	ss = message.Append(ss)

	return &Session{
		session: ss,
		Setup:   options,
	}
}
