package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"testing"
)

// creating a new test object for each function to test
type VerifyCliqueTest struct {
	graph    [][]int
	vertices []int
	result   bool
}

type EnumerateCliquesSerialTest struct {
	graph  [][]int
	k      int
	result [][]int
}

type GetSizeKSubgraphTest struct {
	k      int
	n      int
	result [][]int
}

type EnumerateCliquesParallelTest struct {
	graph  [][]int
	k      int
	result [][]int
}

/*
--------------------------------
Testing Functions
--------------------------------
*/

// TestVerifyClique tests the VerifyClique function getting the tests
// from the ReadVerifyCliqueTests function
func TestVerifyClique(t *testing.T) {
	//read in all tests from the directory and run them
	tests := ReadVerifyCliqueTests("testcases/VerifyClique")
	//run each test
	for _, test := range tests {
		res := VerifyClique(test.graph, test.vertices)
		if res != test.result {
			t.Errorf("VerifyClique(%v, %v) = %v; want %v", test.graph, test.vertices, res, test.result)
		}
	}
}

// ReadVerifyCliqueTests will read in the test cases for VerifyClique from a directory
// and return the test cases as a slice of VerifyCliqueTest objects
func ReadVerifyCliqueTests(directory string) []VerifyCliqueTest {
	//read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "/input/")
	numFiles := len(inputFiles)

	tests := make([]VerifyCliqueTest, numFiles)
	for i, inputFile := range inputFiles {
		//I'll actually use the same graph for all of these, called graph.txt, for simplicity
		tests[i].graph = ReadGraphFromFile("Tests/graph.txt")
		tests[i].vertices = ReadIntegersFromFile(directory + "/input/" + inputFile.Name())
	}

	//now, read output files
	outputFiles := ReadDirectory(directory + "/output/")
	if len(outputFiles) != numFiles {
		panic("Error: number of input and output files do not match!")
	}

	for i, outputFile := range outputFiles {
		//read in the test's result
		tests[i].result = ReadBooleanFromFile(directory + "/output/" + outputFile.Name())
	}

	return tests
}

// TestEnumerateCliquesSerial tests the EnumerateCliquesSerial function getting the tests
// from the ReadEnumerateCliquesSerialTests function
func TestEnumerateCliquesSerial(t *testing.T) {
	//read in all tests from the directory and run them
	tests := ReadEnumerateCliquesSerialTests("testcases/EnumerateCliquesSerial")
	//run each test
	for _, test := range tests {
		res := EnumerateCliquesSerial(test.graph, test.k)
		if !equal2D(res, test.result) {
			t.Errorf("EnumerateCliquesSerial(%v, %v) = %v; want %v", test.graph, test.k, res, test.result)
		}
	}
}

// ReadEnumerateCliquesSerialTests will read in the test cases for EnumerateCliquesSerial from a directory
// and return the test cases as a slice of EnumerateCliquesSerialTest objects
func ReadEnumerateCliquesSerialTests(directory string) []EnumerateCliquesSerialTest {
	//read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "/input/")
	numFiles := len(inputFiles)

	tests := make([]EnumerateCliquesSerialTest, numFiles)
	for i, inputFile := range inputFiles {
		//I'll actually use the same graph for all of these, called graph.txt, for simplicity
		tests[i].graph = ReadGraphFromFile("Tests/graph.txt")
		tests[i].k = ReadIntegersFromFile(directory + "/input/" + inputFile.Name())[0]
	}

	//now, read output files
	outputFiles := ReadDirectory(directory + "/output/")
	if len(outputFiles) != numFiles {
		panic("Error: number of input and output files do not match!")
	}

	for i, outputFile := range outputFiles {
		//read in the test's result
		tests[i].result = ReadGraphFromFile(directory + "/output/" + outputFile.Name())
	}

	return tests
}

// TestGetSizeKSubgraph tests the GetSizeKSubgraph function getting the tests
// from the ReadGetSizeKSubgraphTests function
func TestGetSizeKSubgraph(t *testing.T) {
	//read in all tests from the directory and run them
	tests := ReadGetSizeKSubgraphTests("testcases/GetSizeKSubgraph")
	//run each test
	for _, test := range tests {
		res := GetSizeKSubgraph(test.k, test.n)
		if !equal2D(res, test.result) {
			t.Errorf("GetSizeKSubgraph(%v, %v) = %v; want %v", test.k, test.n, res, test.result)
		}
	}
}

// ReadGetSizeKSubgraphTests will read in the test cases for GetSizeKSubgraph from a directory
// and return the test cases as a slice of GetSizeKSubgraphTest objects
func ReadGetSizeKSubgraphTests(directory string) []GetSizeKSubgraphTest {
	//read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "/input/")
	numFiles := len(inputFiles)

	tests := make([]GetSizeKSubgraphTest, numFiles)
	for i, inputFile := range inputFiles {
		tests[i].k = ReadIntegersFromFile(directory + "/input/" + inputFile.Name())[0]
		tests[i].n = ReadIntegersFromFile(directory + "/input/" + inputFile.Name())[1]
	}

	//now, read output files
	outputFiles := ReadDirectory(directory + "/output/")
	if len(outputFiles) != numFiles {
		panic("Error: number of input and output files do not match!")
	}

	for i, outputFile := range outputFiles {
		//read in the test's result
		tests[i].result = ReadGraphFromFile(directory + "/output/" + outputFile.Name())
	}

	return tests
}

// TestEnumerateCliquesParallel tests the EnumerateCliquesParallel function getting the tests
// from the ReadEnumerateCliquesParallelTests function
func TestEnumerateCliquesParallel(t *testing.T) {
	//read in all tests from the directory and run them
	tests := ReadEnumerateCliquesParallelTests("testcases/EnumerateCliquesParallel")
	//run each test
	for _, test := range tests {
		res := EnumerateCliquesParallel(test.graph, test.k, 1)
		if !equal2D(res, test.result) {
			t.Errorf("EnumerateCliquesParallel(%v, %v) = %v; want %v", test.graph, test.k, res, test.result)
		}
	}
}

// ReadEnumerateCliquesParallelTests will read in the test cases for EnumerateCliquesParallel from a directory
// and return the test cases as a slice of EnumerateCliquesParallelTest objects
func ReadEnumerateCliquesParallelTests(directory string) []EnumerateCliquesParallelTest {
	//read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "/input/")
	numFiles := len(inputFiles)

	tests := make([]EnumerateCliquesParallelTest, numFiles)
	for i, inputFile := range inputFiles {
		//I'll actually use the same graph for all of these, called graph.txt, for simplicity
		tests[i].graph = ReadGraphFromFile("Tests/graph.txt")
		tests[i].k = ReadIntegersFromFile(directory + "/input/" + inputFile.Name())[0]
	}

	//now, read output files
	outputFiles := ReadDirectory(directory + "/output/")
	if len(outputFiles) != numFiles {
		panic("Error: number of input and output files do not match!")
	}

	for i, outputFile := range outputFiles {
		//read in the test's result
		tests[i].result = ReadGraphFromFile(directory + "/output/" + outputFile.Name())
	}

	return tests
}

/*
--------------------------------
I/O Functions
--------------------------------
*/

// ReadDirectory reads in a directory and returns a slice of fs.DirEntry objects containing file info for the directory
func ReadDirectory(dir string) []fs.DirEntry {
	//read in all files in the given directory
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	return files
}

/*
First, I want to make a function ReadGraphFromFile() that takes in a filename containing the
entries of a 2D integer array representing a graph adjacency matrix, and returns
the adjacency matrix as a 2D integer array.
*/
func ReadGraphFromFile(filename string) [][]int {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	// Read the file line by line
	var graph [][]int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Split the line by spaces
		line := scanner.Text()
		entries := strings.Fields(line)
		// Convert the entries to integers
		var row []int
		for _, entry := range entries {
			num, err := strconv.Atoi(entry)
			if err != nil {
				panic(err)
			}
			row = append(row, num)
		}
		// Append the row to the graph
		graph = append(graph, row)
	}

	return graph
}

/*
next we need a function ReadIntegersFromFile() that takes in a file name
and returns the integers in the file in a 1D array.
*/
func ReadIntegersFromFile(filename string) []int {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	// Read the file line by line
	var integers []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Split the line by spaces
		line := scanner.Text()
		entries := strings.Fields(line)
		// Convert the entries to integers
		for _, entry := range entries {
			num, err := strconv.Atoi(entry)
			if err != nil {
				panic(err)
			}
			integers = append(integers, num)
		}
	}

	return integers
}

/*
	finally we need a function ReadBooleanFromFile() that takes in a file name

and returns the boolean value in the file.
*/
func ReadBooleanFromFile(filename string) bool {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Read the boolean value
		line := scanner.Text()
		boolean, err := strconv.ParseBool(line)
		if err != nil {
			panic(err)
		}
		return boolean
	}

	return false
}

/*
TestGraphsEqual() takes in two graphs as 2D integer arrays and returns a boolean value
as to whether or not each value at each position in the graph is the same or not.
*/
func equal2D(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}
