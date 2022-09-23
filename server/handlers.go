package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jordhan-carvalho/belphegorv2/interfaces"
	"github.com/jordhan-carvalho/belphegorv2/utils"
)

type GameEventsHandler struct {
	GameEventsChan      chan interfaces.GameEvents
	VoiceStarted        *bool
	GameEventsReceivers *int
	// entirePayload interface{}
}

func (g *GameEventsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	userIP := utils.GetIP(r)
	persistedStatus := getIPPersistence()
	if !shouldListenToRequest(userIP, persistedStatus) {
		http.Error(w, "Listening to a game.", http.StatusBadRequest)
		return
	}

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

type IPStatus struct {
	LastActiveIp string `json:"lastActiveIP"`
	LastReqTime  int64  `json:"lastReqTime"`
}

func writteIPPersistence(userIP string) {
	now := time.Now().Unix()

	ipStatus := &IPStatus{
		LastActiveIp: userIP,
		LastReqTime:  now,
	}

	content, err := json.Marshal(ipStatus)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("ippersistence.json", content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func getIPPersistence() IPStatus {
	content, err := ioutil.ReadFile("ippersistence.json")
	if err != nil {
		log.Fatal(err)
	}
	ipStatus := IPStatus{}
	err = json.Unmarshal(content, &ipStatus)
	if err != nil {
		log.Fatal(err)
	}

	return ipStatus
}

func shouldListenToRequest(userIP string, ps IPStatus) bool {
	now := time.Now().Unix()
	secondsToInvalidateGame := 180
	secondsToRewritte := 30

	if userIP == ps.LastActiveIp {
		if now >= (ps.LastReqTime + int64(secondsToRewritte)) {
			writteIPPersistence(userIP)
		}
		return true
	} else {
		if now > ps.LastReqTime+int64(secondsToInvalidateGame) {
			writteIPPersistence(userIP)
			return true
		}
	}
	return false
}
