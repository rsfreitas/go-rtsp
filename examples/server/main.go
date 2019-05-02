//
// Description: A server example.
// Author: Rodrigo Freitas
// Created at: Tue Apr 23 09:25:33 -03 2019
//
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rsfreitas/go-rtsp"
)

type Handler struct{}

//func (h *Handler) Setup(w *rtsp.Response, r *rtsp.Request) {
//}

func (h *Handler) Play() {
	fmt.Println("Client Play")
}

func (h *Handler) Pause() {
	fmt.Println("Client Pause")
}

/*func (h *Handler) Record(w *rtsp.Response, r *rtsp.Request) {
}

func (h *Handler) Announce(w *rtsp.Response, r *rtsp.Request) {
}

func (h *Handler) Teardown(w *rtsp.Response, r *rtsp.Request) {
}

func (h *Handler) SetParameter(w *rtsp.Response, r *rtsp.Request) {
}

func (h *Handler) GetParameter(w *rtsp.Response, r *rtsp.Request) {
}*/

func monitorOurself(server *rtsp.Server) {
	quit := make(chan os.Signal)
	signal.Notify(quit,
		os.Interrupt,
		syscall.SIGTERM)

	go func() {
		<-quit
		fmt.Println("Finishing application")
		server.Close()
	}()
}

func main() {
	handler := &Handler{}

	server, err := rtsp.NewServer(
		rtsp.ServerSetup{
			Port:       8554,
			UDPPortMin: 39000,
			UDPPortMax: 45001,
			MediaSetup: &rtsp.MediaSetup{
				Port:       8001,
				ClientHost: "127.0.0.1",
			},
			//				Video: rtsp.MediaSetup{
			//					Port: 8001,
			//				},
			//				ClientHost: "127.0.0.1",
			//			},
		}, handler)

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	go monitorOurself(server)
	fmt.Println("Starting server")
	server.Start()
}
