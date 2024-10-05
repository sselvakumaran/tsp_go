package algorithms

import (
	"hw2/tsp"
)

type NN2OOptions struct {
	StartingNode int
}

func NearestNeighbor2OptSearch(g *tsp.Graph, N int, options NN2OOptions) (tsp.Weight, int, []int) {
	cost, nodes_expanded, tour := NearestNeighborSearch(g, N, NNOptions(options))
	continue_opt := true
	var delta tsp.Weight = 0
	for continue_opt {
		continue_opt = false
		for i := 0; i < N && !continue_opt; i++ {
			for j := i + 1; j < N; j++ {
				delta = Calculate2OptDelta(g, N, tour, i, j)
				if delta < 0 {
					continue_opt = true
					SwapEdges(tour, i, j)
					cost += delta
					nodes_expanded++
					break
				}
			}
		}
	}
	return cost, nodes_expanded, tour
}