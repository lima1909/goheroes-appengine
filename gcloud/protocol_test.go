package gcloud

import (
	"testing"
)

var (
	protocol = NewProtocol("Add", 2, "A note ...")
)

func TestProtocol2Map(t *testing.T) {
	m := protocol2Map(protocol)
	if m["Action"] != "Add" {
		t.Errorf("Add != %v", m["Action"])
	}
	if m["HeroID"] != "2" {
		t.Errorf("2 != %v", m["HeroID"])
	}
	if m["Note"] != "A note ..." {
		t.Errorf("A note ... != %v", m["Note"])
	}
}

func TestMap2Protocol(t *testing.T) {
	m := map[string]string{
		"Action": "Delete",
		"HeroID": "7",
		"Time":   protocol.GetTimeString(),
	}
	p := map2Protocol(m)
	if p.Action != "Delete" {
		t.Errorf("Delete != %v", p.Action)
	}
	if p.HeroID != 7 {
		t.Errorf("7 != %v", p.HeroID)
	}
	if p.GetTimeString() != protocol.GetTimeString() {
		t.Errorf("%v != %v", p.GetTimeString(), protocol.GetTimeString())
	}
}
