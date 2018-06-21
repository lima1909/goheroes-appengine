package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/lima1909/goheroes-appengine/service"
)

func init() {
	os.Setenv("HERO_SERVICE_IMPL", "NotSet")
}

func TestGetHeroes_heroList(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost:8080/api/heroes", nil)
	w := httptest.NewRecorder()
	heroes(w, r)

	// check status code
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status ok (200), but is: %v", resp.StatusCode)
	}

	// check size of heroes
	heroes := make([]service.Hero, 0)
	body, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(body, &heroes)
	if err != nil {
		t.Errorf("No err expected: %v", err)
	}

	hs, _ := app.List(context.TODO(), "")
	if len(hs) != len(heroes) {
		t.Errorf("heroes expected: %v and get: %v", len(hs), len(heroes))
	}

	// check Header: contenttype
	if resp.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf(`expect "text/plain; charset=utf-8" but get: %v`, resp.Header.Get("Content-Type"))
	}
	// check Header: Access-Control-Allow-Origin
	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf(`expect "*" but get: %v`, resp.Header.Get("Access-Control-Allow-Origin"))
	}
}

func TestGetHeroID_getHero(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost:8080/api/heroes", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	heroID(w, r)

	// check status code
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status ok (200), but is: %v (%v)", resp.StatusCode, string(body))
	}

	// check result - hero
	hero := service.Hero{}
	err := json.Unmarshal(body, &hero)
	if err != nil {
		t.Errorf("No err expected: %v", err)
	}

	hr, _ := app.GetByID(context.TODO(), int64(1))
	if hero.ID != 1 {
		t.Errorf("expect ID=1, but is: %v", hero.ID)
	}
	if hero.Name != hr.Name {
		t.Errorf("expect Name: %v, but is: %v", hr.Name, hero.Name)
	}

	// check Header: Access-Control-Allow-Origin
	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf(`expect "*" but get: %v`, resp.Header.Get("Access-Control-Allow-Origin"))
	}
}

func TestGetHeroID_searchHeroes(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost:8080/api/heroes", nil)
	q := r.URL.Query()
	q.Add("name", "Jasmin")
	r.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()
	searchHeroes(w, r)

	// check status code
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status ok (200), but is: %v (%v)", resp.StatusCode, string(body))
	}

	heroes := make([]service.Hero, 0)
	err := json.Unmarshal(body, &heroes)
	if err != nil {
		t.Errorf("No err expected: %v", err)
	}
	// check result: one Hero
	if len(heroes) != 1 {
		t.Errorf("expect one hero as search-result, but %v", len(heroes))
	}

	// check Header: Access-Control-Allow-Origin
	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf(`expect "*" but get: %v`, resp.Header.Get("Access-Control-Allow-Origin"))
	}
}
