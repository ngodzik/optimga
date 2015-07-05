// Work in progress

package optimga

import (
	"math/rand"
	"sort"
)

type Tournament struct {
	size int
}

func NewTournament(size int) *Tournament {
	if size < 2 {
		size = 2
	}

	return &Tournament{size: size}
}

type tFitness struct {
	fitness  []Fitness
	indexPop []int
}

func (f *tFitness) init(size int) {
	f.fitness = make([]Fitness, size)
	f.indexPop = make([]int, size)
}

func (f *tFitness) set(fitness Fitness, index int, pos int) {
	f.fitness[pos] = fitness
	f.indexPop[pos] = index
}

func (f tFitness) Len() int { return len(f.fitness) }
func (f tFitness) Less(i, j int) bool {
	return f.fitness[i] < f.fitness[j]
}

func (f tFitness) Swap(i, j int) {
	f.fitness[i], f.fitness[j] = f.fitness[j], f.fitness[i]
	f.indexPop[i], f.indexPop[j] = f.indexPop[j], f.indexPop[i]
}

// Select Tournament selection, select randomly 2 parents and take the best of them for child
func (t *Tournament) Select(parents *Pop, offspring *Pop) {

	var fitarray tFitness
	fitarray.init(t.size)

	for i := 0; i < offspring.size(); i++ {
		for ti := 0; ti < t.size; ti++ {
			index := rand.Intn(parents.size())
			fitarray.set(parents.fitness[index], index, ti)
		}

		sort.Sort(fitarray)

		index := fitarray.indexPop[0]

		offspring.genotypes[i] = parents.genotypes[index].clone()
		offspring.fitness[i] = parents.fitness[index]
	}
}
