//
// Description:
// Author: Rodrigo Freitas
// Created at: Wed May  1 21:04:07 -03 2019
//
package header

import (
	"errors"
	"fmt"
	"strings"
)

type TransportDeliveryType int

const (
	TransportUnicast TransportDeliveryType = iota + 1
	TransportMulticast
)

type TransportMode int

const (
	TransportRTP TransportMode = iota + 1
)

type TransportLowerType int

const (
	TransportLowerUDP TransportLowerType = iota + 1
	TransportLowerTCP
)

// Transport holds all RTSP Transport header parameter already parsed.
type Transport struct {
	Transport      string
	Profile        string
	LowerTransport string
	Delivery       TransportDeliveryType
	Destination    string
	Source         string
	Layers         int
	Mode           string
	Append         bool
	Interleaved    []int
	TTL            int
	RTPPort        []int
	ClientPort     []int
	ServerPort     []int
	Ssrc           string

	parameters       map[string]string
	singleParameters []string
}

func (t *Transport) SetDelivery(mode TransportDeliveryType) {
	switch mode {
	case TransportUnicast:
		t.singleParameters = append(t.singleParameters, "unicast")

	case TransportMulticast:
		t.singleParameters = append(t.singleParameters, "multicast")
	}
}

func (t *Transport) SetTransport(transport TransportMode) error {
	switch transport {
	case TransportRTP:
		t.Transport = "RTP"

	default:
		return errors.New("unsupported transport mode")
	}

	return nil
}

func (t *Transport) SetProfile() {
	t.Profile = "AVP"
}

func (t *Transport) SetLowerTransport(lowerTransport TransportLowerType) {
	switch lowerTransport {
	case TransportLowerUDP:
		t.LowerTransport = "UDP"

	case TransportLowerTCP:
		t.LowerTransport = "TCP"
	}
}

// String returns the Transport object in the format required by the
// RTSP Transport header parameter.
func (t *Transport) String() string {
	var p string

	p += t.Transport

	if t.Profile != "" {
		p += "/" + t.Profile
	}

	if t.LowerTransport != "" {
		p += "/" + t.LowerTransport
	}

	s := fmt.Sprintf("%s", p)

	for _, v := range t.singleParameters {
		s += fmt.Sprintf(";%s", v)
	}

	for k, v := range t.parameters {
		s += fmt.Sprintf(";%s=%s", k, v)
	}

	return s
}

func (t *Transport) HasParameter(parameter string) bool {
	_, ok := t.parameters[parameter]

	return ok
}

func (t *Transport) AppendSingleParameter(parameter string) {
	t.singleParameters = append(t.singleParameters, parameter)
}

func (t *Transport) AppendParameter(key string, values ...interface{}) error {
	if len(values) == 0 {
		return errors.New("invalid parameter value")
	}

	if len(values) == 1 {
		switch values[0].(type) {
		case string:
			t.parameters[key] = values[0].(string)

		default:
			t.parameters[key] = fmt.Sprintf("%d", values[0].(int))
		}
	} else {
		var s string

		switch values[0].(type) {
		case string:
			s += fmt.Sprintf("%s-%s", values[0].(string), values[1].(string))

		case int:
			s += fmt.Sprintf("%d-%d", values[0].(int), values[1].(int))
		}

		t.parameters[key] = s
	}

	return nil
}

// NewTransport creates a new empty Transport object.
func NewTransport() *Transport {
	return &Transport{
		parameters: make(map[string]string),
	}
}

// NewTransportFromString parses a string, in the RTSP Transport header format,
// to a Transport object.
func NewTransportFromString(s string) (*Transport, error) {
	var err error

	t := NewTransport()

	for _, p := range strings.Split(s, ";") {
		if strings.Contains(p, "/") {
			f := strings.Split(p, "/")
			t.Transport = f[0]

			if len(f) > 1 {
				t.Profile = f[1]

				if len(f) > 2 {
					t.LowerTransport = f[2]
				} else {
					// Since we received nothing, set as UDP as our default lower
					// transport protocol
					t.LowerTransport = "UDP"
				}
			}
		} else if strings.Contains(p, "=") {
			f := strings.Split(p, "=")
			t.parameters[f[0]] = f[1]
		} else {
			t.singleParameters = append(t.singleParameters, p)
		}
	}

	if f, ok := t.parameters["port"]; ok {
		t.RTPPort, err = parseIntSlice(f, "-")

		if err != nil {
			return nil, err
		}
	}

	if f, ok := t.parameters["client_port"]; ok {
		t.ClientPort, err = parseIntSlice(f, "-")

		if err != nil {
			return nil, err
		}
	}

	if f, ok := t.parameters["server_port"]; ok {
		t.ServerPort, err = parseIntSlice(f, "-")

		if err != nil {
			return nil, err
		}
	}

	if f, ok := t.parameters["interleaved"]; ok {
		t.Interleaved, err = parseIntSlice(f, "-")

		if err != nil {
			return nil, err
		}
	}

	if f, ok := t.parameters["ttl"]; ok {
		t.TTL, err = parseInt(f)

		if err != nil {
			return nil, err
		}
	}

	if f, ok := t.parameters["layers"]; ok {
		t.Layers, err = parseInt(f)

		if err != nil {
			return nil, err
		}
	}

	if f, ok := t.parameters["destination"]; ok {
		t.Destination = f
	}

	if f, ok := t.parameters["ssrc"]; ok {
		t.Ssrc = f
	}

	if f, ok := t.parameters["mode"]; ok {
		t.Mode = f
	}

	for _, v := range t.singleParameters {
		if v == "append" {
			t.Append = true
		} else if v == "unicast" {
			t.Delivery = TransportUnicast
		} else if v == "multicast" {
			t.Delivery = TransportMulticast
		}
	}

	return t, nil
}
