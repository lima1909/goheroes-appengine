package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lima1909/goheroes-appengine/db"
	"github.com/lima1909/goheroes-appengine/service"

	"google.golang.org/appengine"
)

var (
	app *service.App
)

func getHostAndPort() string {
	port := "8080"
	if s := os.Getenv("PORT"); s != "" {
		port = s
	}

	host := ""
	if appengine.IsDevAppServer() {
		host = "127.0.0.1"
	}

	return host + ":" + port
}

func main() {

	app = service.NewApp(db.NewMemService())
	// app = service.NewApp(db.DatastoreService{})

	router := mux.NewRouter()

	router.Handle("/", http.RedirectHandler("/api/heroes", http.StatusFound))

	router.HandleFunc("/example", example)
	router.HandleFunc("/api/heroes", heroes)
	router.HandleFunc("/api/heroes/{id:[0-9]+}", heroID)
	http.Handle("/", router)

	log.Println("Server is started on: ", getHostAndPort())
	appengine.Main()
}

func example(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/index.html")
	if err != nil {
		fmt.Fprintf(w, "Err: %v\n", err)
		return
	}

	err = t.Execute(w, "Hello World!")
	if err != nil {
		fmt.Fprintf(w, "Err: %v\n", err)
		return
	}
}

func heroes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		heroList(w, r)
	case "POST":
		addHero(w, r)
	case "PUT":
		updateHero(w, r)

	default:
		http.Error(w, "invalid method: "+r.Method, http.StatusBadRequest)
	}
}

func heroID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

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
		http.Error(w, fmt.Sprintf("invalid id: %v", varID), http.StatusBadRequest)
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

	err = app.Delete(appengine.NewContext(r), int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "")
}

func updateHero(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	hero := service.Hero{}
	err := json.NewDecoder(r.Body).Decode(&hero)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h, err := app.Update(appengine.NewContext(r), hero)
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
