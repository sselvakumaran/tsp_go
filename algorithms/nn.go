package algorithms

import (
	"hw2/tsp"

	"github.com/emirpasic/gods/sets/hashset"
)

type NNOptions struct {
	StartingNode int
}

func findNearestUnvisitedNeighbor(g *tsp.Graph, from int, V *hashset.Set) *tsp.Edge {
	var low *tsp.Edge = nil
	for k, v := range g.V[from].E {
		if (low == nil || v.W < low.W) && !V.Contains(k) {
			low = v
		}
	}
	return low
}

func NearestNeighborSearch(g *tsp.Graph, N int, options NNOptions) (tsp.Weight, int, []int) {
	return GreedyDepthSearch(g, N, findNearestUnvisitedNeighbor, options.StartingNode)
}
