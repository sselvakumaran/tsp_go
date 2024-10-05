package algorithms

import (
	"fmt"
	"hw2/tsp"
	"os"
	"testing"
)

// TestMSTHeuristic tests the MSTHeuristic function
func TestMSTHeuristic1(t *testing.T) {
	// Read the graph from an adjacency matrix
	reader, _ := os.Open("../input/infile05_01.txt")
	g, N, err := tsp.ReadAdjacencyMatrixFromReader(reader)
	if err != nil {
		t.Fatalf("Error reading adjacency matrix: %v", err)
	}

	// Define the start node and avoidance array
	startNode := 0
	avoid := make([]bool, N)
	// avoid[1] = true // Example: Avoid node 1 for testing purposes

	// Calculate MST Heuristic for the given graph, avoiding node 1
	// remaining := N - 2 // Example: Remaining nodes to connect to the MST
	mstCost := MSTHeuristic(g, N, avoid, 5, startNode)

	// Expected cost (you can adjust this based on what you expect the cost to be)
	expectedCost := tsp.Weight(1062) // Example: Replace with the correct value

	// Compare the calculated MST cost with the expected cost
	if mstCost != expectedCost {
		t.Errorf("MSTHeuristic returned %v, but expected %v", mstCost, expectedCost)
	}
}

func TestMSTHeuristic2(t *testing.T) {
	// Read the graph from an adjacency matrix
	reader, _ := os.Open("../input/infile05_01.txt")
	g, N, err := tsp.ReadAdjacencyMatrixFromReader(reader)
	if err != nil {
		t.Fatalf("Error reading adjacency matrix: %v", err)
	}

	// Define the start node and avoidance array
	startNode := 0
	avoid := make([]bool, N)
	avoid[1] = true
	// avoid[1] = true // Example: Avoid node 1 for testing purposes

	// Calculate MST Heuristic for the given graph, avoiding node 1
	// remaining := N - 2 // Example: Remaining nodes to connect to the MST
	mstCost := MSTHeuristic(g, N, avoid, 3, startNode)

	// Expected cost (you can adjust this based on what you expect the cost to be)
	expectedCost := tsp.Weight(1062 - 236) // Example: Replace with the correct value

	// Compare the calculated MST cost with the expected cost
	if mstCost != expectedCost {
		t.Errorf("MSTHeuristic returned %v, but expected %v", mstCost, expectedCost)
	}
}

func TestAStar(t *testing.T) {
	reader, _ := os.Open("../input/infile05_01.txt")
	g, N, err := tsp.ReadAdjacencyMatrixFromReader(reader)
	if err != nil {
		t.Fatalf("Error reading adjacency matrix: %v", err)
	}

	cost, _, tour := AStarMSTSearch(g, N, AStarOptions{
		StartingNode: 0,
		HConstant: 1,
		BeamWidth: 5,
		FScoreCutoff: 30,
	})
	fmt.Printf("%d, %v\n", cost, tour)

	if cost != 1582 {
		t.Fail()
	}

}