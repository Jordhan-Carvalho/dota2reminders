package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jordhan-carvalho/belphegorv2/interfaces"
)

type GameEventsHandler struct {
	GameEventsChan chan interfaces.GameEvents
	VoiceStarted   *bool
  GameEventsReceivers *int
	// entirePayload interface{}
}


func (g *GameEventsHandler) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Println("Chegou o request")
	gameEvent := interfaces.GameEvents{}
	// logEntirePayload := g.entirePayload

	err := json.NewDecoder(r.Body).Decode(&gameEvent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if *g.VoiceStarted {
		// If we dont specify the receiver count the channel will be amepty after the fist consume
    fmt.Println("handlers.go, sending the gameEvent to channel, receivers:", *g.GameEventsReceivers)
		for i := 0; i < *g.GameEventsReceivers; i++ {
			g.GameEventsChan <- gameEvent
		}
	}

	fmt.Fprintf(w, "Game Event: %+v", gameEvent)
}
