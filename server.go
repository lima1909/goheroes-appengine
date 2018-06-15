package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

// Hero type
type Hero struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var (
	// Heroes list from Hero examples
	Heroes = []Hero{
		Hero{"1", "Jasmin"},
		Hero{"2", "Mario"},
		Hero{"3", "Alex M"},
		Hero{"4", "Adam O"},
		Hero{"5", "Shauna C"},
		Hero{"6", "Lena H"},
		Hero{"7", "Chris S"},
	}
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/heroes", heros)
	router.HandleFunc("/api/heroes/{id:[0-9]+}", getHeroById)

	http.Handle("/", router)

	log.Println("Start Server: http://localhost:8081")
	log.Fatalln(http.ListenAndServe(":8081", nil))
}

func heros(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(Heroes)
	if err != nil {
		fmt.Fprintf(w, "Err by marshal heros: %v", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	fmt.Fprintln(w, string(b))
}

func getHeroById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	varId := vars["id"]  
	i, err := strconv.Atoi(varId)

	if err != nil {
		fmt.Fprintf(w, "Err during string convert: %v", err)
		return
	}

	b, err := json.Marshal(Heroes[i-1]);

	if err != nil {
		fmt.Fprintf(w, "Err by marshal hero: %v", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	fmt.Fprintf(w, string(b))
}

/* func updateHero(w http.ResponseWriter, r *http.Request) {
	idString := r.URL.Path[len("/api/heroes/"):]
	i, err := strconv.Atoi(idString)

	if err != nil {
		fmt.Fprintf(w, "Err during string converge: %v", err)
		return
	}

	if r.Method == "PUT" {
		fmt.Fprintf(w, "method PUT")
		fmt.Fprintf(w, string(r.Body))
	}


	b, err := json.Marshal(Heroes[i-1]);

	// vars := mux.Vars(r)
	// varId := vars["id"]  
	// i, err := strconv.Atoi(varId)

	// if err != nil {
	// 	fmt.Fprintf(w, "Err during string convert: %v", err)
	// 	return
	// }

	//b, err := json.Marshal(Heroes[i-1]);

	if err != nil {
		fmt.Fprintf(w, "Err by marshal hero: %v", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	fmt.Fprintf(w, string(b))
} */


