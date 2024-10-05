package algorithms

import (
	"hw2/tsp"
	"strings"

	"github.com/emirpasic/gods/queues/priorityqueue"
)

// returns the inverse
func getInverseEncoding(slice []bool, middle int8, start int) string {
	if len(slice) == 0 {
		return ""
	}
	middle32 := int(middle)
	// Preallocate a strings.Builder for efficient string concatenation
	var builder strings.Builder
	builder.Grow(len(slice)) // Grow to the expected size to minimize reallocations

	for i, b := range slice {
		if middle32 != i && (b || i == start) {
			builder.WriteByte('0') // Write '1' for true
		} else {
			builder.WriteByte('1') // Write '0' for false
		}
	}

	return builder.String()
}

func combineTours(slice1 []int8, slice2 []int8, start int) []int {
	L1, L2 := len(slice1), len(slice2)
	var result []int = make([]int, L1 + L2)
	for i, v := range slice1 {
		result[i] = int(v)
	}
	for i, v := range slice2 {
		result[L1 + L2 - i - 2] = int(v)
	}
	result[L1 + L2 - 1] = start
	return result
}
func BidirectionalAStarMSTSearch(g *tsp.Graph, N int, options AStarOptions) (tsp.Weight, int, []int) {
	Q := priorityqueue.NewWith(aStarStateComparator)
	var h_cache map[string]tsp.Weight = make(map[string]tsp.Weight, 1 << int(N / 2))
	var g_scores map[string]tsp.Weight = make(map[string]tsp.Weight, 1 << int(N / 2))
	var path_cache map[string]*[]int8 = make(map[string]*[]int8, 1 << int(N / 2))

	best_f := tsp.Weight(float32(MSTHeuristic(g, N, make([]bool, N), N - 1, options.StartingNode)) * options.HConstant)

	for i := 0; i < N; i++ {
		if i == options.StartingNode {continue}

		V := make([]bool, N)
		P := make([]int8, 1)
		P[0] = int8(i)

		h := MSTHeuristic(g, N, V, N - 1, options.StartingNode)

		V[i] = true

		weight := g.V[options.StartingNode].E[i].W
		encoding := encodeV(V)
		h_cache[encoding] = h
		g_scores[encoding] = weight
		path_cache[encoding] = &P
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
	best_d := -1
	for !Q.Empty() {
		uncasted_current, ok := Q.Dequeue()
		if !ok {break}
		current := uncasted_current.(aStarState)
		best_f = min(best_f, current.FScore)
		nodes_expanded++

		last := current.Path[current.Count - 1]
		inv_encoding := getInverseEncoding(current.Visited, last, options.StartingNode)
		if compliment, ok := g_scores[inv_encoding]; ok {
			p_inv := path_cache[inv_encoding]
			if (*p_inv)[len(*p_inv) - 1] == last {
				return current.Cost + compliment, nodes_expanded,
					combineTours(current.Path, *path_cache[inv_encoding], options.StartingNode)
			}
		}

		if current.Count > best_d {
			best_d = current.Count
		}
		if current.Count > (N / 2) + 1 {
			continue
		}
		
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
						h_cache[encoding] = MSTHeuristic(g, N, current.Visited, N - current.Count, options.StartingNode)
					}
					path_cache[encoding] = &new_p

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
				if node.(aStarState).FScore > 
					tsp.Weight(float32(best_f) * options.FScoreCutoff) {
						break
				}
				Q.Enqueue(node)
			}
		}
		beam.Clear()
	}

	return -1, nodes_expanded, nil
}