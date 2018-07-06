package service

import (
	"os"
	"reflect"
	"time"
)

// App is the Entrypoint
type App struct {
	ProtocolHeroService

	Info    Info
	Version string
}

// Info to the current system
type Info struct {
	HeroesService      string
	EnvHeroServiceImpl string
}

// NewApp create a new App instance
func NewApp(svc ProtocolHeroService) *App {
	return &App{
		ProtocolHeroService: svc,
		Info: Info{
			EnvHeroServiceImpl: os.Getenv("HERO_SERVICE_IMPL"),
			HeroesService:      reflect.TypeOf(svc).String(),
		},
		Version: "dev-snapshot_" + time.Now().Local().Format("2006.01.02 15:04:05"),
	}
}
