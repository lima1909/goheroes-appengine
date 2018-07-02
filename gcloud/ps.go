// Package gcloud package contains specific gcloud API calls, like Pub/Sub
//
// https://developers.google.com/apis-explorer/#p/pubsub/v1/pubsub.projects.subscriptions.pull
//
// to read the message-queue:
// gcloud.cmd pubsub subscriptions pull HERO_SUB
// with: --auto-ack you can clear the queue
//
// PubSub in the App Engine runs only with OLD impl!!!
// good example find here: https://github.com/d2g/dg-pubsubtest
//
// find current project dynamic: appengine.RequestID(c) (ProjectID = "goheros-207118")
//
package gcloud

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/oauth2/google"

	"github.com/lima1909/goheroes-appengine/service"
	pubsub "google.golang.org/api/pubsub/v1"
	"google.golang.org/appengine/log"
)

const (
	// ProjectID from Cloud Project
	ProjectID = "goheros-207118"
)

// HeroService wrapeer to HeroService
// protocoll HeroService calls
type HeroService struct {
	hs service.HeroService
}

// NewHeroService create a new instance
func NewHeroService(hs service.HeroService) *HeroService {
	return &HeroService{hs: hs}
}

// List protocoll list call
func (hs HeroService) List(c context.Context, name string) ([]service.Hero, error) {
	l, err := hs.hs.List(c, name)
	pub(c, Protocol{
		Action: "LIST",
		Time:   time.Now(),
		Note:   fmt.Sprintf("get list with name: %s and size: %v", name, len(l)),
	})
	return l, err
}

// GetByID delegate to HeroService
func (hs HeroService) GetByID(c context.Context, id int64) (*service.Hero, error) {
	hero, err := hs.hs.GetByID(c, id)
	pub(c, Protocol{
		Action: "GetByID",
		HeroID: id,
		Time:   time.Now(),
		Note:   fmt.Sprintf("GetByID find Hero: %v", hero),
	})
	return hero, err
}

// Add delegate to HeroService
func (hs HeroService) Add(c context.Context, n string) (*service.Hero, error) {
	h, err := hs.hs.Add(c, n)
	pub(c, Protocol{
		Action: "Add",
		HeroID: h.ID,
		Time:   time.Now(),
		Note:   fmt.Sprintf("Add Hero: %v with Name: %s", h, n),
	})
	return h, err
}

// Update delegate to HeroService
func (hs HeroService) Update(c context.Context, h service.Hero) (*service.Hero, error) {
	pub(c, Protocol{
		Action: "Update",
		HeroID: h.ID,
		Time:   time.Now(),
		Note:   fmt.Sprintf("Update Hero: %v", h),
	})
	return hs.hs.Update(c, h)
}

// UpdatePosition delegate to HeroService
func (hs HeroService) UpdatePosition(c context.Context, h service.Hero, pos int64) (*service.Hero, error) {
	hero, err := hs.hs.UpdatePosition(c, h, pos)
	pub(c, Protocol{
		Action: "UpdatePosition",
		HeroID: h.ID,
		Time:   time.Now(),
		Note:   fmt.Sprintf("UpdatePosition Hero: %v with new Pos: %v", hero, pos),
	})
	return hero, err
}

// Delete delegate to HeroService
func (hs HeroService) Delete(c context.Context, id int64) (*service.Hero, error) {
	h, err := hs.hs.Delete(c, id)
	pub(c, Protocol{
		Action: "Delete",
		HeroID: id,
		Time:   time.Now(),
		Note:   fmt.Sprintf("Delete Hero: %v with ID: %v", h, id),
	})
	return h, err
}

// Protocol are the logging information by each HeroService call
type Protocol struct {
	Action string    `json:"action"`
	HeroID int64     `json:"heroid"`
	Note   string    `json:"note"`
	Time   time.Time `json:"time"`
}

func protocol2Map(p Protocol) map[string]string {
	return map[string]string{
		"Action": p.Action,
		"HeroID": strconv.Itoa(int(p.HeroID)),
		"Note":   p.Note,
		"Time":   p.Time.String(),
	}
}

func map2Protocol(m map[string]string) Protocol {
	t, _ := time.Parse("2006.01.02 15:04:05", m["Time"])
	id, _ := strconv.Atoi(m["HeroID"])
	return Protocol{
		Action: m["Action"],
		HeroID: int64(id),
		Note:   m["Note"],
		Time:   t,
	}
}

func createSevice(c context.Context) (*pubsub.Service, error) {
	hc, err := google.DefaultClient(c, pubsub.PubsubScope)
	if err != nil {
		return nil, fmt.Errorf("can not create new default client: %v", err)
	}

	svc, err := pubsub.New(hc)
	if err != nil {
		return nil, fmt.Errorf("can not create new service: %v", err)
	}

	return svc, nil
}

func pub(c context.Context, p Protocol) {
	svc, err := createSevice(c)
	if err != nil {
		log.Errorf(c, "Publish create service error: %v", err)
		return
	}

	_, err = svc.Projects.Topics.Publish(
		"projects/goheros-207118/topics/HERO",
		&pubsub.PublishRequest{
			Messages: []*pubsub.PubsubMessage{
				{
					Attributes: protocol2Map(p),
					Data:       base64.StdEncoding.EncodeToString([]byte("pub protcol message")),
				},
			},
		},
	).Do()
	if err != nil {
		log.Errorf(c, "Publish error: %v", err)
	}
}

// Sub ...
func Sub(c context.Context) ([]Protocol, error) {
	svc, err := createSevice(c)
	if err != nil {
		return nil, err
	}

	result, err := svc.Projects.Subscriptions.Pull(
		"projects/goheros-207118/subscriptions/HERO_SUB",
		&pubsub.PullRequest{
			MaxMessages:       5,
			ReturnImmediately: true,
		},
	).Do()
	if err != nil {
		e := fmt.Errorf("Publish error:  %v", err)
		log.Errorf(c, "%v", e)
		return nil, e
	}

	ps := make([]Protocol, len(result.ReceivedMessages))
	for i, m := range result.ReceivedMessages {
		// b, _ := base64.StdEncoding.DecodeString(m.Message.Data)
		ps[i] = map2Protocol(m.Message.Attributes)
		ack(c, m.AckId)
		err = Add(c, ps[i])
		if err != nil {
			log.Errorf(c, "err by add protocol to datastore: %v", err)
		}
	}

	return ps, nil
}

// acknowledge all messages with the ack ID
func ack(c context.Context, ackIDs ...string) {
	svc, err := createSevice(c)
	if err != nil {
		log.Errorf(c, "Acknowledge create Service error: %v", err)
		return
	}

	_, err = svc.Projects.Subscriptions.Acknowledge(
		"projects/goheros-207118/subscriptions/HERO_SUB",
		&pubsub.AcknowledgeRequest{AckIds: ackIDs},
	).Do()

	if err != nil {
		log.Errorf(c, "Acknowledge error by execute acknowledge-request: %v", err)
	}
}
