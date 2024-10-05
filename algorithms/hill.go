package algorithms

import (
	"hw2/tsp"
	"math/rand"
)

type HillClimbingOptions struct {
	Restarts int 
}

func HillClimbingLocalSearch(g *tsp.Graph, N int, options HillClimbingOptions) (tsp.Weight, int, []int) {
	
	var best_tour []int = nil
	var best_cost tsp.Weight = -1

	var tour []int
	var value tsp.Weight
	var delta, best_delta tsp.Weight
	var temp_i, temp_j int

	for restart := 0; restart < options.Restarts; restart++ {
		tour = rand.Perm(N)
		value = g.GetTourCost(tour)
		best_delta = -1
		temp_i, temp_j = 1, 1

		for best_delta < 0 {
			SwapEdges(tour, temp_i, temp_j)
			best_delta = 0
			for i := 0; i < N - 1; i++ {
				for j := i + 1; j < N; j++ {
					delta = Calculate2OptDelta(g, N, tour, i, j)
					if delta < best_delta {
						temp_i, temp_j = i, j
						best_delta = delta
					}
				}
			}
			value += best_delta
		}

		if best_cost == -1 || value < best_cost {
			best_tour = tour
			best_cost = value
		}
	}
	
	return best_cost, N * (N - 1) * N, best_tour
}