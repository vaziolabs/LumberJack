package _01

import (
	"engine"
	"testing"
)

func TestNewFlower(t *testing.T) {
	f := NewFlower("test")

	if f.Name != "test" {
		t.Errorf("Expected name to be 'test', got %s", f.Name)
	}

	engine.Log(engine.DebugLevel, " > NewFlowerTest")
	engine.Log(engine.DebugLevel, "Flower: %s", f.Name)
	engine.Break()
}

func TestAddFlower(t *testing.T) {
	b := NewBranch("test")
	f := NewFlower("test")
	b.AddFlower(f)

	if b.Flower["test"] == nil {
		t.Error("Expected flower to be added")
	}

	engine.Log(engine.DebugLevel, " > AddFlowerTest")
	engine.Log(engine.DebugLevel, "Flowers: %d", len(b.Flower))
	engine.Break()
}

func TestRemoveFlower(t *testing.T) {
	b := NewBranch("test")
	f := NewFlower("test")
	b.AddFlower(f)
	b.RemoveFlower(f.Name)

	if b.Flower["test"] != nil {
		t.Error("Expected flower to be removed")
		return
	}

	engine.Log(engine.DebugLevel, " > RemoveFlowerTest")
	engine.Log(engine.DebugLevel, "Flowers: %d", len(b.Flower))
	engine.Break()
}

func TestGetFlower(t *testing.T) {
	b := NewBranch("test")
	f := NewFlower("test")
	b.AddFlower(f)
	flower := b.GetFlower("test")

	if flower == nil {
		t.Error("Expected flower to be found")
		return
	}

	if flower.Name != "test" {
		t.Errorf("Expected name to be 'test', got %s", flower.Name)
		return
	}

	engine.Log(engine.DebugLevel, " > GetFlowerTest")
	engine.Log(engine.DebugLevel, "Flower: %s", flower.Name)
	engine.Break()
}
