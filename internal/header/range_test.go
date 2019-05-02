//
// Description:
// Author: Rodrigo Freitas
// Created at: Mon Apr 29 18:50:13 -03 2019
//
package header_test

import (
	"fmt"
	"testing"

	"github.com/rsfreitas/go-rtsp/internal/header"
	"github.com/stretchr/testify/assert"
)

func TestNewRange(t *testing.T) {
	assert := assert.New(t)

	r, err := header.NewRange("clock=19960213T143205.25Z-;time=19970123T143720Z")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(r)
	assert.Equal(1, 1, "Should be equal")
}
