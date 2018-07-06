package service

import (
	"os"
	"reflect"
	"strconv"
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
	HeroesService string
	RunInCloud    bool
}

// NewApp create a new App instance
func NewApp(svc ProtocolHeroService) *App {

	inCloud, _ := strconv.ParseBool(os.Getenv("RUN_IN_CLOUD"))

	return &App{
		ProtocolHeroService: svc,
		Info: Info{
			RunInCloud:    inCloud,
			HeroesService: reflect.TypeOf(svc).String(),
		},
		Version: "dev-snapshot_" + time.Now().Local().Format("2006.01.02 15:04:05"),
	}
}
