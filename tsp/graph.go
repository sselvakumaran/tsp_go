package tsp

func NewGraph(num_vertices int) *Graph {
	return &Graph{
			V: make(map[int]*Vertex),
			N: num_vertices,
	}
}

func (g *Graph) AddVertex(index int) {
	if _, exists := g.V[index]; !exists {
			g.V[index] = &Vertex{
					I: index,
					E: make(map[int]*Edge),
			}
	}
}

func (g *Graph) AddBidirectionalEdge(from, to int, weight Weight) {

	// Add edge from 'from' to 'to'
	g.V[from].E[to] = &Edge{W: weight, To: g.V[to]}

	// Add edge from 'to' to 'from' (since the graph is undirected)
	g.V[to].E[from] = &Edge{W: weight, To: g.V[from]}
}

func (g *Graph) GetEdgeWeight(from, to int) Weight {
	return g.V[from].E[to].W
}

func (g *Graph) GetNeighboringEdges(v int) []*Edge {
	keys := make([]*Edge, len(g.V[v].E))
	i := 0
	for _, v := range g.V[v].E {
		keys[i] = v
		i++
	}
	return keys
}

func (g *Graph) GetNeighboringUnvisitedEdges(v int, unvisited func(v int) bool) []*Edge {
	keys := make([]*Edge, len(g.V[v].E))
	i := 0
	for k, v := range g.V[v].E {
		if unvisited(k) {
			keys[i] = v
			i++
		}
	}
	return keys[:i]
}

func (g *Graph) GetTourCost(slice []int) Weight {
	if len(slice) != g.N {
		return -1
	}
	var cost Weight = 0
	for i := 1; i < g.N; i++ {
		if e, ok := g.V[slice[i-1]].E[slice[i]]; ok {
			cost += e.W
		} else {
			return -1
		}
	}
	if e, ok := g.V[slice[g.N-1]].E[slice[0]]; ok {
		cost += e.W
	} else {
		return -1
	}
	return cost
}