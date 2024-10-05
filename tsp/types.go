package tsp

type Weight int

// Graph represents a set of vertices connected by edges.
type Edge struct {
	W Weight
	To *Vertex
}

type Vertex struct {
	I int
	E map[int]*Edge
}

type Graph struct {
	V map[int]*Vertex
	N int
}