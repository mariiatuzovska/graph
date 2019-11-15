package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

type graph struct {
	graph map[string][]string
	road  map[string]int
	q     []queue
}

type queue struct {
	vertex string
	step   int
}

func main() {

	g := new(graph)
	g.getDataCSV("G1.csv")
	g.road = make(map[string]int)
	g.q = make([]queue, 0)

	// g.bfs("192.168.0.5", 0)
	// fmt.Println(g.road)

	// g.maxTime()

	g.searchDisabled()

}

func (g *graph) getDataCSV(path string) {

	csvfile, _ := os.Open(path)
	r := csv.NewReader(csvfile)
	record, _ := r.ReadAll()
	g.graph = make(map[string][]string)

	for _, node := range record {
		for _, vertex := range node {
			_, exist := g.graph[vertex]
			if exist == false {
				g.graph[vertex] = make([]string, 0)
			}
		}
		g.graph[node[0]] = append(g.graph[node[0]], node[1])
		g.graph[node[1]] = append(g.graph[node[1]], node[0])
	}

	return
}

func (g *graph) bfs(vertex string, step int) {

	g.road[vertex] = step

	if len(g.q) != 0 {
		g.q = append(g.q[:0], g.q[1:]...)
	}

	for _, v := range g.graph[vertex] {
		_, exist := g.road[v]
		if exist == false {
			g.q = append(g.q, queue{vertex: v, step: step + 1})
		}
	}

	if len(g.q) != 0 {
		g.bfs(g.q[0].vertex, g.q[0].step)
	}

	return
}

func (g *graph) maxTime() {

	v1, v2, max := "", "", 4

	for vertex := range g.graph {

		g.road, g.q = make(map[string]int), make([]queue, 0)
		g.bfs(vertex, 0)

		for v, steps := range g.road {
			if steps >= max {
				v1, v2, max = vertex, v, steps
				fmt.Println("From", v1, "to", v2, "is", max*2, "ms")
			}
		}
	}

}

func (g *graph) searchDisabled() {

	min, disabled := 1000, ""

	for vertex := range g.graph {

		g.road, g.q = make(map[string]int), make([]queue, 0)
		g.bfs(vertex, 0)

		if min > len(g.road) {
			min = len(g.road)
			disabled = vertex
		}
	}

	fmt.Println("Probably disabled", disabled, "has", min, "nodes")

}

func writeDataCSV(path string, data [][]string) {

	csvfile, _ := os.Create(path)
	csvwriter := csv.NewWriter(csvfile)
	for _, row := range data {
		_ = csvwriter.Write(row)
	}
	csvwriter.Flush()
	csvfile.Close()

	fmt.Println("Data has written successfully into", path)
}
