package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type mapEvent struct {
  name string
}

type gameEvents struct {
	Map mapEvent
}

func main() {
	http.HandleFunc("/", gameEventsHandler)

	fmt.Printf("Starting server at port 3000\n")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}

func gameEventsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

  fmt.Println("Chegou o request")
	a := gameEvents{}

	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

  json.NewEncoder(os.Stdout).Encode(&m)
  fmt.Println("Printandoooo", a.Map.name)
	// Do something with the Person struct...
	fmt.Fprintf(w, "Game Event: %+v", a)

}
