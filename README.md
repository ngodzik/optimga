# Optimga

## Evolutionary optimization using Go

Work in progress, for now, only one part of the evolution strategies are implemented

for example, try to optimize the sphere function:

```
func sphereFunction(genes []float64) optimga.Fitness {

	var fitness float64

	for i := range genes {
		fitness += math.Pow(genes[i], 2)
	}

	return optimga.Fitness(fitness)
}
```

with the following example main program:

```
func main() {

	logger := log.New(os.Stdout, fmt.Sprintf("[opti %d] ", os.Getpid()), log.LstdFlags)

	// The stopping criterium is generations max
	// In evolution strategies, the crossover may be useless, it is set to 0 here
	opti, err := optimga.NewAlgoGAES(0.0, // crossover probability
		1.0, // mutation probability
		500, // generations max
		logger)

	if err != nil {
		logger.Fatalf("Could not create GAES algorithm: %s", err)
	}

	opti.SetRandSeed(time.Now().Unix())

	opti.SetSelecter(optimga.NewTournament(2))

	// Create the population with the choosen genotype type (reals evolution strategies in this case)
	// We use 50 genes, the sphere function is so x1^2 + x2^2 + ... + x50^2
	optimga.MakeRealESPop(opti,
		5,    // range
		50,   // number of genes
		1000, // population size
		sphereFunction)

	opti.SetNbThreads(4)

	opti.Init()

	opti.Run()
}
```
