package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Hero type
type Hero struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var (
	// Heros list from Hero examples
	Heros = []Hero{
		Hero{"1", "Jasmin"},
		Hero{"2", "Mario"},
		Hero{"3", "Alex M"},
		Hero{"4", "Adam O"},
		Hero{"5", "Shauna C"},
	}
)

func main() {
	http.HandleFunc("/api/heroes", heros)
	log.Println("Start Server: http://localhost:8081")
	log.Fatalln(http.ListenAndServe(":8081", nil))
}

func heros(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(Heros)
	if err != nil {
		fmt.Fprintf(w, "Err by marshal heros: %v", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	fmt.Fprintln(w, string(b))
}
