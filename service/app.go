package service

import (
	"os"
	"reflect"
)

// App is the Entrypoint
type App struct {
	HeroService
	Info Info
}

// Info to the current system
type Info struct {
	HeroesService      string
	EnvHeroServiceImpl string
}

// NewApp create a new App instance
func NewApp(svc HeroService) *App {
	return &App{
		HeroService: svc,
		Info: Info{
			EnvHeroServiceImpl: os.Getenv("HERO_SERVICE_IMPL"),
			HeroesService:      reflect.TypeOf(svc).String(),
		},
	}
}
