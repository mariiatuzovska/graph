package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Graph structure
type Graph struct {
	graph map[string][]string // моя структура графу, по феншую, я так бачу
	road  map[string]int      // певний шлях з деякої вершини
	q     []queue             // черга, необхідна для коректного визначення шляху
}

type queue struct {
	vertex string // ім'я вершини
	step   int    // крок, на якому ми вперше зустріли задану вершину
}

func main() {

	g := new(Graph)
	g.getDataCSV("G1.csv")
	g.road = make(map[string]int)
	g.q = make([]queue, 0)

	// // можна подивитись кількість вузлів від заданої вершини до кожної доступної
	// g.bfs("192.168.0.5", 0)
	// fmt.Println(g.road)

	//

	// (1)
	// ping 192.168.0.1 з хоста 192.168.0.5
	g.ping("192.168.0.5", "192.168.0.1")

	//

	// // (2)
	// // визначено максимальний час затримки 6 мс
	// g.maxTime(3)

	//

	// // (3)
	// //визначено, що всі вершини пов'язані
	g.searchDisabled()

	//

	// // визначено 126 вершин і записано до файлу
	g.writeVertexCSV("vertex.csv")

	// // визначено матрицю досяжності
	// m := g.reachabMatrix()
	// fmt.Println(m)
	// // можна перевірити матрицю на симетричність
	// testMatrix()

	// // визначено кількість витих пар
	g.searchTwistedPairs()

	//

	// // (4)
	// // створено довільний граф, у якого максимальний час затримки 4 мс
	// // у мапі ключі не заіндексовані, тому при кожному виклику в нас різні графи
	g.writeGraphCSV("G2.csv")
	// // коректність графу
	g2 := new(Graph)
	g2.getDataCSV("G2.csv")
	// // визначено максимальний час затримки 4 мс
	g2.maxTime(2)

}

func (g *Graph) getDataCSV(path string) {

	csvfile, _ := os.Open(path)
	r := csv.NewReader(csvfile)
	record, _ := r.ReadAll()
	g.graph = make(map[string][]string)

	for _, node := range record {
		for _, vertex := range node {
			_, exist := g.graph[vertex]
			if exist == false {
				g.graph[vertex] = make([]string, 0) // запишемо деяку вершину, яку вперше зустріли
			}
		}
		g.graph[node[0]] = append(g.graph[node[0]], node[1]) // кожна вершина це ключ, який вміщує масив інших вершин
		g.graph[node[1]] = append(g.graph[node[1]], node[0]) // з якими пов'язаний
	}

	return
}

func (g *Graph) bfs(vertex string, step int) { // складність O(V+E)

	if step != 0 {
		g.q = append(g.q[:0], g.q[1:]...) // зміщуємо чергу, видаливши перший елемент, який є вхідним
	} else {
		g.road[vertex] = step // визначимо деяку вершину, на першому кроці
	}

	for _, v := range g.graph[vertex] { // подивимося, з якими вершинами v пов'язана вхідна вершина vertex
		_, exist := g.road[v] // якщо у нашому шляху такої вершини немає
		if exist == false {
			g.q = append(g.q, queue{vertex: v, step: step + 1}) // запишимо її до черги
			g.road[v] = step + 1                                // запам'ятовуємо, що ми вже десь бачили таку вершину
		}
	}

	if len(g.q) != 0 {
		g.bfs(g.q[0].vertex, g.q[0].step) // подаємо на вхід перший елемент черги
	}

	return
}

func (g *Graph) ping(v1, v2 string) {

	g.road = make(map[string]int)
	g.q = make([]queue, 0)

	g.bfs(v1, 0)
	s, _ := g.road[v2]

	fmt.Println("PING", v2, ": time =", s*2, "ms")

}

func (g *Graph) maxTime(max int) { // видасть усі пари, для яких час затримки >= max*2 ms

	v1, v2 := "", ""

	for vertex := range g.graph {

		fmt.Println("Check", vertex)

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

func (g *Graph) searchDisabled() {

	min, disabled := 126, ""

	for vertex := range g.graph {

		g.road, g.q = make(map[string]int), make([]queue, 0)
		g.bfs(vertex, 0) // побудуємо пошук в глибину для деякої вершини

		if len(g.road) == 126 { // всього вершин 126, якщо кількість рівна 126 -
			break // немає сенсу шукати не пов'язані вершини
		}

		if min > len(g.road) {
			min = len(g.road)
			disabled = vertex
		}
	}

	if min < 126 {
		fmt.Println("Probably disabled", disabled, "has", min, "nodes")
	} else {
		fmt.Println("All nodes concatenated.")
	}

}

func (g *Graph) searchTwistedPairs() {

	matrix, count := g.reachabMatrix(), 0
	for i := 0; i < 125; i++ {
		for j := i + 1; j < 125; j++ {
			if matrix[i][j] == true {
				count++
			}
		}
	}
	fmt.Println("There are", count, "twisted pairs")
}

func (g *Graph) reachabMatrix() [][]bool {

	rMatrix := make([][]bool, 126)
	for i := 0; i < 126; i++ {

		rMatrix[i] = make([]bool, 126)

		for j := 0; j < 126; j++ {
			rMatrix[i][j] = false
		}
		for _, v := range g.graph["192.168.0."+strconv.Itoa(i+1)] { // пам'ятаємо, що вершина 192.168.0.х
			j, _ := strconv.Atoi(strings.Replace(v, "192.168.0.", "", -1)) // знаходиться під індексом х-1
			rMatrix[i][j-1] = true
		}
		//rMatrix[i][i] = true
	}

	return rMatrix
}

func (g *Graph) writeGraphCSV(path string) {

	count := 771 // the count of twisted pairs

	csvfile, _ := os.Create(path)
	csvwriter := csv.NewWriter(csvfile)
	data := make([][]string, count)

	// v[1] : 125  type of relationship 1 : all/{v[1]} (125 ребер)
	// v[2] : 124  type of relationship 1 : all/{v[1], v[2]} (124 ребер)
	// v[3] : 123  type of relationship 1 : all/{v[1], v[2], v[3]} (123 ребра)
	// v[4] : 122  type of relationship 1 : all/{v[1], v[2], v[3], v[4]} (122)
	// v[5] : 121  type of relationship 1 : all/{v[1], v[2], v[3], v[4], v[5]} (121)
	// v[6] : 120  type of relationship 1 : all/{v[1], v[2], v[3], v[4], v[5], v[6]} (120)
	// total : 735
	// v[7] : 36   type of relationship 1 : v1,v2..,v36 vertexes in V : all/{v[1], v[2], v[3], v[4], v[5], v[6], v[7]} (36)
	// total : 771

	i, t, selectedV := 0, 0, make(map[string]int)

	for v1 := range g.graph {

		selectedV[v1] = t + 1 // обираємо 7 довільних вершин

		if t < 6 {

			for v2 := range g.graph { // поєднуємо усі вершини із заданою v1,
				_, exist := selectedV[v2] // не поєднуємо, якщо вже з'єднали
				if exist == false {
					data[i] = make([]string, 2)
					data[i][0] = v1
					data[i][1] = v2
					i++
				}
			}

		} else { // для сьомої вершини необхідно 36 ребер

			z := i

			for v2 := range g.graph { // поєднуємо усі вершини із заданою v1,

				if i-z > 35 {
					break
				}
				_, exist := selectedV[v2] // не поєднуємо, якщо вже з'єднали
				if exist == false {
					data[i] = make([]string, 2)
					data[i][0] = v1
					data[i][1] = v2
					i++
				}
			}

		}

		t++

		if t == 7 {
			break
		}
	}

	for _, row := range data {
		_ = csvwriter.Write(row)
	}
	csvwriter.Flush()
	csvfile.Close()

	fmt.Println(i, "nodes has written successfully into", path)
}

func (g *Graph) writeVertexCSV(path string) {

	csvfile, _ := os.Create(path)
	csvwriter := csv.NewWriter(csvfile)
	data, i := make([][]string, len(g.graph)), 0

	for v := range g.graph {
		data[i] = make([]string, 1)
		data[i][0] = v
		i++
	}
	for _, row := range data {
		_ = csvwriter.Write(row)
	}
	csvwriter.Flush()
	csvfile.Close()

	fmt.Println(i, "vertices has written successfully into", path)
}

func testMatrix(path string) { // перевіримо матрицю на симетричність

	g := new(Graph)
	g.getDataCSV(path)

	matrix := g.reachabMatrix()

	for i := 0; i < 126; i++ {
		for j := 0; j < 126; j++ {
			if matrix[i][j] != matrix[j][i] {
				log.Fatal("Упс, матриця не симетрична! Дивись: [", i, "] [", j, "] =", matrix[i][j], "!= [", j, "] [", i, "] =", matrix[j][i])
			}
		}
	}
}
