package main

import (
	"testing"
)

func TestProcessor_start(t *testing.T) {
	p := NewProcessor()
	p.ip = 100
	p.start()

	if p.depth != 1 {
		t.Errorf("wrong depth: %v", p.debug)
	}

	if p.getCorrespondingStart(p.depth) != 100 {
		t.Errorf("wrong starting position: %v", p.getCorrespondingStart(p.depth))
	}
}
