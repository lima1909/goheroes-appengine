package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lima1909/goheroes-appengine/db"
	"github.com/lima1909/goheroes-appengine/service"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

var (
	svc service.HeroService
	// Heroes list from Hero examples
	Heroes = []service.Hero{
		service.Hero{ID: 1, Name: "Jasmin"},
		service.Hero{ID: 2, Name: "Mario"},
		service.Hero{ID: 3, Name: "Alex M"},
		service.Hero{ID: 4, Name: "Adam O"},
		service.Hero{ID: 5, Name: "Shauna C"},
		service.Hero{ID: 6, Name: "Lena H"},
		service.Hero{ID: 7, Name: "Chris S"},
	}
)

func main() {

	svc = db.DataStoreService{}

	router := mux.NewRouter()

	router.HandleFunc("/api/heroes", heroList)
	// router.HandleFunc("/api/heroes/add", addHero)

	// router.HandleFunc("/api/heroes/{id:[0-9]+}", getHeroByID)

	http.Handle("/", router)

	// log.Println("Start Server: http://localhost:8081")
	// log.Fatalln(http.ListenAndServe(":8081", nil))
	appengine.Main()
}

func heroList(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, fmt.Errorf("Only GET is allowed").Error(), http.StatusMethodNotAllowed)
		return
	}

	name := ""
	names, ok := r.URL.Query()["name"]
	if ok || len(names) == 1 {
		name = names[0]
	}
	log.Infof(appengine.NewContext(r), "Names=%v -> Name=%v", names, name)

	// v := mux.Vars(r)
	// name := v["name"]

	w.Header().Set("Access-Control-Allow-Origin", "*")
	heroes, err := svc.List(appengine.NewContext(r), name)
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

// func addHero(w http.ResponseWriter, r *http.Request) {
// 	h := db.Hero{ID: "99", Name: "Foo"}
// 	h, err := db.AddHero(appengine.NewContext(r), h)
// 	if err != nil {
// 		fmt.Fprintf(w, "%v", err)
// 	}
// 	fmt.Fprintln(w, toJSON(h))
// }

// func heros(w http.ResponseWriter, r *http.Request) {

// 	switch r.Method {
// 	case "GET":
// 		// loadHeroes(w)
// 		h, err := db.ListHeroes(appengine.NewContext(r))
// 		if err != nil {
// 			fmt.Fprintf(w, "%v", err)
// 		}
// 		fmt.Fprintln(w, toJSON(h))
// 	case "OPTIONS":
// 		writeToClient(w, string(http.StatusOK))
// 	case "PUT":
// 		updateHero(w, r)
// 	case "POST":
// 		addHero(w, r)
// 	}
// }

// func toJSON(v interface{}) string {
// 	b, err := json.Marshal(v)

// 	if err != nil {
// 		return fmt.Sprintf("Err by marshal heros: %v", err)
// 	}

// 	return string(b)
// }

// func loadHeroes(w http.ResponseWriter) {
// 	b, err := json.Marshal(Heroes)

// 	if err != nil {
// 		fmt.Fprintf(w, "Err by marshal heros: %v", err)
// 		return
// 	}

// 	writeToClient(w, string(b))
// }

// func getHeroByID(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	varID := vars["id"]
// 	i, err := strconv.Atoi(varID)

// 	if err != nil {
// 		fmt.Fprintf(w, "Err during string convert: %v", err)
// 		return
// 	}

// 	b, err := json.Marshal(Heroes[i-1])

// 	if err != nil {
// 		fmt.Fprintf(w, "Err by marshal hero: %v", err)
// 		return
// 	}

// 	writeToClient(w, string(b))
// }

// func updateHero(w http.ResponseWriter, r *http.Request) {
// 	hero, err := getHeroFromRequest(r, w)

// 	if err != nil {
// 		fmt.Fprintf(w, "Err by getHeroFromRequest %v", err)
// 		return
// 	}

// 	i, err := strconv.Atoi(hero.ID)

// 	if err != nil {
// 		fmt.Fprintf(w, "Err during string convert: %v", err)
// 		return
// 	}

// 	//update Hero in List
// 	Heroes[i-1] = hero

// 	writeToClient(w, "")
// }

// func addHero2(w http.ResponseWriter, r *http.Request) {
// 	hero, err := getHeroFromRequest(r, w)

// 	if err != nil {
// 		fmt.Fprintf(w, "Err by getHeroFromRequest %v", err)
// 		return
// 	}

// 	hero.ID = strconv.Itoa(len(Heroes) + 1)

// 	Heroes = append(Heroes, hero)

// 	b, err := json.Marshal(hero)

// 	if err != nil {
// 		fmt.Fprintf(w, "Err by marshal hero: %v", err)
// 		return
// 	}

// 	writeToClient(w, string(b))
// }

// func getHeroFromRequest(r *http.Request, w http.ResponseWriter) (db.Hero, error) {
// 	//read transfered data
// 	body, err := ioutil.ReadAll(r.Body)
// 	defer r.Body.Close()

// 	if err != nil {
// 		return db.Hero{}, err
// 	}

// 	//convert to Hero
// 	var hero db.Hero
// 	err = json.Unmarshal(body, &hero)

// 	if err != nil {
// 		return db.Hero{}, err
// 	}

// 	return hero, nil
// }

// func writeToClient(w http.ResponseWriter, s string) {
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

// 	fmt.Fprintf(w, s)
// }
