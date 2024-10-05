package algorithms

import (
	"hw2/tsp"
	"strings"

	"github.com/emirpasic/gods/queues/priorityqueue"
)

/*
	1. TSP IS REVERSIBLE! Lookup opposite set and see if its in there
	2. store path differently
*/

type AStarOptions struct {
	StartingNode int
	HConstant float32
	BeamWidth int
	FScoreCutoff float32
}

type aStarState struct {
	FScore tsp.Weight
	Cost tsp.Weight
	Count int
	Visited []bool
	Path []int8
}

func MSTHeuristic(g *tsp.Graph, N int, avoid []bool, remaining, start int) tsp.Weight {
	fringe := priorityqueue.NewWith(func(a, b interface{}) int {
		return WeightComparator(a.(*tsp.Edge).W, b.(*tsp.Edge).W)
	})

	V := make([]bool, N)
	copy(V, avoid)
	for v, e := range g.V[start].E {
		if !V[v] {
			fringe.Enqueue(e)
		}
	}
	V[start] = true
	var cost tsp.Weight = 0
	for remaining > 0 {
		e_temp, ok := fringe.Dequeue()
		if !ok {break}
		edge := e_temp.(*tsp.Edge)
		i := edge.To.I
		if V[i] {continue}
		cost += edge.W
		V[i] = true
		remaining--
		for v, e := range g.V[i].E {
			if !V[v] {
				fringe.Enqueue(e)
			}
		}
	}
	return cost
}

func encodeV(slice []bool) string {
	if len(slice) == 0 {
		return ""
	}
	
	// Preallocate a strings.Builder for efficient string concatenation
	var builder strings.Builder
	builder.Grow(len(slice)) // Grow to the expected size to minimize reallocations

	for _, b := range slice {
		if b {
			builder.WriteByte('1') // Write '1' for true
		} else {
			builder.WriteByte('0') // Write '0' for false
		}
	}

	return builder.String()
}

func int8SliceToIntSlice(slice []int8) []int {
	var result []int = make([]int, len(slice))
	for i, v := range slice {
		result[i] = int(v)
	}
	return result
}

func aStarStateComparator(a, b interface{}) int {
	var a_score, b_score tsp.Weight = a.(aStarState).FScore, b.(aStarState).FScore
	switch {
	case a_score > b_score:
		return 1
	case a_score < b_score:
		return -1
	default:
		return 0
	}
}

func AStarMSTSearch(g *tsp.Graph, N int, options AStarOptions) (tsp.Weight, int, []int) {
	Q := priorityqueue.NewWith(aStarStateComparator)
	var h_cache map[string]tsp.Weight = make(map[string]tsp.Weight, 1 << int(N / 2))
	var g_scores map[string]tsp.Weight = make(map[string]tsp.Weight, 1 << int(N / 2))

	best_f := tsp.Weight(float32(MSTHeuristic(g, N, make([]bool, N), N, options.StartingNode)) * options.HConstant)

	for i := 0; i < N; i++ {
		if i == options.StartingNode {continue}

		V := make([]bool, N)
		P := make([]int8, 1)
		P[0] = int8(i)

		h := MSTHeuristic(g, N, V, N, options.StartingNode)

		V[i] = true

		weight := g.V[options.StartingNode].E[i].W
		encoding := encodeV(V)
		h_cache[encoding] = h
		g_scores[encoding] = weight
		
		Q.Enqueue(aStarState{
			FScore: weight + tsp.Weight(options.HConstant * float32(h)),
			Cost: weight,
			Count: 1,
			Visited: V,
			Path: P,
		})
	}
	
	beam := priorityqueue.NewWith(aStarStateComparator)
	nodes_expanded := 0
	for !Q.Empty() {
		uncasted_current, ok := Q.Dequeue()
		if !ok {break}
		current := uncasted_current.(aStarState)
		if current.Count == N {
			return current.Cost, nodes_expanded, int8SliceToIntSlice(current.Path)
		}
		best_f = min(best_f, current.FScore)
		nodes_expanded += 1
		
		for v, e := range g.V[int(current.Path[current.Count - 1])].E {
			if v == options.StartingNode && current.Count != N - 1 {continue}
			if !current.Visited[v] {
				new_cost := current.Cost + e.W

				new_v := make([]bool, N)
				copy(new_v, current.Visited)
				new_v[v] = true
				encoding := encodeV(new_v)
				
				if cost, ok := g_scores[encoding]; !ok || new_cost < cost {

					new_p := make([]int8, current.Count + 1)
					copy(new_p, current.Path)
					new_p[current.Count] = int8(v)

					g_scores[encoding] = new_cost
					if _, ok := h_cache[encoding]; !ok {
						h_cache[encoding] = MSTHeuristic(g, N, current.Visited, N - current.Count + 1, options.StartingNode)
					}

					new_state := aStarState{
						FScore: new_cost + tsp.Weight(options.HConstant * float32(h_cache[encoding])),
						Cost: new_cost,
						Count: current.Count + 1,
						Visited: new_v,
						Path: new_p,
					}
					beam.Enqueue(new_state)
				}
			}
		}
		
		for k := 0; k < options.BeamWidth && !beam.Empty(); k++ {
			if node, ok := beam.Dequeue(); ok {
				// if node.(aStarState).FScore > 
				// 	tsp.Weight(float32(best_f) * options.FScoreCutoff) {
				// 		break
				// }
				Q.Enqueue(node)
			}
		}
		beam.Clear()
	}

	return -1, nodes_expanded, nil
}