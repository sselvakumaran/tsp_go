package main

import (
	"fmt"
	"hw2/algorithms"
	"hw2/tsp"
	"os"
	"syscall"
	"time"
)

type Algorithm string
const (
	NN = "NN"
	NN2O = "NN2O"
	RNN = "RNN"
	A_MST = "A_MST" 
	Bi_A_MST = "Bi_A_MST"
	hillClimbing = "hillClimbing" 
	simuAnnealing = "simuAnnealing" 
	genetic = "genetic"
)

func measure(fn func(*tsp.Graph, int) (tsp.Weight, int), g *tsp.Graph, N int) (
	tsp.Weight, int, float64, float64,
) {
	var usage1, usage2 *syscall.Rusage = new(syscall.Rusage), new(syscall.Rusage)

	start := time.Now()
	syscall.Getrusage(syscall.RUSAGE_SELF, usage1)

	cost, nodes := fn(g, N)

	elapsed := time.Since(start)
	syscall.Getrusage(syscall.RUSAGE_SELF, usage2)

	return cost, nodes, 
	float64(elapsed.Microseconds()) / 1000, 
	float64(CalcuateCPUTimeDifference(usage1, usage2).Microseconds()) / 1000
}

func main() {
	var start time.Time
	var elapsed time.Duration
	var usage1, usage2 *syscall.Rusage = new(syscall.Rusage), new(syscall.Rusage)

	var ALGORITHMS = [...]Algorithm{NN, NN2O, RNN, 
		A_MST, Bi_A_MST, 
		hillClimbing, simuAnnealing, genetic}

	FILENAMES := generateFileNames()

	if (len(os.Args) == 1) {
		g, N, err := tsp.ReadAdjacencyMatrix()

		if err != nil {
			panic(fmt.Sprintf("error reading: %s\n", err))
		}

		for _, v := range ALGORITHMS {
			run_algorithm := BenchmarkAlgorithm(v)

			start = time.Now()
			syscall.Getrusage(syscall.RUSAGE_SELF, usage1)

			cost, num_nodes := run_algorithm(g, N)
			
			syscall.Getrusage(syscall.RUSAGE_SELF, usage2)
			elapsed = time.Since(start)
			
			fmt.Printf("%-15s (%5d): %s / %s\n", v, 
				cost, 
				elapsed, 
				CalcuateCPUTimeDifference(usage1, usage2))
			output := fmt.Sprintf("cost,nodes,cpu,wall\n%d,%d,%.3f,%.3f", 
				cost, num_nodes, 
				float64(elapsed.Microseconds()) / 1000, 
				float64(CalcuateCPUTimeDifference(usage1, usage2).Microseconds()) / 1000)

			os.WriteFile(fmt.Sprintf("output/%s.csv", v), []byte(output), 0644)

		}
		return
	}
	if len(os.Args) == 2 {panic("wrong arguments")}

	if os.Args[1] == "pt1" && os.Args[2] == "all" {
		for _, alg := range []Algorithm{NN, NN2O, RNN} {

			run_alg := BenchmarkAlgorithm(alg)
			os.Mkdir("output/pt1_all", 0777)

			file, err := os.Create(fmt.Sprintf("output/pt1_all/%s.csv", alg))

			if err != nil {
        panic(err)
			}
			file.WriteString("N,I,cost,nodes,cpu,wall\n")

			fmt.Printf("%+v\n", FILENAMES)

			for size := 0; size < 6; size++ {
				V := FILENAMES[size]
				for id, name := range V {
					reader, _ := os.Open(name)
					g, N, err := tsp.ReadAdjacencyMatrixFromReader(reader)

					if err != nil {
						panic(fmt.Sprintf("error reading: %s\n", err))
					}

					reader.Close()
					
					cost, nodes, cpu, wall := measure(run_alg, g, N)
					
					line := fmt.Sprintf("%d,%d,%d,%d,%.3f,%.3f\n",
					5 * (size + 1), id + 1, cost, nodes, cpu, wall)
					fmt.Printf("%s\n", name)
					file.WriteString(line)
					
				}
			}
			file.Close()
		}
		return
	}
	if os.Args[1] == "pt1" && os.Args[2] == "rnn" {
		os.Mkdir("output/parameter_search", 0777)

		file, err := os.Create(fmt.Sprintf("output/parameter_search/%s.csv", RNN))
		if err != nil {
			panic(err)
		}

		file.WriteString("N,I,n,cost,nodes,cpu,wall\n")

		V := FILENAMES[5]

		for n := 1; n <= 10; n++ {
			for i := 0; i < 50; i++ {
				id := i
				name := V[0]
				fmt.Printf("%s\n", name)
				reader, _ := os.Open(name)
				g, N, err := tsp.ReadAdjacencyMatrixFromReader(reader)
				if err != nil {
					panic(fmt.Sprintf("error reading: %s\n", err))
				}
	
				run_alg := func(*tsp.Graph, int) (tsp.Weight, int) {
					return getBenchmarks(algorithms.RRNearestNeighbor2OptSearch(g, N, 
						algorithms.RNNOptions{
							Repeats: N,
							NearestNodeBreadth: n,
							NumThreads: 4,
					}))
				}
				reader.Close()
				cost, nodes, cpu, wall := measure(run_alg, g, N)
				line := fmt.Sprintf("%d,%d,%d,%d,%d,%.3f,%.3f\n",
				30, id + 1, n, cost, nodes, cpu, wall)
				file.WriteString(line)
			}
		}
		file.Close()
		return
	}
	if os.Args[1] == "pt2" && os.Args[2] == "all" {
		for _, alg := range []Algorithm{A_MST, Bi_A_MST} {

			run_alg := BenchmarkAlgorithm(alg)
			os.Mkdir("output/pt2_all", 0777)

			file, err := os.Create(fmt.Sprintf("output/pt2_all/%s.csv", alg))

			if err != nil {
        panic(err)
			}
			file.WriteString("N,I,cost,nodes,cpu,wall\n")

			for size := 0; size < 6; size++ {
				V := FILENAMES[size]
				for id, name := range V {
					reader, _ := os.Open(name)
					g, N, err := tsp.ReadAdjacencyMatrixFromReader(reader)

					if err != nil {
						panic(fmt.Sprintf("error reading: %s\n", err))
					}

					reader.Close()
					
					cost, nodes, cpu, wall := measure(run_alg, g, N)
					
					line := fmt.Sprintf("%d,%d,%d,%d,%.3f,%.3f\n",
					5 * (size + 1), id + 1, cost, nodes, cpu, wall)
					fmt.Printf("%s\n", name)
					file.WriteString(line)
					
				}
			}
			file.Close()
		}
	}
	if os.Args[1] == "pt3" && os.Args[2] == "all" {
		os.Mkdir("output/pt3_all", 0777)
		for _, alg := range []Algorithm{genetic} {

			run_alg := BenchmarkAlgorithm(alg)

			file, err := os.Create(fmt.Sprintf("output/pt3_all/%s.csv", alg))

			if err != nil {
        panic(err)
			}
			file.WriteString("N,I,cost,nodes,cpu,wall\n")

			for size := 0; size < 6; size++ {
				V := FILENAMES[size]
				for id, name := range V {
					fmt.Printf("%s\n", name)
					reader, _ := os.Open(name)
					g, N, err := tsp.ReadAdjacencyMatrixFromReader(reader)
					if err != nil {
						panic(fmt.Sprintf("error reading: %s\n", err))
					}

					reader.Close()
					cost, nodes, cpu, wall := measure(run_alg, g, N)
					line := fmt.Sprintf("%d,%d,%d,%d,%.3f,%.3f\n",
					5 * (size + 1), id + 1, cost, nodes, cpu, wall)
					file.WriteString(line)
					
				}
			}
			file.Close()
		}
		return
	}
	if os.Args[1] == "pt3" && os.Args[2] == "hill" {
		os.Mkdir("output/parameter_search", 0777)

		file, err := os.Create(fmt.Sprintf("output/parameter_search/%s.csv", RNN))
		if err != nil {
			panic(err)
		}

		file.WriteString("N,I,restarts,cost,nodes,cpu,wall\n")

		V := FILENAMES[5]

		for n := 10; n <= 100; n += 10 {
			for id, name := range V {
				fmt.Printf("%s\n", name)
				reader, _ := os.Open(name)
				g, N, err := tsp.ReadAdjacencyMatrixFromReader(reader)
				if err != nil {
					panic(fmt.Sprintf("error reading: %s\n", err))
				}
	
				run_alg := func(g *tsp.Graph, N int) (tsp.Weight, int) {
					return getBenchmarks(algorithms.HillClimbingLocalSearch(g, N, 
						algorithms.HillClimbingOptions{
							Restarts: n,
					}))
				}
				reader.Close()
				cost, nodes, cpu, wall := measure(run_alg, g, N)
				line := fmt.Sprintf("%d,%d,%d,%d,%d,%.3f,%.3f\n",
				30, id + 1, n, cost, nodes, cpu, wall)
				file.WriteString(line)
			}
		}
		file.Close()
		return
	}
	// if os.Args[1] == "pt3" && os.Args[2] == "simu" && os.Args[3] == "restarts" {
	// 	panic("didn't write yet")
	// }
	// if os.Args[1] == "pt3" && os.Args[2] == "simu" && os.Args[3] == "temp" {
	// 	panic("didn't write yet")
	// }
	// if os.Args[1] == "pt3" && os.Args[2] == "simu" && os.Args[3] == "alpha" {
	// 	panic("didn't write yet")
	// }
	// if os.Args[1] == "pt3" && os.Args[3] == "genetic" {
	// 	panic("didn't write yet")
	// }
	panic("no instructions available for input")
}

func generateFileNames() map[int][]string {
	A := make(map[int][]string, 6)
	for i, x := range []int{5,10,15,20,25,30} {
		A[i] = make([]string, 30)
		for j := 0; j < 30; j++ {
			A[i][j] = fmt.Sprintf("input/infile%02d_%02d.txt", x, j + 1)
		}
	}
	return A
}

func getBenchmarks(cost tsp.Weight, nodes int, _ []int) (tsp.Weight, int) {
	return cost, nodes
}

func BenchmarkAlgorithm(i Algorithm) (func (*tsp.Graph, int) (tsp.Weight, int)) {
	switch (i) {
	case NN:
		return func(g *tsp.Graph, N int) (tsp.Weight, int) {
			return getBenchmarks(algorithms.NearestNeighborSearch(g, N, 
				algorithms.NNOptions{
					StartingNode: 0,
				}))
		}
	case NN2O:
		return func(g *tsp.Graph, N int) (tsp.Weight, int) {
			return getBenchmarks(algorithms.NearestNeighbor2OptSearch(g, N, 
				algorithms.NN2OOptions{
					StartingNode: 0,
			}))
		}
	case RNN:
		return func(g *tsp.Graph, N int) (tsp.Weight, int) {
			return getBenchmarks(algorithms.RRNearestNeighbor2OptSearch(g, N, 
				algorithms.RNNOptions{
					Repeats: 10 * N,
					NearestNodeBreadth: 3,
					NumThreads: 4,
			}))
		}
	case A_MST:
		return func(g *tsp.Graph, N int) (tsp.Weight, int) {
			h := float32(1)
			if N > 25 {
				h = 1.05
			}
			return getBenchmarks(algorithms.AStarMSTSearch(g, N, 
				algorithms.AStarOptions{
					StartingNode: 0,
					HConstant: h,
					BeamWidth: min(20, N),
					FScoreCutoff: 100,
			}))
		}
	case Bi_A_MST:
		return func(g *tsp.Graph, N int) (tsp.Weight, int) {
			return getBenchmarks(algorithms.BidirectionalAStarMSTSearch(g, N, 
				algorithms.AStarOptions{
					StartingNode: 0,
					HConstant: 1,
					BeamWidth: N,
					FScoreCutoff: 100,
			}))
		}
	case hillClimbing:
		return func(g *tsp.Graph, N int) (tsp.Weight, int) {
			return getBenchmarks(algorithms.HillClimbingLocalSearch(g, N, 
				algorithms.HillClimbingOptions{
					Restarts: N * 5,
			}))
		}
	case simuAnnealing:
		return func(g *tsp.Graph, N int) (tsp.Weight, int) {
			return getBenchmarks(algorithms.SimulatedAnnealingLocalSearch(g, N, 
				algorithms.SimulatedAnnealingOptions{
					Temperature: 20.0, 
					Alpha: 0.95, 
					NumIterations: 1500,
					NumRestarts: 10,
			}))
		}
	case genetic:
		return func(g *tsp.Graph, N int) (tsp.Weight, int) {
			pop_size := 10 * N
			if N == 5 {
				pop_size = 5
			}
			return getBenchmarks(algorithms.GeneticLocalSearch(g, N, 
				algorithms.GeneticOptions{
					PopulationSize: pop_size, 
					NumGenerations: 400 + N * 20,
					RolloverCount: 0.4, 
					BTreeOrder: 3, 
					MutationChance: 0.05,
			}))
		}
	}
	return nil
}

func CalcuateCPUTimeDifference(u1, u2 *syscall.Rusage) time.Duration {
	return time.Duration(u2.Utime.Nano() + u2.Stime.Nano() - 
		u1.Utime.Nano() - u1.Stime.Nano())
}