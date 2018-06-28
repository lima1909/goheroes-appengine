package ps

import (
	"context"

	// "cloud.google.com/go/pubsub"
	"github.com/lima1909/goheroes-appengine/service"
	"google.golang.org/appengine/log"
)

const (
	// ProjectID from Cloud Project
	ProjectID = "goheros-207118"
	// Topic for pub sub
	Topic = "HERO"
)

// HeroService wrapeer to HeroService
// protocoll HeroService calls
type HeroService struct {
	service.HeroService
	hs service.HeroService
}

// NewHeroService create a new instance
func NewHeroService(hs service.HeroService) *HeroService {
	return &HeroService{hs: hs}
}

// List protocoll list call
func (hs HeroService) List(c context.Context, name string) ([]service.Hero, error) {

	heroes, err := hs.hs.List(c, name)
	if err != nil {
		log.Errorf(c, "%v", err)
		return nil, err
	}
	log.Infof(c, "Get Hero-List: %v", heroes)

	// client, err := pubsub.NewClient(c, ProjectID)
	// if err != nil {
	// 	e := fmt.Errorf("can no create new client for project %s: %v", ProjectID, err)
	// 	log.Errorf(c, "%v", e)
	// 	return nil, e
	// }
	// log.Infof(c, "PubSub client created: %v", client)

	// topic, err := client.Topic(Topic)
	// if err != nil {
	// 	e := fmt.Errorf("can not create topic %s: %v", Topic, err)
	// 	log.Errorf(c, "%v", e)
	// 	return nil, e
	// }
	// log.Infof(c, "Get PubSub topic: %s -- %v", Topic, client)

	// msg := &pubsub.Message{
	// 	Data: []byte("payload: " + string(len(heroes))),
	// }
	// serverID, err := topic.Publish(c, msg).Get(c)
	// if err != nil {
	// 	e := fmt.Errorf("can not publish the message: %v", err)
	// 	log.Errorf(c, "%v", e)
	// 	return nil, e
	// }

	// log.Infof(c, "Publish on Topic: %s the message: %v with ServerID: %s", Topic, msg, serverID)
	return heroes, nil
}
