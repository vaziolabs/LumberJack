package _01

import (
	"fmt"
	"regexp"
	"strings"
)

// Parent -> It can have a fork or a distribution
type Seed struct {
	name          string
	forks         map[string]*Seed
	distributions map[string][]string
}

func (s *Seed) AddDistribution(dist string, forks []string) {
	for _, fork := range forks {
		s.distributions[dist] = append(s.distributions[dist], fork)
	}

}

func (s *Seed) AddFork(fork_str string) {
	s.forks[fork_str] = &Seed{
		name:          fork_str,
		forks:         make(map[string]*Seed),
		distributions: make(map[string][]string),
	}
}

func parseSeed(input string) *Seed {
	seed := &Seed{
		name:          "",
		forks:         make(map[string]*Seed),
		distributions: make(map[string][]string),
	}
	parseHelper(input, seed)
	return seed
}

func parseHelper(input string, seed *Seed) {
	tokenRegex := regexp.MustCompile(`(::|\{|\}|\[|\]|\(|\)|,)`)
	tokens := tokenRegex.Split(input, -1)
	tokens = tokens[:len(tokens)-1] // This ends up splitting all the words without regarding the type of token

	currentFork := seed
	var stack []*Seed

	for i := 0; i < len(tokens); i++ {
		token := strings.TrimSpace(tokens[i])

		if token == "::" {
			continue
		} else if token == "{" {
			// Start a new fork context
			stack = append(stack, currentFork)
			currentFork = &Seed{
				forks:         make(map[string]*Seed),
				distributions: make(map[string][]string),
			}
		} else if token == "}" {
			// End current fork context
			last := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			for k := range currentFork.forks {
				last.forks[k] = currentFork.forks[k]
			}
			currentFork = last
		} else if token == "[" {
			// Start a new distribution context
			i++
			distributionName := strings.TrimSpace(tokens[i])
			i += 2 // skip over '('
			var distValues []string
			for tokens[i] != ")" {
				distValues = append(distValues, strings.TrimSpace(tokens[i]))
				i++
			}
			currentFork.distributions[distributionName] = distValues
		} else if token != "," {
			// Regular fork or root name
			if currentFork.name == "" {
				currentFork.name = token
			} else {
				newFork := &Seed{
					name:          token,
					forks:         make(map[string]*Seed),
					distributions: make(map[string][]string),
				}
				currentFork.forks[token] = newFork
				currentFork = newFork
			}
		}
	}
}

func printSeed(seed *Seed, indent int) {
	prefix := strings.Repeat("  ", indent)
	fmt.Printf("%sSeed: %s\n", prefix, seed.name)
	for k, v := range seed.forks {
		fmt.Printf("%s  Fork: %s\n", prefix, k)
		printSeed(v, indent+1)
	}
	for k, v := range seed.distributions {
		fmt.Printf("%s  Distribution: %s -> %v\n", prefix, k, v)
	}
}

func UnwrapPath(path string) (*Seed, error) {
	return parseSeed(path), nil
}
