package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lima1909/goheroes-appengine/db"
	"github.com/lima1909/goheroes-appengine/service"

	"google.golang.org/appengine"
)

var (
	app  *service.App
	info Info
)

// Info to the current system
type Info struct {
	HeroesService      string
	EnvHeroServiceImpl string
}

// handle CORS and the OPION method
func corsAndOptionHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			h.ServeHTTP(w, r)
		}
	}
}

// create all used Handler
func handler() http.Handler {
	router := mux.NewRouter()

	router.Handle("/", http.RedirectHandler("/info", http.StatusFound))
	router.HandleFunc("/info", infoPage)

	router.HandleFunc("/api/heroes", addHero).Methods("POST")

	router.HandleFunc("/api/heroes", heroes)
	router.HandleFunc("/api/heroes/", searchHeroes)
	router.HandleFunc("/api/heroes/{id:[0-9]+}", heroID)

	return corsAndOptionHandler(router)
}

func init() {

	http.Handle("/", handler())

	app = service.NewApp(db.NewMemService())
	info = Info{HeroesService: "MemService", EnvHeroServiceImpl: "Not use in the moment"}
	log.Println("Init is ready and start the server on: http://localhost:8080")
}

func main() {
	appengine.Main()
}

func infoPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/info.html")
	if err != nil {
		fmt.Fprintf(w, "Err: %v\n", err)
		return
	}

	err = t.Execute(w, info)
	if err != nil {
		fmt.Fprintf(w, "Err: %v\n", err)
		return
	}
}

func heroes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		heroList(w, r)
	// case "POST":
	// 	addHero(w, r)
	case "PUT":
		updateHero(w, r)

	default:
		http.Error(w, "invalid method: "+r.Method, http.StatusBadRequest)
	}
}

func heroID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getHero(w, r)
	case "DELETE":
		deleteHero(w, r)

	default:
		http.Error(w, "invalid method: "+r.Method, http.StatusBadRequest)
	}
}

func heroList(w http.ResponseWriter, r *http.Request) {
	heroes, err := app.List(appengine.NewContext(r), "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(heroes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", string(b))
}

func addHero(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	hero := service.Hero{}
	err := json.NewDecoder(r.Body).Decode(&hero)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h, err := app.Add(appengine.NewContext(r), hero)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(h)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", string(b))
}

func getHero(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	varID := vars["id"]
	id, err := strconv.Atoi(varID)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid id: %v in params: %v", vars, varID), http.StatusBadRequest)
		return
	}

	hero, err := app.GetByID(appengine.NewContext(r), int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(hero)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", string(b))
}

func deleteHero(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	varID := vars["id"]
	id, err := strconv.Atoi(varID)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid id: %v", varID), http.StatusBadRequest)
		return
	}

	hero, err := app.Delete(appengine.NewContext(r), int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(hero)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", string(b))

}

func updateHero(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	hero := service.Hero{}
	err := json.NewDecoder(r.Body).Decode(&hero)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//is the position changed?
	posString := ""
	res, ok := r.URL.Query()["pos"]
	if ok || len(res) == 1 {
		posString = res[0]
	}

	var h *service.Hero
	if posString == "" {
		//no new position - update name of hero

		h, err = app.Update(appengine.NewContext(r), hero)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// new position - update list of heroes

		pos, err := strconv.Atoi(posString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		h, err = app.UpdatePosition(appengine.NewContext(r), hero, int64(pos))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	b, err := json.Marshal(h)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", string(b))
}

func searchHeroes(w http.ResponseWriter, r *http.Request) {
	name := ""
	names, ok := r.URL.Query()["name"]
	if ok || len(names) == 1 {
		name = names[0]
	}

	heroes, err := app.List(appengine.NewContext(r), name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(heroes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", string(b))
}
