// Work in progress

package optimga

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sort"
	"sync"
)

type State int

const (
	NotRunning State = iota + 1
	Initializing
	Started
	Stopped
)

type Fitness float64

type Algo interface {
	init()
	SetFitness(f func(Genotype, int, int) Fitness)
	SetSelecter(s Selecter)
	appendg(g Genotype)
	GetCurrentBestResult() (Genotype, Fitness)
	setNbThreads(nbThreads int)
	Run()
	GetState() State
	Stop()
}

type algoOperator interface {
	crossover(offspring *Pop, pcross float64)
	mutation(offspring *Pop, pmut float64)
}

type AlgoPop struct {
	pop             *Pop
	algoOp          algoOperator
	selecter        Selecter
	pcross          float64
	pmut            float64
	isStopRequired  bool
	state           State
	bestGenotype    Genotype
	bestFitness     Fitness
	generationIndex int
	maxGeneration   int
	nbThreads       int
	seed            int64
	mutex           sync.Mutex
	logger          *log.Logger
}

func (a *AlgoPop) SetNbThreads(nbThreads int) {
	// Ignore previous settings return
	fmt.Println("Set nb Threads: ", nbThreads)
	runtime.GOMAXPROCS(nbThreads)
	a.nbThreads = nbThreads
}

func (a *AlgoPop) SetRandSeed(seed int64) {
	a.seed = seed
	rand.Seed(seed)
}

func (a *AlgoPop) Appendg(g Genotype) {

	a.pop.genotypes = append(a.pop.genotypes, g)
}

func (a *AlgoPop) Init() {
	a.state = Initializing

	a.pop.fitness = make([]Fitness, len(a.pop.genotypes))

	if a.nbThreads == 0 {
		a.nbThreads = 1
	}

	a.pop.initThreadsParam(a.nbThreads)

	a.evalPop()

	sort.Sort(a.pop)

	a.mutex.Lock()
	a.bestGenotype = a.pop.genotypes[0]
	a.bestFitness = a.pop.fitness[0]
	a.mutex.Unlock()

	a.state = Started
}

func (a *AlgoPop) GetState() State {
	return a.state
}

func (a *AlgoPop) GetCurrentGeneration() int {

	return a.generationIndex
}

// The evaluation of an individual can be time consuming, goroutines are used increase the speed of the generation computation
func (a *AlgoPop) evalPop() {
	a.evalThisPop(a.pop)
}

func (a *AlgoPop) evalPopPart(threadIndex int, pop *Pop) {

	index := pop.tpops[threadIndex].index
	nb := pop.tpops[threadIndex].size
	result := pop.tpops[threadIndex].result

	var sum Fitness

	for i := index; i < index+nb; i++ {
		if pop.genotypes[i].isEvaluated() == false {
			pop.fitness[i] = pop.genotypes[i].eval()
			pop.genotypes[i].setEval(true)
		}

		sum += pop.fitness[i]
	}

	result <- sum
}

func (a *AlgoPop) evalThisPop(pop *Pop) {

	// For now, a new thread is started even if nbThread == 1

	for i := 0; i < a.nbThreads; i++ {
		go a.evalPopPart(i, pop)
	}

	var sumfitness Fitness

	// Wait the end of all previous started threads
	for i := 0; i < a.nbThreads; i++ {
		sumfitness += <-pop.tpops[i].result
	}
}

func (a *AlgoPop) SetSelecter(selecter Selecter) {
	a.selecter = selecter
}

func (a *AlgoPop) Stop() {
	a.isStopRequired = true
}

func (a *AlgoPop) SetBestResult(bestGenotype Genotype, bestFitness Fitness) {
	a.mutex.Lock()
	a.bestGenotype = bestGenotype
	a.bestFitness = bestFitness
	a.mutex.Unlock()

	return
}

func (a *AlgoPop) GetCurrentBestResult() (Genotype, Fitness) {
	a.mutex.Lock()
	bestGenotype := a.bestGenotype
	bestFitness := a.bestFitness
	a.mutex.Unlock()

	return bestGenotype, bestFitness
}
