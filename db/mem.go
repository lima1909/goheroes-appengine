package db

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/lima1909/goheroes-appengine/service"
)

// MemService is a Impl from service.HeroService
type MemService struct {
	heroes []service.Hero
	maxID  int64
}

// NewMemService create a new instance of MemService
func NewMemService() *MemService {
	heroes := []service.Hero{
		service.Hero{ID: 1, Name: "Jasmin", ScoreData: service.ScoreData{Name: "jasmin-roeper", City: "Nuremberg", Country: "de"}},
		service.Hero{ID: 2, Name: "Mario", ScoreData: service.ScoreData{Name: "mario-linke", City: "Nürnberg", Country: "de"}},
		service.Hero{ID: 3, Name: "Alex M"},
		service.Hero{ID: 4, Name: "Adam O"},
		service.Hero{ID: 5, Name: "Shauna C"},
		service.Hero{ID: 6, Name: "Lena H"},
		service.Hero{ID: 7, Name: "Chris S"},
	}
	maxID := int64(7)

	return &MemService{heroes, maxID}
}

// Protocols impl from ProtocolService
func (MemService) Protocols(c context.Context) ([]service.Protocol, error) {
	return []service.Protocol{
		service.NewProtocolf("Add", 1, "add new Hero with ID: 1"),
		service.NewProtocolf("List", 0, "List from Heroes with len: 5"),
	}, nil
}

// List all Heroes, there are saved in the heroes array
func (m MemService) List(c context.Context, name string) ([]service.Hero, error) {
	if name == "" {
		return m.heroes, nil
	}

	hs := make([]service.Hero, 0)
	for _, h := range m.heroes {
		if strings.Contains(strings.ToUpper(h.Name), strings.ToUpper(name)) { //need uppercase to make it case insensitiv
			hs = append(hs, h)
			log.Printf("find hero: %v\n", h)
		}
	}
	return hs, nil
}

// GetByID get Hero by the ID
func (m MemService) GetByID(c context.Context, id int64) (*service.Hero, error) {
	for _, h := range m.heroes {
		if h.ID == id {
			return &h, nil
		}
	}
	return nil, service.ErrHeroNotFound
}

// Add an Hero
func (m *MemService) Add(c context.Context, name string) (*service.Hero, error) {
	m.maxID++

	h := service.Hero{Name: name, ID: m.maxID}
	m.heroes = append(m.heroes, h)
	log.Printf("add hero: %v\n", h)
	return &h, nil
}

// Update an Hero
func (m *MemService) Update(c context.Context, h service.Hero) (*service.Hero, error) {

	for i, hero := range m.heroes {
		if hero.ID == h.ID {
			m.heroes[i] = h
			log.Printf("update hero from: %v to: %v\n", hero, h)
			return &m.heroes[i], nil
		}
	}

	return nil, service.ErrHeroNotFound
}

// UpdatePosition of Hero
func (m *MemService) UpdatePosition(c context.Context, h service.Hero, pos int64) (*service.Hero, error) {

	if pos > int64(len(m.heroes)+1) {
		return nil, service.ErrPosNotFound
	}

	oldPos := 0
	//need to get the hero on the server because of additional datas like scoreData
	heroOnServer := service.Hero{}
	for i, hero := range m.heroes {
		if hero.ID == h.ID {
			oldPos = i
			heroOnServer = hero
			break
		}
	}

	newHeroesSlice := append(m.heroes[:oldPos], m.heroes[oldPos+1:]...)
	m.heroes = append(newHeroesSlice[:pos], append([]service.Hero{heroOnServer}, newHeroesSlice[pos:]...)...)

	//just for debugging and logging
	for i, hero := range m.heroes {
		if hero.ID == h.ID {
			log.Printf("update pos of %v from: %v to: %v\n", hero.Name, oldPos, i)
			break
		}
	}

	return &m.heroes[pos], nil
}

// Delete an Hero
func (m *MemService) Delete(c context.Context, id int64) (*service.Hero, error) {
	hero := service.Hero{ID: -1}

	for i, h := range m.heroes {
		if h.ID == id {
			hero = h
			//remove from List
			log.Printf("delete hero: %v\n", hero)
			m.heroes = append(m.heroes[:i], m.heroes[i+1:]...)

			return &hero, nil
		}
	}

	return nil, service.ErrHeroNotFound
}

// CreateScoreMap to get the scores from 8a.nu
func (m *MemService) CreateScoreMap(c context.Context) map[int64]int {
	scoreMap := map[int64]int{}

	for _, h := range m.heroes {
		scoreMap[h.ID] = 0
		if h.ScoreData.Name != "" {
			url := fmt.Sprintf("https://www.8a.nu/%s/scorecard/ranking/?City=%s", h.ScoreData.Country, h.ScoreData.City)
			score, err := getScore(url, h.ScoreData.Name)
			if err == nil {
				scoreMap[h.ID] = score
			}
		}
	}

	return scoreMap
}

func getScore(urlString string, name string) (int, error) {
	// Make HTTP GET request
	response, err := http.Get(urlString)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	// Get the response body as a string
	dataInBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	pageContent := string(dataInBytes)
	if pageContent == "" {
		return 0, service.ErrNoContent
	}

	// Find a substr
	startIndex := strings.Index(pageContent, name)
	if startIndex == -1 {
		return 0, fmt.Errorf("Can not find %v on %v ", name, urlString)
	}

	subString := pageContent[startIndex:(startIndex + 200)]

	// Find score
	indexStart := strings.Index(subString, "\">")
	indexEnd := strings.Index(subString, "</a>")

	if indexStart == -1 || indexEnd == -1 {
		return 0, fmt.Errorf("Can not find score for %v " + name)
	}

	return convertToNumber(subString[(indexStart + 2):indexEnd]), nil
}

func convertToNumber(s string) int {
	re := regexp.MustCompile("[0-9]+")
	scoreNumberArray := re.FindAllString(s, -1)

	scoreNumberString := ""
	for _, c := range scoreNumberArray {
		scoreNumberString = scoreNumberString + c
	}

	nb, err := strconv.Atoi(scoreNumberString)
	if err != nil {
		return 0
	}

	return nb
}
