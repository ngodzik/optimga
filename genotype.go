// Work in progress

package optimga

type Genotype interface {
	init()
	cross(gen Genotype)
	mutate()
	clone() Genotype
	eval() Fitness
	setEval(isEvalDone bool)
	isEvaluated() bool
	GetGenes() []float64
}

type genotypeCommon struct {
	isEvalDone bool
}

func (g *genotypeCommon) setEval(isEvalDone bool) {
	g.isEvalDone = isEvalDone
}

func (g *genotypeCommon) isEvaluated() bool {
	return g.isEvalDone
}
