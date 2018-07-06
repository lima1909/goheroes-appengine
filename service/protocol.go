package service

import (
	"fmt"
	"strconv"
	"time"
)

const (
	dateFormat = "2006.01.02 15:04:05"
)

// Protocol are the logging information by each HeroService call
type Protocol struct {
	Action string    `json:"action"`
	HeroID int64     `json:"heroid"`
	Note   string    `json:"note"`
	Time   time.Time `json:"time"`
}

// GetTimeString convert Time in the right format (const: dateFormat)
func (p Protocol) GetTimeString() string {
	return p.Time.Format(dateFormat)
}

// NewProtocol create a new Protocol instance with Time = Now
func NewProtocol(action string, hID int64, note string) Protocol {
	return Protocol{
		Time:   time.Now(),
		Action: action,
		HeroID: hID,
		Note:   note,
	}
}

// NewProtocolf Sprintf for the note
func NewProtocolf(action string, hID int64, note string, a ...interface{}) Protocol {
	return NewProtocol(action, hID, fmt.Sprintf(note, a...))
}

// Protocol2Map convert Protocol to map
func Protocol2Map(p Protocol) map[string]string {
	return map[string]string{
		"Action": p.Action,
		"HeroID": strconv.Itoa(int(p.HeroID)),
		"Note":   p.Note,
		"Time":   p.GetTimeString(),
	}
}

// Map2Protocol convert map to Protocol
func Map2Protocol(m map[string]string) Protocol {
	t, err := time.Parse(dateFormat, m["Time"])
	if err != nil {
		fmt.Printf("err by parse date: %v", err)
	}
	id, _ := strconv.Atoi(m["HeroID"])
	return Protocol{
		Action: m["Action"],
		HeroID: int64(id),
		Note:   m["Note"],
		Time:   t,
	}
}
