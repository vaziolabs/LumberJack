package _01

import (
	"engine"
	"errors"
	"strings"
)

type Forest struct {
	Name         string              `json:"name"`         // equivalent to the name of a database
	Trees        map[string]*Branch  `json:"trees"`        // equivalent to a lookup table in a database
	Distribution map[string][]string `json:"distribution"` // equivalent to a distribution list in a database
	Tulips       map[string]*Flower  `json:"tulips"`       // equivalent to a leaf node in a database - used for a fast lookup table Needs to be a hashtable
}

func NewForest(name string) *Forest {
	return &Forest{
		Name:  name,
		Trees: make(map[string]*Branch),
	}
}

func (f *Forest) addTree(name string) *Branch {
	engine.Log(engine.TraceLevel, " > addTree: %s", name)

	root := NewBranch(name)
	f.Trees[name] = root
	return root
}

func (f *Forest) GetTree(name string) *Branch {
	engine.Log(engine.TraceLevel, " > GetTree: %s", name)

	if tree, ok := f.Trees[name]; ok {
		return tree
	}

	return f.addTree(name)
}

func (f *Forest) RemoveForest(name string) {
	engine.Log(engine.TraceLevel, " > RemoveForest: %s", name)
	delete(f.Trees, name)
}

func (f *Forest) Log() {
	engine.Log(engine.TraceLevel, "Forest: %s", f.Name)

	for _, tree := range f.Trees {
		tree.PrintAll()
	}
}

func (f *Forest) TreeConstructor(absolute_path string) (*Forest, error) {
	engine.Log(engine.TraceLevel, " > TreeConstructor: %s", absolute_path)
	return f, nil
}

// TODO: Add Test
func (f *Forest) SearchBranches(name string) (*Branch, error) {
	for key, branch := range f.Trees {
		if key == name {
			return branch, nil
		}

		b, e := branch.FindBranch(name)

		if e == nil {
			return b, nil
		}
	}

	return nil, errors.New("Branch not found")
}

func (f *Forest) GetBranch(name string) (*Branch, error) {
	if branch, ok := f.Trees[name]; ok {
		return branch, nil
	}
	return nil, errors.New("Branch not found")
}

func (f *Forest) ClimbBranch(branches string) (*Branch, error) {
	keys := strings.Split(branches, "::")
	branch := f.Trees[keys[0]]

	for _, key := range keys {
		branch = branch.Branches[key]
		if branch == nil {
			return nil, errors.New("Branch not found")
		}
	}

	return branch, nil
}

// Tree Should Exist already AND this should return the branch pointer
func (f *Forest) PopulateBranches(path string) ([]*Branch, error) {
	//f.Distribution[distribution] = branches
	engine.Log(engine.TraceLevel, " > PopulateBranches: %s", path)

	return nil, nil
}

func (f *Forest) PlantTree(name string) {
	engine.Log(engine.TraceLevel, " > PlantTree: %s", name)
}
