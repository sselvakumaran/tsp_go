package algorithms

import (
	"hw2/tsp"

	"github.com/emirpasic/gods/sets/hashset"
)

type EdgeSelector func(*tsp.Graph, int, *hashset.Set) *tsp.Edge

func GreedyDepthSearch(g *tsp.Graph, N int, fn EdgeSelector, start int) (tsp.Weight, int, []int) {
	V := hashset.New()
	var current_vertex int = start
	var cost tsp.Weight = 0

	slice := make([]int, N)
	slice[0] = start

	V.Add(start)
	for i := 1; i < N; i++ {
		best_edge := fn(g, current_vertex, V)

		if best_edge == nil {
			break
		}
		
		V.Add(best_edge.To.I)
		current_vertex = best_edge.To.I
		slice[i] = best_edge.To.I

		cost += best_edge.W
	}
	cost += g.V[current_vertex].E[start].W

	return cost, N, slice
}

func Calculate2OptDelta(g *tsp.Graph, N int, route []int, v1, v2 int) tsp.Weight {
	return -g.GetEdgeWeight(route[v1], route[(v1 + 1) % N]) -
		g.GetEdgeWeight(route[v2], route[(v2 + 1) % N]) +
		g.GetEdgeWeight(route[(v1+1) % N], route[(v2 + 1) % N]) + 
		g.GetEdgeWeight(route[v1], route[v2])
}

func SwapEdges(route []int, v1, v2 int) {
	v1++
	for v1 < v2 {
		temp := route[v1]
		route[v1] = route[v2]
		route[v2] = temp
		v1++
		v2--
	}
}

func WeightComparator(a, b interface{}) int {
	var a_score, b_score tsp.Weight = a.(tsp.Weight), b.(tsp.Weight)
	switch {
	case a_score > b_score:
		return 1
	case a_score < b_score:
		return -1
	default:
		return 0
	}
}