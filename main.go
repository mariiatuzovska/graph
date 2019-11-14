package main

import (
	"encoding/csv"
	"log"
	"os"
)

// Data structure for quick getting data
type Data struct {
	data map[string][]string
}

// Road structure 
type Road stuct {
	head 	string
	road 	map[string]int
	visited map[string]bool
}

// type graph struct {
// 	reachabMatrix [][]bool
// 	vertex        map[string]int
// 	key           []string
// 	count         int
// }

// func initialise(path string) *Data {

// 	data := new(Data)
// 	data.getDataCSV(path)

// 	return data
// }

func main() {

	g := initialise("G1.csv")

}

func newRoad()

// func newGraph(data map[string][]string) *graph {

// 	g := new(graph)
// 	g.count = 0

// 	for key := range data {
// 		g.key[g.count] = key
// 		g.vertex[key] = g.count
// 		g.count++
// 	}

// 	g.reachabMatrix = make([][]bool, g.count)
// 	for i := range g.reachabMatrix {
// 		g.reachabMatrix[i] = make([]bool, g.count)
// 		for j := range g.reachabMatrix[i] {
// 			g.reachabMatrix[i][j] = false
// 		}
// 	}

// 	for i, vertex := range g.key {

// 		node := data[vertex]

// 		for j, v := range node {
// 			k, _ := g.vertex[v]
// 			g.reachabMatrix[j][k] = true
// 		}

// 	}

// 	return g
// }

func (g *Data) getDataCSV(path string) {

	csvfile, err := os.Open(path)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", path, err)
	}
	r := csv.NewReader(csvfile)
	record, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	g.data = make(map[string][]string)

	for _, node := range record {
		for _, vertex := range node {
			_, exist := g.data[vertex]
			if exist == false {
				g.data[vertex] = make([]string, 0)
			}
		}
		g.data[node[0]] = append(g.data[node[0]], node[1])
		g.data[node[1]] = append(g.data[node[1]], node[0])
	}

	return
}

func (g *Data) dfs(vertex string) *Road {

	road := new(Road)
	road.head = vertex

	for v := range g.data {
		road.visited[v] = false
	}
	road.visited[vertex] = true


	return road
}

func (g *Data) babyDFS(road *Road) *Road{

	for _, v := range g.data[road.head] {
		road.visited[v] = true
		_, exist := road[v]
		if exist == true {
			road[v]++
		} else {
			road[v] = 1
		}
	}

	return road
}
