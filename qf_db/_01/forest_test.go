package _01

import (
	"engine"
	"testing"
)

func TestUnwrapSimple(t *testing.T) {
	wrapped := "simple::nested::branches"
	branches, err := UnwrapAll(wrapped)
	engine.Log(engine.InfoLevel, "Branches: %v, Error: %v", branches, err)

	if len(branches) != 1 {
		t.Errorf("Expected 1 branch, got %d", len(branches))
	}
}

func TestUnwrapComplex(t *testing.T) {
	wrapped := "{fork1,fork2}::{branch1,branch2::[declaration(1,2,3)]}::your_mom"
	f := NewForest("Complex")

	branches, err := f.PopulateBranches(wrapped)
	engine.Log(engine.InfoLevel, "Branches: %v, Error: %v", branches, err)

	if len(branches) != 2 {
		t.Errorf("Expected 2 forks, got %d", len(branches))
	}

	if len(branches[0].Branches) != 2 {
		t.Errorf("Expected 2 branches, got %d", len(branches[0].Branches))
	}

	if len(branches[0].Branches["branch1"].Branches) != 1 {
		t.Errorf("Expected 1 branch, got %d", len(branches[0].Branches["branch1"].Branches))
	}

}

func TestMultipleTrees(t *testing.T) {
	forest := NewForest("My First Forest")

	health_tree := forest.GetTree("Health")
	health_tree.GrowBranch("Diet")
	health_tree.GrowBranch("Exercise")
	health_tree.GrowBranch("Sleep")
	health_tree.PrintAll()

	finance_tree := forest.GetTree("Finance")
	finance_tree.GrowBranches([]string{"Budget", "Investing", "Retirement"})
	engine.Log(engine.InfoLevel, "Finance Tree: %v", finance_tree)

	relationship_tree := forest.GetTree("Relationships")
	relationship_tree.GrowBranches([]string{"Family", "Friends"})
	engine.Log(engine.InfoLevel, "Relationships Tree: %v", relationship_tree)

	business_tree := forest.GetTree("Business")
	business_tree.GrowBranches([]string{"Marketing", "Operations", "Finance"})
	engine.Log(engine.InfoLevel, "Business Tree: %v", business_tree)

	sales := NewBranch("Sales")
	finance_tree.AddBranch(sales)
	business_tree.AddBranch(sales)
	engine.Log(engine.InfoLevel, "Sales Branch: %v", sales)

	if len(forest.Trees) != 4 {
		t.Errorf("Expected 4 trees, got %d", len(forest.Trees))
	}

	engine.Log(engine.TraceLevel, " > MultipleTreesTest")
	engine.Log(engine.TraceLevel, " ~~~ Forest ~~~")
	forest.Log()
}

func TestMultipleParents(t *testing.T) {
	forest := NewForest("MultipleParents")

	health_tree := forest.GetTree("Health")
	fitness_branch := health_tree.GrowBranch("Fitness")

	career_tree := forest.GetTree("Career")
	training_branch := career_tree.GrowBranch("Training")

	workout_routing := NewBranch("Workout Routine").AddParents([]*Branch{fitness_branch, training_branch})

	engine.Log(engine.InfoLevel, "Workout Routine: %v", workout_routing)

	if len(workout_routing.Parents) != 2 {
		t.Errorf("Expected 2 parents, got %d", len(workout_routing.Parents))
	}

	engine.Log(engine.TraceLevel, " > MultipleParentsTest")
	forest.Log()
}

func TestAbstractForest(t *testing.T) {
	f, err := NewForest("test_forest").TreeConstructor("simple::nested::branches")

	engine.Log(engine.DebugLevel, " >")
	f.Log()
	engine.Log(engine.DebugLevel, " >")

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if f.Name != "test_forest" {
		t.Errorf("Expected name to be 'test_forest', got %s", f.Name)
	}

	if len(f.Trees) != 1 {
		t.Errorf("Expected 1 tree, got %d", len(f.Trees))
	}

	tree := f.GetTree("simple")
	if tree.Name != "simple" {
		t.Errorf("Expected name to be 'simple', got %s", tree.Name)
	}

	if len(tree.Branches) != 1 {
		t.Errorf("Expected 1 branch, got %d", len(tree.Branches))
	}
}

func TestComplexTree(t *testing.T) {
	f, err := NewForest("test_forest").TreeConstructor("{rootA, rootB}::{branchA, branchB}::nested")

	engine.Log(engine.DebugLevel, " >")
	f.Log()
	engine.Log(engine.DebugLevel, " >")

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if len(f.Trees) != 2 {
		t.Errorf("Expected 2 trees, got %d", len(f.Trees))
	}

	tree := f.GetTree("rootA")
	if tree.Name != "rootA" {
		t.Errorf("Expected name to be 'rootA', got %s", tree.Name)
	}

	if len(tree.Branches) != 2 {
		t.Errorf("Expected 2 branches, got %d", len(tree.Branches))
	}
}

func TestDistributionList(t *testing.T) {
	f, err := NewForest("test_forest").TreeConstructor("tree::[menial(trivial,complex)]::nested")
	engine.Log(engine.DebugLevel, " >")
	f.Log()
	engine.Log(engine.DebugLevel, " >")

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if len(f.Trees) != 1 {
		t.Errorf("Expected 1 tree, got %d", len(f.Trees))
	}

	root := f.GetTree("root")
	if len(root.Branches) != 2 {
		t.Errorf("Expected 2 branch, got %d", len(f.GetTree("root").Branches))
	}

	if root.Branches["trivial"] == nil {
		t.Errorf("Expected branch 'menial' to exist")
	}

	if root.Branches["complex"] == nil {
		t.Errorf("Expected branch 'complex' to exist")
	}

}
