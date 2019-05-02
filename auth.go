//
// Description:
// Author: Rodrigo Freitas
// Created at: Sun Apr 28 16:23:16 -03 2019
//
package rtsp

type AuthorizationType int

const (
	AuthorizationUnused AuthorizationType = iota + 1
	AuthorizationBasic
	AuthorizationDigest
)
