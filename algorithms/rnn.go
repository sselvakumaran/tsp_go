package algorithms

import (
	"hw2/tsp"
	"math/rand"
	"sort"
	"sync"

	"github.com/emirpasic/gods/sets/hashset"
)

type RNNOptions struct {
	Repeats int
	NearestNodeBreadth int
	NumThreads int
}

type trial struct {
	cost tsp.Weight
	tour []int
}

func findRandomNearestNeighbor(g *tsp.Graph, from int, V *hashset.Set, n int) *tsp.Edge {
	edges := g.GetNeighboringUnvisitedEdges(from, func(v int) bool {
		return !V.Contains(v)
	})
	num_options := len(edges)
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].W < edges[j].W
	})
	max_pick := min(num_options, n)
	if max_pick == 0 {
		return nil
	}
	return edges[rand.Intn(max_pick)]
}

func RRNearestNeighbor2OptSearch(g *tsp.Graph, N int, options RNNOptions) (tsp.Weight, int, []int) {
	edge_selector := func(g *tsp.Graph, from int, V *hashset.Set) *tsp.Edge {
		return findRandomNearestNeighbor(g, from, V, options.NearestNodeBreadth)
	}
	
	trials := make([]trial, options.Repeats)
	
	var wait_group sync.WaitGroup

	for i := 0; i < options.NumThreads; i++ {
		wait_group.Add(1)
		go func(threadIndex int) {
			defer wait_group.Done()

			for k := threadIndex; k < options.Repeats; k += options.NumThreads {
				x := rand.Intn(N)
				t_cost, _, t_tour := GreedyDepthSearch(g, N, edge_selector, x)
				trials[k] = trial{
					cost: t_cost, 
					tour: t_tour,
				}
			}
		}(i)
	}

	wait_group.Wait()

	var best_trial int = 0
	for i := 1; i < options.Repeats; i++ {
		if trials[i].cost < trials[best_trial].cost {
			best_trial = i
		}
	}
	tour := trials[best_trial].tour
	cost := trials[best_trial].cost
	continue_opt := true
	var delta tsp.Weight
	nodes_expanded := options.Repeats * N

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