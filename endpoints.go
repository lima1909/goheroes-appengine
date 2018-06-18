package main

import (
	"log"

	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	"github.com/lima1909/goheroes-appengine/service"
)

func init() {
	heroService := &service.HeroService{}
	_, err := endpoints.RegisterService(heroService, "heroes", "v1", "Heroes API", true)
	if err != nil {
		log.Fatalf("Register service: %v", err)
	}

	endpoints.HandleHTTP()
}
