package tsp

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func ReadAdjacencyMatrix() (*Graph, int, error) {
	return ReadAdjacencyMatrixFromReader(os.Stdin)
}

func ReadAdjacencyMatrixFromReader(r io.Reader) (*Graph, int, error) {
	scanner := bufio.NewScanner(r);

	if !scanner.Scan() {
		return nil, 0, fmt.Errorf("cannot read input")
	}
	numVertices, err := strconv.Atoi(scanner.Text())
	if err != nil || numVertices <= 0 {
		return nil, 0, fmt.Errorf("invalid number of vertices %v", err)
	}

	graph := NewGraph(numVertices)

	for i:= 0; i < numVertices; i++ {
		graph.AddVertex(i)
	}

	for i := 0; i < numVertices; i++ {
		if !scanner.Scan() {
			return nil, 0, fmt.Errorf("failed to read row")
		}

		row := strings.Fields(scanner.Text())

		if len(row) != numVertices {
			return nil, 0, fmt.Errorf("invalid number of vertices in row")
		}

		for j := 0; j < numVertices; j++ {
			weight, err := strconv.Atoi(row[j])
			if err != nil {
				return nil, 0, fmt.Errorf("invalid weight")
			}

			if weight != 0 {
				graph.AddBidirectionalEdge(i, j, Weight(weight))
			}
		}
	}
	return graph, numVertices, nil
}