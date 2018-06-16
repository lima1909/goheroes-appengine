package main

import "testing"

func TestFoo(t *testing.T) {
	a := 1
	b := 2
	if a != b {
		t.Errorf("a: %v is not equal to b: %v", a, b)
	}
}
