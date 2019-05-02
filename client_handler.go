//
// Description:
// Author: Rodrigo Freitas
// Created at: Mon Apr 22 17:18:38 -03 2019
//
package rtsp

type ClientPlay interface {
	Play()
}

type ClientPause interface {
	Pause()
}

type ClientTeardown interface {
	Teardown()
}

type ClientRecord interface {
	Record()
}

type ClientAnnounce interface {
	Announce()
}

type ClientGetParameter interface {
	GetParameter()
}

type ClientSetParameter interface {
	SetParameter()
}
