package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type whatcomes struct {
	Test string
  aff string
}

func main() {
	http.HandleFunc("/", gameEventsHandler)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func gameEventsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	a := whatcomes{}

	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

  fmt.Println("Printandoooo", a.Test)
	// Do something with the Person struct...
	fmt.Fprintf(w, "Person: %+v", a)

}
