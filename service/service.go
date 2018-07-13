// Package service is the interface to interact with Heroes
package service

import (
	"context"
	"errors"
	"os"
	"strconv"
)

var (
	// ErrHeroNotFound if no Hero was found
	ErrHeroNotFound = errors.New("Hero not Found")
	// ErrPosNotFound if new Position is out of range
	ErrPosNotFound = errors.New("Out of Range")
	// ErrNoContent if reading 8a.nu returns empty string
	ErrNoContent = errors.New("No content found on 8a.nu")
)

// Hero the struct
type Hero struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	ScoreData ScoreData `json:"-"`
}

// ScoreData - to create the correct search url
type ScoreData struct {
	Name    string
	City    string
	Country string
}

// HeroService access to Heroes methods
type HeroService interface {
	List(c context.Context, name string) ([]Hero, error)
	GetByID(c context.Context, id int64) (*Hero, error)
	Add(c context.Context, n string) (*Hero, error)
	Update(c context.Context, h Hero) (*Hero, error)
	UpdatePosition(c context.Context, h Hero, pos int64) (*Hero, error)
	Delete(c context.Context, id int64) (*Hero, error)
	CreateScoreMap(c context.Context) (map[int64]int, error)
}

// ProtocolService acces to the Protocols
type ProtocolService interface {
	Protocols(c context.Context) ([]Protocol, error)
}

// ProtocolHeroService combine Hero and ProtocolService
type ProtocolHeroService interface {
	HeroService
	ProtocolService
}

// RunInCloud check Env: RUN_IN_CLOUD is set tue true
func RunInCloud() bool {
	inCloud, _ := strconv.ParseBool(os.Getenv("RUN_IN_CLOUD"))
	return inCloud
}
