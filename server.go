package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"errors"
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
	router.HandleFunc("/api/heroes", heros)

	http.Handle("/", router)

	log.Println("Start Server: http://localhost:8081")
	log.Fatalln(http.ListenAndServe(":8081", nil))
}

func heros(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
		case "GET":
			loadHeroes(w)
		case "OPTIONS":
			writeToClient(w, string(http.StatusOK))
		case "PUT":
			updateHero(w, r)
	}
}

func loadHeroes(w http.ResponseWriter) {
	b, err := json.Marshal(Heroes)

	if err != nil {
		fmt.Fprintf(w, "Err by marshal heros: %v", err)
		return
	}

	writeToClient(w, string(b))
}

func getHeroById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	varId := vars["id"]  

	i, err := getIndex(varId)

	if err != nil {
		fmt.Fprintf(w, "Err by getIndex ", err)
		return
	}

	b, err := json.Marshal(Heroes[i]);

	if err != nil {
		fmt.Fprintf(w, "Err by marshal hero: %v", err)
		return
	}

	writeToClient(w, string(b))
}

func updateHero(w http.ResponseWriter, r *http.Request) {
	//read transfered data
	body, err := ioutil.ReadAll(r.Body)
	
	if err != nil {
        fmt.Fprintf(w, "can't read body", err)
        return
	}
	
	//convert to Hero
	var hero Hero
	err = json.Unmarshal(body, &hero)

	if err != nil {
		fmt.Fprintf(w, "Err by unmarshal hero: %v", err)
		return
	}

	//update Hero in List
	i, err := getIndex(hero.ID)

	if err != nil {
		fmt.Fprintf(w, "Err by getIndex ", err)
		return
	}

	Heroes[i] = hero;

	writeToClient(w, "")
} 

//don't know if the ids are really the same as the indices in the future!
func getIndex(id string) (int, error) {
	for index, h := range Heroes {
		if (h.ID == id) {
			return index, nil;
		}
	}	
	return 0, errors.New("No matching id found for "+id)
}

func writeToClient(w http.ResponseWriter, s string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	fmt.Fprintf(w, s)
}


