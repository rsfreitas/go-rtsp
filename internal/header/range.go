//
// Description:
// Author: Rodrigo Freitas
// Created at: Wed May  1 21:04:07 -03 2019
//
package header

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type NptType int

const (
	NptNow NptType = iota + 1
	NptSec
	NptHHMMSS
)

type Npt struct {
	Type     NptType
	Hours    int
	Minutes  int
	Seconds  int
	Fraction int
}

type SmpteType int

const (
	SmpteSimple SmpteType = iota + 1
	Smpte30Drop
	Smpte25
)

type SmpteTime struct {
	Hours     int
	Minutes   int
	Seconds   int
	Frames    int
	Subframes int
}

type Smpte struct {
	Type SmpteType
	Time []SmpteTime
}

type Range struct {
	parameters map[string][]string

	Npt   []*Npt
	Smpte *Smpte
	Clock []time.Time
	Time  time.Time
}

func (r *Range) String() string {
	var s string

	for k, v := range r.parameters {
		if len(v) == 1 {
			s += fmt.Sprintf("%s=%s", k, v[0])
		} else {
			s += fmt.Sprintf("%s=%s-%s", k, v[0], v[1])
		}

		s += ";"
	}

	return strings.TrimSuffix(s, ";")
}

func newSmpteTime(s string) (SmpteTime, error) {
	fields := strings.Split(s, ":")

	if len(fields) == 0 {
		return SmpteTime{}, errors.New("invalid 'smpte' field")
	}

	var err error
	st := SmpteTime{}
	d := make(map[int]string)

	for i, j := range fields {
		d[i] = j
	}

	if t, ok := d[0]; ok {
		st.Hours, err = strconv.Atoi(t)

		if err != nil {
			return SmpteTime{}, errors.New("invalid 'smpte' hours")
		}
	}

	if t, ok := d[1]; ok {
		st.Minutes, err = strconv.Atoi(t)

		if err != nil {
			return SmpteTime{}, errors.New("invalid 'smpte' minutes")
		}
	}

	if t, ok := d[2]; ok {
		st.Seconds, err = strconv.Atoi(t)

		if err != nil {
			return SmpteTime{}, errors.New("invalid 'smpte' seconds")
		}
	}

	if t, ok := d[3]; ok {
		if strings.Contains(t, ".") {
			tt := strings.Split(t, ".")
			st.Frames, err = strconv.Atoi(tt[0])

			if err != nil {
				return SmpteTime{}, errors.New("invalid 'smpte' frames")
			}

			st.Subframes, err = strconv.Atoi(tt[1])

			if err != nil {
				return SmpteTime{}, errors.New("invalid 'smpte' subframes")
			}
		} else {
			st.Frames, err = strconv.Atoi(t)

			if err != nil {
				return SmpteTime{}, errors.New("invalid 'smpte' frames")
			}
		}
	}

	return st, nil
}

func parseSeconds(s string) (int, int, error) {
	var (
		seconds  int
		fraction int
		err      error
	)

	if strings.Contains(s, ".") {
		t := strings.Split(s, ".")
		seconds, err = strconv.Atoi(t[0])

		if err != nil {
			return 0, 0, errors.New("invalid 'npt' seconds field")
		}

		fraction, err = strconv.Atoi(t[1])

		if err != nil {
			return 0, 0, errors.New("invalid 'npt' seconds fraction field")
		}
	} else {
		seconds, err = strconv.Atoi(s)

		if err != nil {
			return 0, 0, errors.New("invalid 'npt' seconds field")
		}
	}

	return seconds, fraction, nil
}

func newNpt(s string) (*Npt, error) {
	if s == "now" {
		return &Npt{
			Type: NptNow,
		}, nil
	}

	var err error
	n := &Npt{}

	if strings.Contains(s, ":") {
		n.Type = NptHHMMSS
		t := strings.Split(s, ":")
		n.Hours, err = strconv.Atoi(t[0])

		if err != nil {
			return nil, errors.New("invalid 'npt' hours field")
		}

		n.Minutes, err = strconv.Atoi(t[1])

		if err != nil {
			return nil, errors.New("invalid 'npt' minutes field")
		}

		n.Seconds, n.Fraction, err = parseSeconds(t[2])

		if err != nil {
			return nil, err
		}
	} else {
		n.Type = NptSec
		n.Seconds, n.Fraction, err = parseSeconds(s)

		if err != nil {
			return nil, err
		}
	}

	return n, nil
}

func NewRange(in string) (*Range, error) {
	r := &Range{
		parameters: make(map[string][]string),
	}

	for _, p := range strings.Split(in, ";") {
		if strings.Contains(p, "=") {
			var n []string
			t := strings.Split(p, "=")

			if strings.Contains(t[1], "-") {
				n = strings.Split(t[1], "-")
			} else {
				n = []string{t[1]}
			}

			r.parameters[t[0]] = n
		}
	}

	if f, ok := r.parameters["time"]; ok {
		t, err := time.Parse("20060102T150405Z", f[0])

		if err != nil {
			return nil, err
		}

		r.Time = t
	}

	if f, ok := r.parameters["clock"]; ok {
		if len(f) != 2 {
			return nil, errors.New("invalid 'clock' field")
		}

		for _, s := range f {
			if t, err := time.Parse("20060102T150405.99Z", s); err == nil {
				r.Clock = append(r.Clock, t)
			}
		}
	}

	var (
		f  []string
		ok bool
	)

	smpte := &Smpte{}
	f, ok = r.parameters["smpte"]

	if !ok {
		f, ok = r.parameters["smpte-30-drop"]

		if !ok {
			f, ok = r.parameters["smpte-25"]

			if ok {
				smpte.Type = Smpte25
			}
		} else {
			smpte.Type = Smpte30Drop
		}
	} else {
		smpte.Type = SmpteSimple
	}

	if ok {
		for _, s := range f {
			st, err := newSmpteTime(s)

			if err != nil {
				return nil, err
			}

			smpte.Time = append(smpte.Time, st)
		}

		r.Smpte = smpte
	}

	if f, ok := r.parameters["npt"]; ok {
		var npt []*Npt

		for _, s := range f {
			n, err := newNpt(s)

			if err != nil {
				return nil, err
			}

			npt = append(npt, n)
		}

		r.Npt = npt
	}

	return r, nil
}
