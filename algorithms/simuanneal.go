package algorithms

import (
	"hw2/tsp"
	"math"
	"math/rand"
)

type SimulatedAnnealingOptions struct {
	Temperature float64
	Alpha float64
	NumIterations int
	NumRestarts int
}

func probability(p float64) bool {
	return rand.Float64() <= p 
}

func SimulatedAnnealingLocalSearch(g *tsp.Graph, N int, options SimulatedAnnealingOptions) (tsp.Weight, int, []int) {
	var T float64 = options.Temperature
	var tour, best_tour []int = nil, nil
	var cost, best_cost tsp.Weight = -1, -1
	var E float64

	var delta tsp.Weight
	var e_i float64

	for repeat := 0; repeat < options.NumRestarts; repeat++ {
		tour = rand.Perm(N)
		cost = g.GetTourCost(tour)
		E = 1.0 / float64(cost)
		
		for i := 0; i < options.NumIterations || T > 1e-8; i++ {
			l, r := rand.Intn(N), rand.Intn(N - 1)
			if l <= r {
				r++
			}
			if l > r {
				l, r = r, l
			}
			delta = Calculate2OptDelta(g, N, tour, l, r)
			if delta < 0 {
				SwapEdges(tour, l, r)
				cost += delta
				E = 1.0 / float64(cost)
			} else {
				e_i = 1.0 / float64(cost + delta)
				p := math.Exp(-(E - e_i) / T)
				if probability(p) {
					SwapEdges(tour, l, r)
					cost += delta
					E = 1.0 / float64(cost)
				}
			}
			T *= options.Alpha
		}
		if best_cost == -1 || cost < best_cost {
			best_cost = cost
			best_tour = tour
		}
	}
	
	return best_cost, options.NumRestarts * options.NumIterations, best_tour
}