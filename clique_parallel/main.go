package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	//read in the params from the command line
	//first one will be the mode
	mode := os.Args[1]
	//second one will be the graph file
	graphFile := os.Args[2]
	//the third one is conditional
	//if mode is "verify" then the third one will be the vertices
	//if mode is "serial" then the third one will be the k
	//if mode is "parallel" then the third one will be the k
	//if mode is "parallel" then the fourth one will be the p
	var vertices []int
	var k int
	var p int
	if mode == "verify" {
		//get vertices from file
		vertFile := os.Args[3]
		vertices = ReadIntegersFromFile(vertFile)
	}
	if mode == "serial" || mode == "parallel" {
		//get k from command line
		k, _ = strconv.Atoi(os.Args[3])
		p, _ = strconv.Atoi(os.Args[4])
	}

	//read in the graph from the file
	graph := ReadGraphFromFile(graphFile)
	//run the appropriate function
	if mode == "verify" {
		res := VerifyClique(graph, vertices)
		if res {
			fmt.Println("The vertices form a clique.")
		} else {
			fmt.Println("The vertices do not form a clique.")
		}
	} else if mode == "serial" {
		cliques := EnumerateCliquesSerial(graph, k)
		fmt.Println("Cliques of size", k, "are:")
		for _, clique := range cliques {
			fmt.Println(clique)
		}
	} else if mode == "parallel" {
		cliques := EnumerateCliquesParallel(graph, k, p)
		fmt.Println("Cliques of size", k, "are:")
		for _, clique := range cliques {
			fmt.Println(clique)
		}
	} else {
		fmt.Println("Invalid mode. Please use 'verify', 'serial', or 'parallel'.")
	}
}

/*
VerifyClique() is a function that takes in a binary 2D array as an adjacency
matrix for a graph. It also takes in a list of vertices. It returns a boolean
value indicating whether the list of vertices forms a clique in the graph
(i.e. whether each vertex is connected to every other vertex in the list).
*/
func VerifyClique(graph [][]int, vertices []int) bool {
	// Iterate through all pairs of vertices
	for i := 0; i < len(vertices); i++ {
		for j := i + 1; j < len(vertices); j++ {
			// If there is no edge between the pair of vertices, return false
			if graph[vertices[i]][vertices[j]] == 0 {
				return false
			}
		}
	}
	// If all pairs of vertices are connected, return true
	return true
}

/*
Serial enumeration of cliques
EnumerateCliquesSerial() is a function that takes in a binary 2D array as an
adjacency matrix for a graph as well as an integer k. It returns a list of all
cliques of size k in the graph. The function uses a serial enumeration approach.
*/
func EnumerateCliquesSerial(graph [][]int, k int) [][]int {
	// Initialize an empty list to store the cliques
	cliques := [][]int{}

	// Iterate through all possible combinations of k vertices
	// first get all combinations
	subgraphs := GetSizeKSubgraph(k, len(graph))
	// then iterate through each combination
	for _, subgraph := range subgraphs {
		// If the combination forms a clique, add it to the list of cliques
		if VerifyClique(graph, subgraph) {
			cliques = append(cliques, subgraph)
		}
	}

	// Return the list of cliques
	return cliques
}

/*
Helper function for above: GetSizeKSubgraph()
GetSizeKSubgraph() is a helper function that takes in an integer k as well as
an integer n. It returns a list of all possible combinations of k vertices out
of n vertices; this is basically a 2D array.
*/
func GetSizeKSubgraph(k int, n int) [][]int {
	// We are going to define this recursively
	if k == 1 {
		// base case: return a list of lists containing each vertex
		subgraph := [][]int{}
		for i := 0; i < n; i++ {
			subgraph = append(subgraph, []int{i})
		}
		return subgraph
	} else {
		// recursive case: get all combinations of k-1 vertices out of n-1 vertices
		subgraph := [][]int{}
		subgraphPrev := GetSizeKSubgraph(k-1, n-1)
		for _, subgraphEntry := range subgraphPrev {
			for i := subgraphEntry[len(subgraphEntry)-1] + 1; i < n; i++ {
				subgraph = append(subgraph, append(subgraphEntry, i))
			}
		}
		return subgraph
	}
}

/*
Parallel enumeration of cliques
EnumerateCliquesParallel() is a function that takes in a binary 2D array as an
adjacency matrix for a graph as well as an integer k. It also takes in an integer p,
the number of processors to use. It returns a list of all
cliques of size k in the graph. The function uses a parallel enumeration approach.
*/
func EnumerateCliquesParallel(graph [][]int, k int, p int) [][]int {
	// Initialize an empty list to store the cliques
	cliques := [][]int{}

	// Create a channel to store the results
	results := make(chan []int)

	// Get all possible combinations of k vertices
	subgraphs := GetSizeKSubgraph(k, len(graph))
	numsubgraphs := len(subgraphs)

	//loop through p, and create p goroutines
	for i := 0; i < p; i++ {
		//the subgraphs should go from i to len(subgraphs) with a step of p
		//so we index each by [numsubgraphs * i / p : numsubgraphs * (i + 1) / p]
		currentSubgraphs := subgraphs[numsubgraphs*i/p : numsubgraphs*(i+1)/p]
		go VerifyCliquesOneProcessor(graph, currentSubgraphs, results)
	}

	// Collect the results from the channel
	for i := 0; i < numsubgraphs; i++ {
		clique := <-results
		cliques = append(cliques, clique)
	}

	// Return the list of cliques
	return cliques
}

// Helper function for above: VerifyCliquesOneProcessor()
/* VerifyCliquesOneProcessor() is a function that takes in a binary 2D array as an
adjacency matrix for a graph as well as a set of k-sized cliques and a channel.
It verifies if the cliques are valid and sends the valid cliques to the channel.
*/
func VerifyCliquesOneProcessor(graph [][]int, cliques [][]int, results chan []int) {
	for _, clique := range cliques {
		if VerifyClique(graph, clique) {
			results <- clique
		}
	}
}
