package algorithms

import (
	"fmt"
	"hw2/tsp"
	"math/rand"
	"sort"
	"time"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/emirpasic/gods/trees/btree"
	"github.com/emirpasic/gods/utils"
)

type GeneticOptions struct {
	PopulationSize int
	NumGenerations int
	RolloverCount float64
	BTreeOrder int
	MutationChance float64
}

type chromosome struct {
	tour []int
	cost tsp.Weight
}

func GeneticLocalSearch(
	g *tsp.Graph, N int, options GeneticOptions,
) (tsp.Weight, int, []int) {
	R := rand.New(rand.NewSource(time.Now().UnixNano()))

	population := initializePopulation(g, N, R, options)
	for generation := 0; generation < options.NumGenerations; generation++ {
		roulette_wheel, roulette_to_chromosome := calculateFitness(population, options)
		new_population := createRolloverPopulation(population, options)
		populateWithChildren(roulette_wheel, roulette_to_chromosome, 
			new_population, g, N, R, options)
		population = new_population
	}
	k_, v_ := population.LeftKey(), population.LeftValue()
	return tsp.Weight(k_.(float64)), 
		options.PopulationSize * options.NumGenerations, v_.(chromosome).tour
}

/*
	create population btree and fill with random chromosomes
*/
func initializePopulation(
	g *tsp.Graph, N int, R *rand.Rand, options GeneticOptions,
) *btree.Tree {
	population := btree.NewWith(options.BTreeOrder, utils.Float64Comparator)
	i := 1
	for population.Size() < options.PopulationSize && i < options.PopulationSize * 10000 {
		tour := R.Perm(N)
		cost := g.GetTourCost(tour)
		population.Put(float64(cost), chromosome{
			tour: tour, cost: cost,
		})
		i++
	}
	if i == options.PopulationSize * 2 {
		panic(fmt.Sprintf("too big population size %d: %d\n", N, options.PopulationSize))
	}
	return population
}

/*
	calculate the fitness of entire population and create cumulative probabilities
*/
func calculateFitness(population *btree.Tree,
	options GeneticOptions,
) ([]float64, []chromosome) {
	roulette_wheel := make([]float64, options.PopulationSize)
	roulette_to_chromosome := make([]chromosome, options.PopulationSize)

	max_cost := population.RightKey().(float64)
	total_cost := 0.0

	iterator := population.Iterator()
	for iterator.Next() {
		total_cost += iterator.Key().(float64)
	}
	inv_total_fitness := 1.0 / (max_cost * float64(options.PopulationSize) - 
		total_cost)
	var cumulative_probability float64 = 0
	for i := 0; i < population.Size() && iterator.Prev(); i++ {
		k_, v_ := iterator.Key(), iterator.Value()
		k, v := k_.(float64), v_.(chromosome)
		roulette_wheel[i] = cumulative_probability
		cumulative_probability += (max_cost - k) * inv_total_fitness
		roulette_to_chromosome[i] = v
	}
	return roulette_wheel, roulette_to_chromosome
}

/*
	remove worst candidates from population
*/
func createRolloverPopulation(prev_population *btree.Tree,
	options GeneticOptions,
) *btree.Tree {
	population := btree.NewWith(options.BTreeOrder, utils.Float64Comparator)

	iterator := prev_population.Iterator()
	end := int(options.RolloverCount * float64(options.PopulationSize))
	for i := 0; i < end && iterator.Next(); i++ {
		population.Put(iterator.Key().(float64), iterator.Value().(chromosome))
	}
	return population
}

func getRandomParents(roulette []float64, roulette_to_chromosome []chromosome,
	R *rand.Rand, options GeneticOptions,
) (chromosome, chromosome) {
	roll1 := R.Float64()
	p1 := sort.Search(options.PopulationSize, 
		func(i int) bool {return roll1 <= roulette[i]}) - 1
	slice := -roulette[p1]
	if p1 + 1 == options.PopulationSize {
		slice += 1
	} else {
		slice += roulette[p1 + 1]
	}
	roll2 := R.Float64() * (1 - slice)
	if roll2 >= roulette[p1] {
		roll2 += slice
	}
	p2 := sort.Search(options.PopulationSize, 
			func(i int) bool {return roll2 <= roulette[i]}) - 1
	return roulette_to_chromosome[p1], roulette_to_chromosome[p2]
}

/*
	fill in rest of missing population with children, mutating
*/
func populateWithChildren(
	roulette []float64, 
	roulette_to_chromosome []chromosome,
	population *btree.Tree,
	g *tsp.Graph, N int, R *rand.Rand, options GeneticOptions,
) bool {
	// var mutex sync.Mutex
	// var wg sync.WaitGroup

	// add_child_thread := func(new_R *rand.Rand) {
	// 	defer wg.Done()
	// 	for {
	// 		p1, p2 := getRandomParents(roulette, roulette_to_chromosome, new_R, options)
	// 		pure_child := crossoverER(p1.tour, p2.tour, N, new_R)
	// 		mutated := mutateChromosome(pure_child, N, new_R, options)
	// 		child_chromosome := chromosome {
	// 			tour: mutated,
	// 			cost: g.GetTourCost(mutated),
	// 		}
	// 		mutex.Lock()
	// 		if population.Size() < options.PopulationSize {
	// 			population.Put(float64(child_chromosome.cost), child_chromosome)
	// 		} else {
	// 			mutex.Unlock()
	// 			break
	// 		}
	// 		mutex.Unlock()
	// 	}
	// }

	// for thread_i := 0; thread_i < options.NumThreads; thread_i++ {
	// 	wg.Add(1)
	// 	go add_child_thread(rand.New(rand.NewSource(R.Int63())))
	// }

	// wg.Wait()
	// return true

	for population.Size() < options.PopulationSize {
		p1, p2 := getRandomParents(roulette, roulette_to_chromosome, R, options)
		pure_child := crossoverER(p1.tour, p2.tour, N, R)
		mutated := mutateChromosome(pure_child, N, R, options)
		child_chromosome := chromosome {
			tour: mutated,
			cost: g.GetTourCost(mutated),
		}
		population.Put(float64(child_chromosome.cost), child_chromosome)

	}
	return true
}

func crossoverER(c1, c2 []int, N int, R *rand.Rand) []int {
	child := make([]int, N)

	edge_map := make(map[int]*hashset.Set, N)
	for i := 0; i < N; i++ {
		edge_map[i] = hashset.New()
	}

	V := make([]bool, N)

	for i := 0; i < N; i++ {
		v1, v2 := c1[i], c2[i]
		p1, p2 := c1[(N + i - 1) % N], c2[(N + i - 1) % N]
		n1, n2 := c1[(i + 1) % N], c2[(i + 1) % N]

		edge_map[v1].Add(p1, n1)
		edge_map[v2].Add(p2, n2)
	}

	current_city := R.Intn(N)
	for i := 0; i < N; i++ {
		child[i] = current_city
		V[current_city] = true

		_fringe := edge_map[current_city].Values()
		fringe := make([]int, len(_fringe))
		for j := 0; j < len(_fringe); j++ {
			fringe[j] = _fringe[j].(int)
		}
		current_city = getLowestConnectedVertex(edge_map, fringe, V, N, R)

		if current_city != -1 {
			continue
		}

		fringe = make([]int, 0, N)
		for j := 0; j < N; j++ {
			if !V[j] {
				fringe = append(fringe, j)
			} 
		}
		current_city = getLowestConnectedVertex(edge_map, fringe, V, N, R)

		if current_city != -1 {
			continue
		}
		break
	}

	return child
}

func getLowestConnectedVertex(
	edge_map map[int]*hashset.Set, candidates []int, V []bool, 
	N int, R *rand.Rand,
) int {
	min_connections := N
	best_candidates := make([]int, N)
	num_candidates := 0
	for _, v := range candidates {
		l := edge_map[v].Size()
		if l == 0 || V[v] {
			continue
		}
		if l == min_connections {
			best_candidates[num_candidates] = v
			num_candidates++
		} else if l < min_connections {
			min_connections = l
			best_candidates[0] = v
			num_candidates = 1
		}
	}
	if num_candidates == 0 {
		return -1
	}
	return best_candidates[R.Intn(num_candidates)]
}

func mutateChromosome(c []int, 
	N int, R *rand.Rand, options GeneticOptions,
) []int {
	if R.Float64() >= options.MutationChance {
		return c
	}
	l, r := R.Intn(N), R.Intn(N - 1)
	if l <= r {
		r++
	}
	if l > r {
		l, r = r, l
	}
	new := make([]int, N)
	copy(new, c)
	new[l] = c[r]
	new[r] = c[l]
	return new
}