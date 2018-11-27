// Work in progress

package optimga

import (
	"fmt"
	"math"
	"math/rand"
)

const nbChildrenRate = 7

type RealES struct {
	genotypeCommon
	genes        []float64
	steps        []float64 // Sigma standard deviations
	rangeGene    float64
	tglobal      float64
	tlocal       float64
	sizeGenotype int // needed for tcglobal and tlocal computation
	fitnessFunc  func([]float64) Fitness
}

func (r *RealES) init() {

	r.genes = make([]float64, r.sizeGenotype)
	r.steps = make([]float64, r.sizeGenotype)

	for i := range r.genes {
		r.genes[i] = rand.Float64() * r.rangeGene
		r.steps[i] = rand.Float64()
	}

	r.tglobal = 1 / math.Sqrt(2*float64(r.sizeGenotype))
	r.tlocal = 1 / math.Sqrt(2*math.Sqrt(float64(r.sizeGenotype)))
}

func (r *RealES) String() string {

	var result string

	result += fmt.Sprintf("genes: ")

	for i := range r.genes {
		result += fmt.Sprintf("%f ", r.genes[i])
	}

	result += fmt.Sprintf("\n\nstdev: ")

	for i := range r.steps {
		result += fmt.Sprintf("%f ", r.steps[i])
	}

	return result
}

func (r *RealES) checkBoundaries(gene float64) float64 {
	if gene < -r.rangeGene {
		gene = -r.rangeGene
	} else if gene > r.rangeGene {
		gene = r.rangeGene
	}

	return gene
}

func (r *RealES) cross(gen Genotype) {
	for i := range r.genes {

		weight := rand.Float64()
		var r2 = gen.(*RealES)

		rgene := r.genes[i]
		r2gene := r2.genes[i]

		r.genes[i] = rgene*weight + r2gene*(1-weight)
		r2.genes[i] = rgene*(1-weight) + r2gene*(weight)

		r.genes[i] = r.checkBoundaries(r.genes[i])
		r2.genes[i] = r2.checkBoundaries(r.genes[i])

	}
}

func (r *RealES) mutate() {
	mindev := 0.0

	global := r.tglobal * rand.NormFloat64()

	for i := range r.steps {

		r.steps[i] *= math.Exp(global + r.tlocal*rand.NormFloat64())

		if r.steps[i] < mindev {
			r.steps[i] = mindev
		} else if r.steps[i] > r.rangeGene/2 {
			r.steps[i] = r.rangeGene / 2
		}
	}

	for i := range r.genes {
		r.genes[i] += r.steps[i] * rand.NormFloat64()
		r.genes[i] = r.checkBoundaries(r.genes[i])
	}
}

func (r *RealES) GetGenes() []float64 {

	return r.genes
}

func (r *RealES) clone() Genotype {

	rclone := new(RealES)

	*rclone = *r

	rclone.genes = make([]float64, len(r.genes))
	copy(rclone.genes, r.genes)

	rclone.steps = make([]float64, len(r.steps))
	copy(rclone.steps, r.steps)

	return rclone
}

func (r *RealES) eval() Fitness {
	return r.fitnessFunc(r.genes)
}

func MakeRealESPop(algo *AlgoGAES, rangeGene float64, sizeGenotype int, sizePop int, fitnessFunc func([]float64) Fitness) {

	algo.pop = new(Pop)
	algo.nbChildren = sizePop * nbChildrenRate

	for i := 0; i < sizePop; i++ {
		realES := new(RealES)
		realES.sizeGenotype = sizeGenotype
		realES.rangeGene = rangeGene
		realES.fitnessFunc = fitnessFunc

		realES.init()

		algo.pop.genotypes = append(algo.pop.genotypes, realES)
	}
}
