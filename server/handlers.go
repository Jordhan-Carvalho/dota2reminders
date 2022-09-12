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
}


func (g *GameEventsHandler) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Println("Chegou o request")
	gameEvent := interfaces.GameEvents{}

	err := json.NewDecoder(r.Body).Decode(&gameEvent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if *g.VoiceStarted {
    fmt.Println("Started to send request to channel")
    // THIS WILL BLOCK THE EXECUTION.... WHY?????????
		g.GameEventsChan <- gameEvent
    fmt.Println("After send evnet to game chan")
	}

	fmt.Fprintf(w, "Game Event: %+v", gameEvent)
	//TODO REMOVE THIS, I THINK ITS BETTER TO JUST LISTEN TO THE MESSAGES... HOW WOULD YOU JOIN A SERVER IF THERE NO DISCORDSERVER
	// WHEN YOU TYPE START, JOIN CHANNEL AND STARTS THE TIMER WITH THE CLOCKTIME
	/* if gameEvent.Map.GameState == "DOTA_GAMERULES_STATE_GAME_IN_PROGRESS" {

	} */

}
