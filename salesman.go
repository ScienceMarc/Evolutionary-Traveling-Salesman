package main

import (
	"math/rand"
)

type Salesman struct {
	path    []string
	fitness float64
}

func (s *Salesman) initializePath(randomizer *rand.Rand) {
	s.path = make([]string, numberOfLocations)
	for i := range s.path {
		s.path[i] = string(alphabet[i])
	}
	//shuffle
	for i := range s.path {
		r := randomizer.Intn(len(s.path))
		s.path[i], s.path[r] = s.path[r], s.path[i] //swap
	}
}

func (s *Salesman) evaluate() float64 {
	totalDistance := 0
	for i := 0; i < len(s.path)-1; i++ {
		totalDistance += pointMap[s.path[i]+s.path[i+1]]
	}
	s.fitness = 1e9 / float64(totalDistance) //arbitrary fitness function
	return s.fitness
}

func (s *Salesman) crossover(p Salesman, randomizer *rand.Rand) Salesman { //Using order crossover
	offspring := Salesman{}
	offspring.path = make([]string, numberOfLocations)
	p1 := randomizer.Intn(numberOfLocations)
	p2 := randomizer.Intn(numberOfLocations)
	start, end := 0, 0
	if p1 < p2 {
		start = p1
		end = p2
	} else {
		start = p2
		end = p1
	}

	for i := start; i <= end; i++ {
		offspring.path[i] = s.path[i]
	}
	containsPoint := func(slice []string, letter string) bool {
		for _, v := range slice {
			if v == letter {
				return true
			}
		}
		return false
	}
	for _, v := range p.path {
		if !containsPoint(offspring.path, v) {
			for i := range offspring.path {
				if offspring.path[i] == "" {
					offspring.path[i] = v
					break
				}
			}
		}
	}
	for i := 0; i < 100; i++ {
		if randomizer.Float64() < mutationRate {
			r1 := randomizer.Intn(len(offspring.path))
			r2 := randomizer.Intn(len(offspring.path))
			offspring.path[r1], offspring.path[r2] = offspring.path[r2], offspring.path[r1]
		}
	}

	//fmt.Println(offspring.path)

	return offspring
}
