// Work in progress

package optimga

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
)

type AlgoGAES struct {
	AlgoPop
	pCross     float64
	pMut       float64
	nbChildren int
}

func NewAlgoGAES(pCross, pMut float64, maxGeneration int, log *log.Logger) (*AlgoGAES, error) {
	algo := &AlgoGAES{pCross: pCross, pMut: pMut}

	algo.maxGeneration = maxGeneration
	algo.logger = log

	return algo, nil
}

// crossover will create variation from 2 individuals
func (a *AlgoGAES) crossover(offspring *Pop, pcross float64) {
	for i := 0; i < offspring.size()/2; i++ {
		if rand.Float64() <= pcross {
			offspring.genotypes[2*i].setEval(false)
			offspring.genotypes[2*i+1].setEval(false)
			offspring.genotypes[2*i].cross(offspring.genotypes[2*i+1])
		}
	}
}

// mutation will create variation from 1 individual
func (a *AlgoGAES) mutation(offspring *Pop, pmut float64) {
	for i := 0; i < offspring.size(); i++ {

		if rand.Float64() <= pmut {
			offspring.genotypes[i].setEval(false)
			offspring.genotypes[i].mutate()
		}
	}
}

// replacement takes the best among the children and the parents to set the new population
func (a *AlgoGAES) replacement(offspring *Pop) {
	mixPop := newPopMix(a.pop, offspring)

	sort.Sort(mixPop)

	a.pop.genotypes = mixPop.genotypes[0:a.pop.size()]
	a.pop.fitness = mixPop.fitness[0:a.pop.size()]
}

// Run using algorithm inspired by typical evolution strategies (mu+lambda)-ES
func (a *AlgoGAES) Run() {
	fmt.Println("Ctrl-C to end algorithm before the planned end")
	a.WaitSignal()

	var offspring = new(Pop)
	offspring.reset(a.nbChildren)
	offspring.initThreadsParam(a.nbThreads)

	for gen := 0; gen < a.maxGeneration; gen++ {

		a.generationIndex = gen

		a.selecter.Select(a.pop, offspring)

		a.crossover(offspring, a.pCross)

		a.mutation(offspring, a.pMut)

		a.evalThisPop(offspring)

		a.replacement(offspring)

		// Set best individual
		a.SetBestResult(a.pop.genotypes[0], a.pop.fitness[0])

		a.logger.Printf("generation : %d, best fitness: %f\n", gen, a.pop.fitness[0])

		if a.isStopRequired == true {
			a.state = Stopped
			break
		}
	}

	fmt.Println()
	a.logger.Println("end of ES algorithm, best individual:")
	a.logger.Printf("%+v\n", a.pop.genotypes[0].GetGenes())
}
