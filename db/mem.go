package db

import (
	"context"
	"log"
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
		service.Hero{ID: 1, Name: "Jasmin"},
		service.Hero{ID: 2, Name: "Mario"},
		service.Hero{ID: 3, Name: "Alex M"},
		service.Hero{ID: 4, Name: "Adam O"},
		service.Hero{ID: 5, Name: "Shauna C"},
		service.Hero{ID: 6, Name: "Lena H"},
		service.Hero{ID: 7, Name: "Chris S"},
	}
	maxID := int64(7)

	return &MemService{heroes, maxID}
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
func (m *MemService) Add(c context.Context, h service.Hero) (*service.Hero, error) {
	m.maxID++
	h.ID = m.maxID
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
	for i, hero := range m.heroes {
		if hero.ID == h.ID {
			oldPos = i
			break
		}
	}

	newHeroesSlice := append(m.heroes[:oldPos], m.heroes[oldPos+1:]...)
	m.heroes = append(newHeroesSlice[:pos], append([]service.Hero{h}, newHeroesSlice[pos:]...)...)

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
