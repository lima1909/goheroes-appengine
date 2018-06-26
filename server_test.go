package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/lima1909/goheroes-appengine/service"
)

var (
	server = httptest.NewServer(handler())
)

func init() {
	os.Setenv("HERO_SERVICE_IMPL", "NotSet")
}

func TestHeroList(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost:8080/api/heroes", nil)
	w := httptest.NewRecorder()
	heroList(w, r)

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

	// check the lengths from Service with http-call
	hs, err := app.List(context.TODO(), "")
	if err != nil {
		t.Errorf("no err expected, got: %v", err)
	}
	if len(hs) != len(heroes) {
		t.Errorf("heroes expected: %v and get: %v", len(hs), len(heroes))
	}
}

func TestHeroList_HandlerCORS(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/api/heroes", server.URL))
	if err != nil {
		t.Errorf("No err expected: %v", err)
	}

	// check Header: Access-Control-Allow-Origin
	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf(`expect "*" but get: %v`, resp.Header.Get("Access-Control-Allow-Origin"))
	}
}

func TestGetHeroID(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost:8080/api/heroes", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	getHero(w, r)

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
}

func TestGetHeroID_HandlerCORS(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/api/heroes/2", server.URL))
	if err != nil {
		t.Errorf("No err expected: %v", err)
	}

	// check Header: Access-Control-Allow-Origin
	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf(`expect "*" but get: %v`, resp.Header.Get("Access-Control-Allow-Origin"))
	}
}

func TestSearchHeroes(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost:8080/api/heroes", nil)
	q := r.URL.Query()
	q.Add("name", "Jasmin")
	r.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()
	heroList(w, r)

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
}

func TestSearchHeroesWithEmptyName(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost:8080/api/heroes", nil)
	w := httptest.NewRecorder()
	heroList(w, r)

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
	if len(heroes) != 7 {
		t.Errorf("expect all heroes as search-result, but %v", len(heroes))
	}
}

func TestSearchHeroes_HandlerCORS(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/api/heroes/?name=%s", server.URL, url.QueryEscape("Adam O")))
	if err != nil {
		t.Errorf("No err expected: %v", err)
	}

	// check Header: Access-Control-Allow-Origin
	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf(`expect "*" but get: %v`, resp.Header.Get("Access-Control-Allow-Origin"))
	}
}

func TestOptionsCORS(t *testing.T) {
	req, err := http.NewRequest("OPTIONS", fmt.Sprintf("%s/api/heroes", server.URL), nil)
	if err != nil {
		t.Errorf("No err expected: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("No err expected: %v", err)
	}

	// check StatusOK
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expect status ok (200), but is: %v", resp.StatusCode)
	}

	// check Header: Access-Control-Allow-Origin
	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf(`expect "*" but get: %v`, resp.Header.Get("Access-Control-Allow-Origin"))
	}

	// by OPTIONS you get no body
	body, _ := ioutil.ReadAll(resp.Body)
	if len(body) != 0 {
		t.Errorf("no body expect, but is: %v", string(body))
	}
}

func TestAddHeroHandlerCORS(t *testing.T) {
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/api/heroes", server.URL),
		strings.NewReader(` { "name" : "Test" } `))
	if err != nil {
		t.Errorf("No err expected: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("No err expected: %v", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	strBody := string(body)
	if strings.Contains(strBody, "Test") == false {
		t.Errorf("expect: Test in body, got: %v", strBody)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status ok (200), but is: %v (%v)", resp.StatusCode, strBody)
	}

	// check Header: Access-Control-Allow-Origin
	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf(`expect "*" but get: %v`, resp.Header.Get("Access-Control-Allow-Origin"))
	}
}

func TestDeleteHero(t *testing.T) {
	heroes, _ := app.List(context.TODO(), "")
	hLen := len(heroes)

	r := httptest.NewRequest("GET", "http://localhost:8080/api/heroes", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "2"})
	w := httptest.NewRecorder()
	deleteHero(w, r)

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

	if hero.ID != 2 {
		t.Errorf("expect ID=2, but is: %v", hero.ID)
	}

	heroes, err = app.List(context.TODO(), "")
	if err != nil {
		t.Errorf("no err expected, got: %v", err)
	}

	if hLen-1 != len(heroes) {
		t.Errorf("expect heroes size: %v, got: %v", (hLen - 1), len(heroes))
	}
}

func TestDeleteHero_HandlerCORS(t *testing.T) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/heroes/3", server.URL), nil)
	if err != nil {
		t.Errorf("No err expected: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("No err expected: %v", err)
	}

	// check Header: Access-Control-Allow-Origin
	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf(`expect "*" but get: %v`, resp.Header.Get("Access-Control-Allow-Origin"))
	}
}

func TestUpdateHero(t *testing.T) {
	req, err := http.NewRequest("PUT",
		fmt.Sprintf("%s/api/heroes", server.URL),
		strings.NewReader(` { "name" : "Test", "id" : 1} `))
	if err != nil {
		t.Errorf("No err expected: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("No err expected: %v", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	strBody := string(body)
	if strings.Contains(strBody, "Test") == false {
		t.Errorf("expect: Test in body, got: %v", strBody)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status ok (200), but is: %v (%v)", resp.StatusCode, strBody)
	}

	// check Header: Access-Control-Allow-Origin
	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf(`expect "*" but get: %v`, resp.Header.Get("Access-Control-Allow-Origin"))
	}
}

func TestGetHeroFromService(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost:8080/api/heroes", strings.NewReader(` { "name" : "Test" } `))
	hero, err := getHeroFromService(r)
	if err != nil {
		t.Errorf("no err expected, got: %v", err)
	}
	if hero.Name != "Test" {
		t.Errorf("expect hero name: Test, got: %v", hero.Name)
	}
	if hero.ID != 0 {
		t.Errorf("expect hero ID == 0, got: %v", hero.ID)
	}
}

func TestGetHeroFromServiceFail(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost:8080/api/heroes", strings.NewReader(` { "name" : "Test"  `))
	_, err := getHeroFromService(r)
	if err == nil {
		t.Errorf("expected err, got nil")
	}
}
