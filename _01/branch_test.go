package _01

import (
	"engine"
	"testing"
)

func TestNewBranch(t *testing.T) {
	b := NewBranch("testies")

	if b.Name != "test" {
		t.Errorf("Expected name to be 'test', got %s", b.Name)
	}

	engine.Log(engine.DebugLevel, " > NewBranchTest")
	engine.Log(engine.DebugLevel, "Branch: %s", b.Name)
	engine.Break()
}

func TestAddBranch(t *testing.T) {
	b := NewBranch("test")
	b.AddBranch(NewBranch("test2"))
	b.AddBranch(NewBranch("test1"))

	if b.Branches["test2"] == nil {
		t.Error("Expected branch to be added")
	}

	engine.Log(engine.DebugLevel, " > AddBranchTest")
	engine.Log(engine.DebugLevel, "N-Branches: %d", len(b.Branches))
	engine.Break()
}

func TestGrowBranch(t *testing.T) {
	b := NewBranch("test")
	b.GrowBranch("test2")

	if b.Branches["test2"] == nil {
		t.Error("Expected branch to be grown")
	}

	engine.Log(engine.DebugLevel, " > GrowBranchTest")
}

func TestGetBranch(t *testing.T) {
	b := NewBranch("test")
	test1 := NewBranch("test1")
	test2 := NewBranch("test2")
	nested := NewBranch("nested")
	b.AddBranch(test1)
	b.AddBranch(test2)
	test2.AddBranch(nested)

	branch, err := b.GetBranch("test2")

	if err != nil {
		t.Error("Expected branch to be found")
	}

	if branch.Name != "test2" {
		t.Errorf("Expected name to be 'test2', got %s", branch.Name)
	}

	engine.Log(engine.DebugLevel, " > GetBranchTest")
	engine.Log(engine.DebugLevel, "Branch: %s", branch.Name)
	engine.Break()
}

func TestRemoveBranch(t *testing.T) {
	b := NewBranch("test")
	test1 := NewBranch("test1")
	test2 := NewBranch("test2")
	nested := NewBranch("nested")
	b.AddBranch(test1)
	b.AddBranch(test2)
	test2.AddBranch(nested)
	if b.Branches["test"] != nil {
		t.Error("Expected branch to be removed")
	}
	engine.Log(engine.DebugLevel, " > RemoveBranchTest")
	engine.Log(engine.DebugLevel, "Branches: %d", len(b.Branches))
	engine.Break()
}

func TestCutBranch(t *testing.T) {
	b := NewBranch("test")
	test1 := NewBranch("test1")
	test2 := NewBranch("test2")
	nested := NewBranch("nested")
	b.AddBranch(test1)
	b.AddBranch(test2)
	test2.AddBranch(nested)
	b.Cut("test2")

	if b.Branches["test2"] != nil {
		t.Error("Expected branch to be pruned")
	}

	engine.Log(engine.DebugLevel, " > CutBranchTest")
	engine.Log(engine.DebugLevel, "Branches: %d", len(b.Branches))
	engine.Break()
}

func TestBranchFlower(t *testing.T) {
	b := NewBranch("test")
	b.AddFlower(NewFlower("test1"))
	b.AddFlower(NewFlower("test2"))

	if b.Flower["test1"] == nil {
		t.Error("Expected flower to be added")
	}

	engine.Log(engine.DebugLevel, " > BranchFlowerTest")
	engine.Log(engine.DebugLevel, "N-Flowers: %d", len(b.Flower))
	engine.Break()
}
