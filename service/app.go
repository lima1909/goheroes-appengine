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

// RunInCloud check Env: RUN_IN_CLOUD is set tue true
func RunInCloud() bool {
	inCloud, _ := strconv.ParseBool(os.Getenv("RUN_IN_CLOUD"))
	return inCloud
}

// NewApp create a new App instance
func NewApp(svc ProtocolHeroService) *App {
	return &App{
		ProtocolHeroService: svc,
		Info: Info{
			RunInCloud:    RunInCloud(),
			HeroesService: reflect.TypeOf(svc).String(),
		},
		Version: "dev-snapshot_" + time.Now().Local().Format("2006.01.02 15:04:05"),
	}
}
