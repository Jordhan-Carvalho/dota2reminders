package server

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/jordhan-carvalho/belphegorv2/interfaces"
)

type GameEventsHandler struct {
	GameEventsChan chan interfaces.GameEvents
	VoiceStarted   *bool
  GameEventsReceivers *int
	// entirePayload interface{}
}


func (g *GameEventsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

  getIP(r)

	gameEvent := interfaces.GameEvents{}
	// logEntirePayload := g.entirePayload

	err := json.NewDecoder(r.Body).Decode(&gameEvent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if *g.VoiceStarted {
		// If we dont specify the receiver count the channel will be empty after the fist consume
    fmt.Println("handlers.go, sending the gameEvent to channel, receivers:", *g.GameEventsReceivers)
		for i := 0; i < *g.GameEventsReceivers; i++ {
			g.GameEventsChan <- gameEvent
		}
	}

	fmt.Fprintf(w, "Game Event: %+v", gameEvent)
}


// https://blog.golang.org/context/userip/userip.go
func getIP(req *http.Request) {
	fmt.Println("Olha req.RemoteAddr", req.RemoteAddr)
	ip, port, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		//return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)

		fmt.Printf("userip: %q is not IP:port", req.RemoteAddr)
	}

	userIP := net.ParseIP(ip)
	fmt.Println("User IP after PArseIp", userIP)
	if userIP == nil {
		//return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
		fmt.Printf("userip: %q is not IP:port", req.RemoteAddr)
		return
	}

	// This will only be defined when site is accessed via non-anonymous proxy
	// and takes precedence over RemoteAddr
	// Header.Get is case-insensitive
	forward := req.Header.Get("X-Forwarded-For")

	fmt.Printf("<p>IP: %s</p>", ip)
	fmt.Printf("<p>Port: %s</p>", port)
	fmt.Printf("<p>Forwarded for: %s</p>", forward)
}
