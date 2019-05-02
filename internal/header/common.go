//
// Description:
// Author: Rodrigo Freitas
// Created at: Thu May  2 09:02:39 -03 2019
//
package header

import (
	"strconv"
	"strings"
)

// parseInt parses a string into a integer value.
func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// parseIntSlice parses a string, splitting it using a delimiter, into a
// slice of integers.
func parseIntSlice(s, delim string) ([]int, error) {
	var (
		v   int
		n   []int
		err error
	)

	if strings.Contains(s, delim) {
		for _, i := range strings.Split(s, delim) {
			v, err = strconv.Atoi(i)

			if err != nil {
				return nil, err
			}

			n = append(n, v)
		}
	} else {
		v, err = strconv.Atoi(s)

		if err != nil {
			return nil, err
		}

		n = append(n, v)
	}

	return n, nil
}
