package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lima1909/goheroes-appengine/db"
	"github.com/lima1909/goheroes-appengine/service"

	"google.golang.org/appengine"
)

var app *service.App

// handle CORS and the OPION method
func corsAndOptionHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
		} else {
			h.ServeHTTP(w, r)
		}
	}
}

// create all used Handler
func handler() http.Handler {
	router := mux.NewRouter()

	router.Handle("/", http.RedirectHandler("/info", http.StatusFound))
	router.HandleFunc("/info", infoPage)

	url := "/api/heroes"
	router.HandleFunc(url, heroList).Methods("GET")
	router.HandleFunc(url, addHero).Methods("POST")
	router.HandleFunc(url, switchHero).Methods("PUT").Queries("pos", "{pos}")
	router.HandleFunc(url, updateHero).Methods("PUT")
	router.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "invalid method: "+r.Method, http.StatusBadRequest)
	}).Methods("DELETE", "PATH", "COPY", "HEAD", "LINK", "UNLINK", "PURGE", "LOCK", "UNLOCK", "VIEW", "PROPFIND")

	urlWithID := "/api/heroes/{id:[0-9]+}"
	router.HandleFunc(urlWithID, getHero).Methods("GET")
	router.HandleFunc(urlWithID, deleteHero).Methods("DELETE")
	router.HandleFunc(urlWithID, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "invalid method: "+r.Method, http.StatusBadRequest)
	}).Methods("PUT", "POST", "PATH", "COPY", "HEAD", "LINK", "UNLINK", "PURGE", "LOCK", "UNLOCK", "VIEW", "PROPFIND")

	urlWithScores := "/api/heroes/scores"
	router.HandleFunc(urlWithScores, getScores).Methods("GET")

	// TODO: not necessary anymore (only for the slash on the end)
	router.HandleFunc("/api/heroes/", heroList)

	return corsAndOptionHandler(router)
}

func init() {
	http.Handle("/", handler())
	app = service.NewApp(db.NewMemService())

	log.Println("Init is ready and start the server on: http://localhost:8080")
}

/** This is needed for tests */
func resetMemService() {
	app = service.NewApp(db.NewMemService())
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

	err = t.Execute(w, app.Info)
	if err != nil {
		fmt.Fprintf(w, "Err: %v\n", err)
		return
	}
}

func getScores(w http.ResponseWriter, r *http.Request) {
	scoreMap := app.CreateScoreMap(appengine.NewContext(r))

	b, err := json.Marshal(scoreMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", b)
}

func heroList(w http.ResponseWriter, r *http.Request) {
	heroes, err := app.List(appengine.NewContext(r), r.URL.Query().Get("name"))
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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	heroName := string(body)
	h, err := app.Add(appengine.NewContext(r), heroName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeHeroToClient(w, r, h)
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

	writeHeroToClient(w, r, hero)
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

	writeHeroToClient(w, r, hero)

}

func updateHero(w http.ResponseWriter, r *http.Request) {
	hero, err := getHeroFromService(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h, err := app.Update(appengine.NewContext(r), hero)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeHeroToClient(w, r, h)
}

func switchHero(w http.ResponseWriter, r *http.Request) {

	hero, err := getHeroFromService(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pos := r.FormValue("pos")
	posNb, err := strconv.Atoi(pos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h, err := app.UpdatePosition(appengine.NewContext(r), hero, int64(posNb))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeHeroToClient(w, r, h)

}

func getHeroFromService(r *http.Request) (service.Hero, error) {
	defer r.Body.Close()

	hero := service.Hero{}
	return hero, json.NewDecoder(r.Body).Decode(&hero)
}

func writeHeroToClient(w http.ResponseWriter, r *http.Request, h *service.Hero) {
	b, err := json.Marshal(h)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", string(b))
}
