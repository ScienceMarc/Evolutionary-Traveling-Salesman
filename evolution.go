package main

import (
	"math/rand"
	"sort"
	"sync"
)

type Population struct {
	members    []Salesman
	matingPool []Salesman
}

func (p *Population) createPopulation(randomizer *rand.Rand) {
	p.members = make([]Salesman, generationSize)
	for i := range p.members {
		p.members[i].initializePath(randomizer)
	}
}

func (p *Population) evaluate() ([]string, float64, float64) {
	totalFitness := 0.0
	highestFitness := 0.0
	bestPath := make([]string, numberOfLocations)
	var wg sync.WaitGroup
	for i := range p.members {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fitness := p.members[i].evaluate()
			totalFitness += fitness
			if fitness > highestFitness {
				highestFitness = fitness
				bestPath = p.members[i].path
			}
		}(i)
	}
	wg.Wait()
	/*
		fmt.Println("Average fitness:", totalFitness/float64(len(p.members)))
		fmt.Println("Best fitness:", highestFitness)
		fmt.Println("Best path:", bestPath)*/
	return bestPath, highestFitness, totalFitness / float64(len(p.members))
}
func (p *Population) selection() []Salesman { //using truncation selection
	sort.Slice(p.members, func(i, j int) bool {
		return p.members[i].fitness > p.members[j].fitness
	})
	p.matingPool = make([]Salesman, 0)
	for i := 0; i < len(p.members)/2; i++ {
		p.matingPool = append(p.matingPool, p.members[i])
	}
	return p.matingPool[:20]
}
func (p *Population) crossover(randomizer *rand.Rand) {
	for i := range p.matingPool {
		offspring := p.matingPool[i].crossover(p.matingPool[randomizer.Intn(len(p.matingPool))], randomizer)
		p.members = append(p.members, offspring) //Every member of the mating pool mates with another, possibly themselves
	}
	remove := func(s []Salesman, i int) []Salesman {
		s[len(s)-1], s[i] = s[i], s[len(s)-1]
		return s[:len(s)-1]
	}
	for len(p.members) > generationSize { //Keep removing random members until population size is appropriate
		p.members = remove(p.members, randomizer.Intn(len(p.members)))
	}
}
