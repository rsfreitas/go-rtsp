//
// Description:
// Author: Rodrigo Freitas
// Created at: Tue Apr 23 21:22:48 -03 2019
//
package packet

import (
	"net/url"
	"strconv"

	"github.com/gortc/sdp"
)

type Request struct {
	URL     *url.URL
	Version string
	Method  string
	Headers map[string][]string
	SDP     sdp.Session

	sequence uint64
}

// loadInfoFromHeaders is where internal Request information are loaded
// according to internal Header fields.
func (r *Request) loadInfoFromHeaders() {
	if s, ok := r.Headers["Cseq"]; ok {
		if seq, err := strconv.ParseUint(s[0], 10, 64); err == nil {
			r.sequence = seq
		}
	}
}
