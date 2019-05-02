//
// Description:
// Author: Rodrigo Freitas
// Created at: Tue Apr 23 15:29:24 -03 2019
//
package packet

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"

	"github.com/gortc/sdp"
)

type Packet struct {
	*Request
	*Response
}

func (p *Packet) parseBody(body []byte) error {
	var err error
	contentType, ok := p.Request.Headers["Content-Type"]

	if !ok {
		return errors.New("header without Content-Type")
	}

	switch contentType[0] {
	case "application/sdp":
		p.Request.SDP, err = sdp.DecodeSession(body, p.Request.SDP)

		if err != nil {
			return err
		}
	}

	return nil
}

func skipBlankLines(in *bufio.Reader) error {
	for {
		b, err := in.Peek(1)

		if err != nil {
			return err
		}

		if b[0] == '\r' || b[0] == '\n' {
			in.Discard(1)
		} else {
			break
		}
	}

	return nil
}

func (p *Packet) parseHeadersAndBody(in []byte, start, length int) error {
	var err error

	// Reads the request headers
	reader := bufio.NewReader(strings.NewReader(string(in[start:length]) + "\r\n"))
	tp := textproto.NewReader(reader)
	p.Request.Headers, err = tp.ReadMIMEHeader()

	if err != nil {
		return err
	}

	// Reads the request body (if we have one)
	if l, ok := p.Request.Headers["Content-Length"]; ok {
		length, err := strconv.ParseInt(l[0], 10, 32)

		if err != nil {
			return err
		}

		// Advances to the next valid character
		if err := skipBlankLines(reader); err != nil {
			return err
		}

		body := make([]byte, int(length))
		n, err := reader.Read(body)

		if err != nil {
			return err
		}

		if n != int(length) {
			return errors.New("read less bytes than Content-Length")
		}

		if err := p.parseBody(body); err != nil {
			return err
		}
	}

	return nil
}

func (p *Packet) StillNeedsRead(in []byte, length int) bool {
	s := in[length-4 : length]

	if s[0] == '\r' && s[1] == '\n' && s[2] == '\r' && s[3] == '\n' {
		return false
	}

	return true
}

func (p *Packet) UnmarshalRequest(in []byte, length int) error {
	var (
		err    error
		ok     bool
		offset int
	)

	fmt.Printf("--BEGIN--\n|%s|\n--END--\n", string(in))

method:
	for i := 0; i < length; i++ {
		switch in[i] {
		case ' ', '\t':
			p.Request.Method = string(in[0:i])
			ok = true
			offset = i + 1
			break method
		}
	}

	if !ok {
		return errors.New("missing method")
	}

	ok = false

url:
	for i := offset; i < length; i++ {
		switch in[i] {
		case ' ', '\t':
			p.Request.URL, err = url.Parse(string(in[offset:i]))

			if err != nil {
				return err
			}

			ok = true
			offset = i + 1
			break url
		}
	}

	if !ok {
		return errors.New("missing URL")
	}

	ok = false

version:
	for i, read := offset, false; i < length; i++ {
		c := in[i]

		switch read {
		case false:
			switch c {
			case '\r':
				p.Request.Version = string(in[offset:i])
				read = true

			case '\n':
				p.Request.Version = string(in[offset:i])
				offset = i + 1
				ok = true
				break version
			}

		case true:
			if c != '\n' {
				return errors.New("missing newline in version")
			}

			offset = i + 1
			ok = true
			break version
		}
	}

	if !ok {
		return errors.New("missing version")
	}

	if err := p.parseHeadersAndBody(in, offset, length); err != nil {
		return err
	}

	p.Request.loadInfoFromHeaders()

	return nil
}

func (p *Packet) MarshalResponse() ([]byte, error) {
	var b bytes.Buffer

	b.WriteString(fmt.Sprintf("%s %d %s\r\n", p.Request.Version,
		p.Response.StatusCode, p.Response.StatusText))

	// Cseq: uses the same Request Cseq if we previously received it
	b.WriteString(fmt.Sprintf("Cseq: %d\r\n", p.Request.sequence))

	// Response.Headers
	for k, v := range p.Response.Headers {
		b.WriteString(fmt.Sprintf("%s: %s\r\n", k, v[0]))
	}

	b.WriteString("\r\n")

	// Body
	if p.Response.Body != nil {
		b.WriteString(fmt.Sprintf("%s\r\n", string(p.Response.Body)))
	}

	b.WriteString("\r\n")

	return b.Bytes(), nil
}

func (p *Packet) MarshalResponseError(err error) []byte {
	// TODO
	return nil
}

// NewPacket creates a new Packet object to hold the client request and its
// response.
func NewPacket() *Packet {
	return &Packet{
		Request:  &Request{},
		Response: newResponse(),
	}
}
