package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	router.HandleFunc("/api/heroes/{id:[0-9]+}", getHeroByID)

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
	case "POST":
		addHero(w, r)
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

func getHeroByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	varID := vars["id"]
	i, err := strconv.Atoi(varID)

	if err != nil {
		fmt.Fprintf(w, "Err during string convert: %v", err)
		return
	}

	b, err := json.Marshal(Heroes[i-1])

	if err != nil {
		fmt.Fprintf(w, "Err by marshal hero: %v", err)
		return
	}

	writeToClient(w, string(b))
}

func updateHero(w http.ResponseWriter, r *http.Request) {
	hero, err := getHeroFromRequest(r, w)

	if err != nil {
		fmt.Fprintf(w, "Err by getHeroFromRequest %v", err)
		return
	}

	i, err := strconv.Atoi(hero.ID)

	if err != nil {
		fmt.Fprintf(w, "Err during string convert: %v", err)
		return
	}

	//update Hero in List
	Heroes[i-1] = hero

	writeToClient(w, "")
}

func addHero(w http.ResponseWriter, r *http.Request) {
	hero, err := getHeroFromRequest(r, w)

	if err != nil {
		fmt.Fprintf(w, "Err by getHeroFromRequest %v", err)
		return
	}

	hero.ID = strconv.Itoa(len(Heroes) + 1)

	Heroes = append(Heroes, hero)

	b, err := json.Marshal(hero)

	if err != nil {
		fmt.Fprintf(w, "Err by marshal hero: %v", err)
		return
	}

	writeToClient(w, string(b))
}

func getHeroFromRequest(r *http.Request, w http.ResponseWriter) (Hero, error) {
	//read transfered data
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		return Hero{}, err
	}

	//convert to Hero
	var hero Hero
	err = json.Unmarshal(body, &hero)

	if err != nil {
		return Hero{}, err
	}

	return hero, nil
}

func writeToClient(w http.ResponseWriter, s string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	fmt.Fprintf(w, s)
}
