package db

import (
	"context"
	"log"
	"strings"
	"time"

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

	return &MemService{heroes: heroes, maxID: maxID}
}

// Protocols impl from ProtocolService
func (MemService) Protocols(c context.Context) ([]service.Protocol, error) {
	t := time.Now()

	dummyProtocols := make([]service.Protocol, 8)

	dummyProtocols[0] = service.NewProtocolf("Add", 1, "add new Hero with ID: 1")
	dummyProtocols[1] = service.NewProtocolf("List", 0, "List from Heroes with len: 5")
	dummyProtocols[2] = service.Protocol{Action: "Delete", HeroID: 2, Note: "delete Hero with ID: 2", Time: t.Add(time.Duration(-10) * time.Minute)}
	dummyProtocols[3] = service.Protocol{Action: "Search", HeroID: 0, Note: "search list", Time: t.Add(time.Duration(-7) * time.Hour)}
	dummyProtocols[4] = service.Protocol{Action: "Add", HeroID: 5, Note: "add Hero with ID: 5", Time: t.Add(time.Duration(-170) * time.Second)}
	dummyProtocols[5] = service.Protocol{Action: "List", HeroID: 0, Note: "List from Heroes", Time: t.AddDate(0, 0, -1)}
	dummyProtocols[6] = service.Protocol{Action: "Delete", HeroID: 23, Note: "delete Hero with ID: 23", Time: t.AddDate(0, -2, 0)}
	dummyProtocols[7] = service.Protocol{Action: "Search", HeroID: 0, Note: "search list", Time: t.AddDate(-3, 0, 0)}

	return dummyProtocols, nil
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
