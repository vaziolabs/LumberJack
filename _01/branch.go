package _01

import (
	"engine"
	"fmt"
)

type Branch struct {
	Name          string              `json:"name"`          // equivalent to the name of a table in a database
	Branches      map[string]*Branch  `json:"forks"`         // These are Forks in the branch
	Parents       []*Branch           `json:"parents"`       // These are the branches that lead to this branch
	Distributions map[string][]string `json:"distributions"` // These are the branches that are distributed from this branch
	Flower        map[string]*Flower  `json:"flowers"`       // These are the leaves in the branch
	Root          *Branch             `json:"root"`          // This is the root of the branch
}

func (*Branch) Prototype() *Branch {
	b := &Branch{}
	b.Name = "Prototype"
	b.Branches = make(map[string]*Branch)
	b.Parents = make([]*Branch, 0)
	b.Distributions = make(map[string][]string)
	b.Flower = make(map[string]*Flower)
	return b
}

func UnwrapAll(path string) ([]*Branch, error) {
	// Logic needs to move from left to right to unwrap the path because we can have nested branches
	// Needs to return a tree where [] is a distribution and {} is a fork
	// Distributions need to have a declaration if they are not already in the distributions map
	return nil, nil
}

func (b *Branch) String() string {
	print_str := fmt.Sprintf("{\n\tName: %s", b.Name)

	print_str += fmt.Sprintf("\n\tParents(%d): [", len(b.Parents))
	for _, parent := range b.Parents {
		print_str += parent.Name
		if parent != b.Parents[len(b.Parents)-1] {
			print_str += ", "
		}
	}
	print_str += "]"

	print_str += fmt.Sprintf("\n\tDistributions(%d): [", len(b.Distributions))
	i := 0
	i_length := len(b.Distributions)
	for key, dist := range b.Distributions {
		print_str += fmt.Sprintf("{%s: %s}", key, dist)
		i++
		if i != i_length {
			print_str += ", "
		}
	}
	print_str += "]"

	i = 0
	i_length = len(b.Branches)
	print_str += fmt.Sprintf("\n\tBranches(%d): [", len(b.Branches))
	for _, branch := range b.Branches {
		print_str += branch.Name
		i++
		if i != i_length {
			print_str += ", "
		}
	}
	print_str += "]"

	i = 0
	i_length = len(b.Flower)
	print_str += fmt.Sprintf("\n\tFlowers(%d): [", len(b.Flower))
	for _, flower := range b.Flower {
		print_str += flower.Name
		i++
		if i != i_length {
			print_str += ", "
		}
	}
	print_str += "]"

	print_str += fmt.Sprintf("\n\tRoot: %s", b.Root.Name)
	print_str += "\n}"
	return print_str
}

func (b *Branch) FindBranch(name string) (*Branch, error) {
	if branch, ok := b.Branches[name]; ok {
		return branch, nil
	}

	for _, child := range b.Branches {
		if branch, err := child.FindBranch(name); err == nil {
			return branch, nil
		}
	}

	return nil, fmt.Errorf("Branch %s not found", name)
}

func (b *Branch) PrintAll() {
	engine.Log(engine.DebugLevel, "Branch: %s", b.String())

	for _, branch := range b.Branches {
		branch.PrintAll()
	}

	for _, flower := range b.Flower {
		engine.Log(engine.DebugLevel, "Flower: %v", flower)
	}
}

func (b *Branch) Traverse(keys []string) (*Branch, error) {
	if len(keys) == 0 {
		return b, nil
	}

	next, err := b.GetBranch(keys[0])

	if err != nil {
		return nil, err
	}

	return next.Traverse(keys[1:])
}

func (b *Branch) AddDistribution(key string, branches []string) {
	b.Distributions[key] = branches
}

func (b *Branch) GetDistributions(key string) []string {
	return b.Distributions[key]
}

func (b *Branch) AddParent(branch *Branch) {
	// Make sure the parent isn't already in the list
	for _, parent := range b.Parents {
		if parent == branch {
			return
		}
	}

	b.Root = branch.Root
	b.Parents = append(b.Parents, branch)

	// We need to check if the parent already has this branch as a child
	for _, child := range branch.Branches {
		if child == b {
			return
		}
	}

	branch.AddBranch(b)
}

func (b *Branch) AddParents(branches []*Branch) *Branch {
	for _, branch := range branches {
		b.AddParent(branch)
	}
	return b
}

func (b *Branch) GrowBranches(keys []string) []*Branch {
	var branches []*Branch

	for _, key := range keys {
		branches = append(branches, b.GrowBranch(key))
	}

	return branches
}

func (b *Branch) GrowBranch(key string) *Branch {
	if branch, ok := b.Branches[key]; ok {
		return branch
	}

	new_branch := NewBranch(key)
	new_branch.Parents = append(new_branch.Parents, b)
	new_branch.Root = b.Root
	b.Branches[key] = new_branch
	return new_branch
}

func (b *Branch) AddBranches(branches []*Branch) []*Branch {
	var branches_out []*Branch

	for _, branch := range branches {
		branches_out = append(branches, b.AddBranch(branch))
	}

	return branches_out
}

// Do we want to add support for multiple branches of the same name?
func (b *Branch) AddBranch(branch *Branch) *Branch {
	if _, ok := b.Branches[branch.Name]; ok {
		return b.Branches[branch.Name]
	}

	branch.AddParent(b)
	branch.Root = b.Root
	b.Branches[branch.Name] = branch
	return branch
}

func (b *Branch) GetBranch(name string) (*Branch, error) {
	if branch, ok := b.Branches[name]; ok {
		return branch, nil
	}
	return nil, fmt.Errorf("Branch %s not found", name)
}

// This clears the branch and all its children so there's no fragments - Maybe this is a bit too much? I need proof either way
func (b *Branch) Clear() {
	for _, branch := range b.Branches {
		branch.Clear()
		delete(b.Branches, branch.Name)
	}

	b.Branches = make(map[string]*Branch)
	b.Flower = make(map[string]*Flower)
}

// This actually removes the branch but leaves any potential children
func (b *Branch) Cut(name string) {
	b.Branches[name].Clear()
	delete(b.Branches, name)
}

// This scrapes all the values off the branch
func (b *Branch) PruneAll() {
	for _, branch := range b.Branches {
		branch.PruneAll()
	}

	for _, flower := range b.Flower {
		delete(b.Flower, flower.Name)
	}
}

func (b *Branch) Prune(name string) {
	b.Branches[name].PruneAll()
}

func (b *Branch) RemoveBranch(branch Branch) {
	delete(b.Branches, branch.Name)
}

func (b *Branch) GrowFlower(name string) *Flower {
	// Check if the flower already exists
	if flower, ok := b.Flower[name]; ok {
		return flower
	}
	new_flower := NewFlower(name)
	b.Flower[name] = new_flower
	return new_flower
}

func (b *Branch) AddFlower(flower *Flower) {
	b.Flower[flower.Name] = flower
}

func (b *Branch) RemoveFlower(name string) {
	delete(b.Flower, name)
}

func (b *Branch) GetFlower(name string) *Flower {
	return b.Flower[name]
}

func NewBranch(name string) *Branch {
	return &Branch{
		Name:     name,
		Branches: make(map[string]*Branch),
		Flower:   make(map[string]*Flower),
	}
}

func (b *Branch) PopulateBranches(path string) {
	branches, err := UnwrapAll(path)

	if err != nil {
		engine.Log(engine.ErrorLevel, "Error unwrapping branches: %v", err)
		return
	}

	for _, branch := range branches {
		engine.Log(engine.DebugLevel, "Branch: %v", branch)
	}
}
