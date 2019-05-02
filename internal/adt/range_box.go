//
// Description:
// Author: Rodrigo Freitas
// Created at: Tue Apr 30 13:44:12 -03 2019
//
package adt

import (
	"errors"
	"sync"
)

// RangeBox is a holder for controlling (requesting and releasing) a specific
// interval of integer values. A requested value is always reserved in pairs.
type RangeBox struct {
	lock     sync.RWMutex
	min      uint32
	max      uint32
	current  uint32
	occupied map[uint32]uint32
	capacity uint32
}

// Min gives the minimum value of a RangeBox
func (r *RangeBox) Min() uint32 {
	return r.min
}

// Max gives the maximum value of a RangeBox
func (r *RangeBox) Max() uint32 {
	return r.max
}

// Capacity gives the RangeBox internal limit of simultaneous requests
func (r *RangeBox) Capacity() uint32 {
	return r.capacity
}

// Request requests a value from a RangeBox. Internally the requested value
// will occupy two resources, so is correct to use the return value and it
// plus 1 as valid values.
func (r *RangeBox) Request() (uint32, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.current == 0 {
		r.current = r.min
		r.occupied[r.current] = r.current + 1
		r.occupied[r.current+1] = r.current

		return r.current, nil
	}

	if uint32(len(r.occupied)) == r.capacity {
		return 0, errors.New("no port is available")
	}

	next := r.current + 2

	for {
		if next > r.max {
			next = r.min
		}

		if _, ok := r.occupied[next]; ok {
			next += 2
			continue
		}

		r.current = next
		r.occupied[r.current] = r.current + 1
		r.occupied[r.current+1] = r.current

		return r.current, nil
	}

	return 0, errors.New("no port is available")
}

// Release releases a previously requested value from a RangeBox to be
// requested again.
func (r *RangeBox) Release(port uint32) {
	r.lock.Lock()
	defer r.lock.Unlock()

	next, ok := r.occupied[port]

	if !ok {
		return
	}

	delete(r.occupied, port)
	delete(r.occupied, next)
}

// Occupied tests if port is currently requested.
func (r *RangeBox) Occupied(port uint32) bool {
	r.lock.Lock()
	defer r.lock.Unlock()

	_, ok := r.occupied[port]

	return ok
}

// NewRangeBox creates a new RangeBox.
func NewRangeBox(min, max uint32) (*RangeBox, error) {
	if min >= max {
		return nil, errors.New("min must be less than max")
	}

	capacity := max - min + 1

	if capacity%2 != 0 {
		return nil, errors.New("the port interval should be even")
	}

	return &RangeBox{
		min:      min,
		max:      max,
		occupied: make(map[uint32]uint32),
		capacity: capacity,
	}, nil
}
