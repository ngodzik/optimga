// Work in progress

package optimga

import "log"

type Pop struct {
	genotypes    []Genotype
	fitness      []Fitness
	isMaximizing bool
	tpops        []threadPop
}

func (p Pop) Len() int { return len(p.fitness) }

func (p Pop) Less(i, j int) bool {
	var returnBool bool

	if p.isMaximizing {
		returnBool = p.fitness[i] < p.fitness[j]
	} else {
		returnBool = p.fitness[i] < p.fitness[j]
	}
	return returnBool
}

func (p Pop) Swap(i, j int) {
	p.genotypes[i], p.genotypes[j] = p.genotypes[j], p.genotypes[i]
	p.fitness[i], p.fitness[j] = p.fitness[j], p.fitness[i]
}

func (p *Pop) size() int {
	return len(p.genotypes)
}

func (p *Pop) Maximize(isMaximizing bool) {
	p.isMaximizing = isMaximizing
}

func (p *Pop) reset(size int) {
	p.genotypes = make([]Genotype, size)
	p.fitness = make([]Fitness, size)
}

func (p *Pop) clone() *Pop {
	copyPop := new(Pop)
	copyPop.reset(len(p.genotypes))

	copy(copyPop.genotypes, p.genotypes)
	copy(copyPop.fitness, p.fitness)

	return copyPop
}

func newPopMix(pop1, pop2 *Pop) *Pop {
	mixPop := new(Pop)

	mixPop.isMaximizing = pop1.isMaximizing

	mixPop.reset(pop1.size() + pop2.size())

	copy(mixPop.genotypes, pop1.genotypes)
	copy(mixPop.fitness, pop1.fitness)

	copy(mixPop.genotypes[pop1.size():], pop2.genotypes)
	copy(mixPop.fitness[pop1.size():], pop2.fitness)

	return mixPop
}

func (p *Pop) display(logger *log.Logger) {
	for i, gen := range p.genotypes {
		logger.Printf("%d: %s\n", i, gen)
	}
}

func (p *Pop) initThreadsParam(nbThreads int) {
	p.tpops = make([]threadPop, nbThreads)

	step := len(p.genotypes) / nbThreads
	sumindexes := 0

	for i := 0; i < nbThreads; i++ {
		p.tpops[i].index = i * step
		if i != nbThreads-1 {
			p.tpops[i].size = step
			sumindexes += step
		} else {
			p.tpops[i].size = len(p.genotypes) - sumindexes
		}
		p.tpops[i].result = make(chan Fitness)
	}
}

type threadPop struct {
	index  int
	size   int
	result chan Fitness
}
