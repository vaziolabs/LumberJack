package _01

import (
	"engine"
	"testing"
)

func TestUnwrapPath(t *testing.T) {
	path := "parent::child::{fork1,fork2}::{branch1,branch2::[declaration(1,2,3)]}::your_mom"
	seed, err := UnwrapPath(path)
	if err != nil {
		t.Error("Expected path to be unwrapped")
	}
	engine.Log(engine.DebugLevel, " > UnwrapPathTest")
	printSeed(seed, 0)
	engine.Log(engine.DebugLevel, " ")
}
