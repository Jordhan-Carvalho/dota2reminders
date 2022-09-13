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
    fmt.Println("Started to send request to channel")
		g.GameEventsChan <- gameEvent
	}

	fmt.Fprintf(w, "Game Event: %+v", gameEvent)
}
