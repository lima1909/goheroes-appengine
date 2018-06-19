package service

// App is the Entrypoint
type App struct {
	HeroService
}

// NewApp create a new App instance
func NewApp(svc HeroService) *App {
	return &App{svc}
}
