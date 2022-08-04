package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"git.ultraware.nl/martin/temp/astar"
)

const (
	nNodes           = 500
	imgSize          = 1000
	neighborDistance = 80
	maxNeighbors     = 10

	random = false
	nEdges = 100

	frametime     = 0.1
	lastFrameTime = 5.0
)

type Graph []*Node

type Node struct {
	X, Y      float64
	Neighbors []*Node
}

type Edge struct {
	A, B *Node
}

func (e Edge) Length() float64 {
	return distance(e.A, e.B)
}

func main() {
	var seed int64
	seedStr, ok := os.LookupEnv(`SEED`)
	seed, err := strconv.ParseInt(seedStr, 10, 64)
	if !ok || err != nil {
		seed = time.Now().UnixNano()
		fmt.Println(`random seed:`, seed)
	}
	rand.Seed(seed)

	graph := getRandomGraph()

	start, goal := graph[0], graph[nNodes+1]
	startTime := time.Now()
	frames, ok := astar.Solve(
		start, goal,
		func(n *Node) []*Node { return n.Neighbors },
		func(n *Node) float64 { return h(n, goal) },
		distance)

	if !ok {
		fmt.Println(`no solution`)
	} else {
		fmt.Printf("solved in %d ms\n", int(time.Since(startTime).Milliseconds()))
	}

	makeSVG(graph, frames)
}

func h(n, goal *Node) float64 {
	return distance(n, goal)
}

func getRandomGraph() Graph {
	graph := []*Node{{0, 0, nil}}
	for i := 0; i < nNodes; i++ {
		graph = append(graph, &Node{float64(rand.Intn(imgSize-1) + 1), float64(rand.Intn(imgSize-1) + 1), nil})
	}
	graph = append(graph, &Node{imgSize, imgSize, nil})

	if random {
		for i := 0; i < nEdges; i++ {
			a, b := rand.Intn(nNodes+2), rand.Intn(nNodes+2)
			graph[a].Neighbors = append(graph[a].Neighbors, graph[b])
			graph[b].Neighbors = append(graph[b].Neighbors, graph[a])
		}
	} else {
		for i := range graph {
			if i == len(graph)-1 {
				break
			}

			for j := range graph[i+1:] {
				j = j + i + 1
				if distance(graph[i], graph[j]) < neighborDistance && len(graph[i].Neighbors) < maxNeighbors && len(graph[j].Neighbors) < maxNeighbors {
					graph[i].Neighbors = append(graph[i].Neighbors, graph[j])
					graph[j].Neighbors = append(graph[j].Neighbors, graph[i])
				}
			}
		}
	}

	return graph
}

func makeSVG(g Graph, frames []astar.Frame[*Node]) {
	var svg bytes.Buffer

	svg.WriteString("<svg viewBox=\"0 0 1000 1000\" xmlns=\"http://www.w3.org/2000/svg\">")
	svg.WriteString("<rect width=\"1000\" height=\"1000\" fill=\"white\" />")

	svg.WriteString("<defs><style type=\"text/css\">")
	for i := range frames {
		if i != 0 {
			svg.WriteByte(',')
		}
		svg.WriteString(fmt.Sprintf("#_%d", i))
	}
	svg.WriteString("{visibility: hidden}</style></defs>")

	edges := make(map[Edge]struct{})

	for _, node := range g {
		svg.WriteString(fmt.Sprintf("<circle cx=\"%d\" cy=\"%d\" r=\"5\" fill=\"black\"/>", int(node.X), int(node.Y)))
		for _, n := range node.Neighbors {
			edges[Edge{node, n}] = struct{}{}
		}
	}

	for edge := range edges {
		svg.WriteString(fmt.Sprintf("<line x1=\"%d\" x2=\"%d\" y1=\"%d\" y2=\"%d\" stroke=\"black\" />", int(edge.A.X), int(edge.B.X), int(edge.A.Y), int(edge.B.Y)))
	}

	for i, frame := range frames {
		svg.WriteString(fmt.Sprintf("<g id=\"_%d\">", i))

		for edge := range frame.Traversed {
			svg.WriteString(fmt.Sprintf("<line x1=\"%d\" x2=\"%d\" y1=\"%d\" y2=\"%d\" stroke=\"red\" stroke-width=\"5\"/>", int(edge[0].X), int(edge[1].X), int(edge[0].Y), int(edge[1].Y)))
		}

		var previousNode *Node
		for i, node := range frame.Path {
			if i == 0 {
				previousNode = node
				continue
			}

			svg.WriteString(fmt.Sprintf("<line x1=\"%d\" x2=\"%d\" y1=\"%d\" y2=\"%d\" stroke=\"#00ff00\" stroke-width=\"5\"/>", int(node.X), int(previousNode.X), int(node.Y), int(previousNode.Y)))

			previousNode = node
		}
		svg.WriteString("</g>")
	}

	svg.WriteString("<defs><style type=\"text/css\">")
	for i := range frames {
		if i == len(frames)-1 {
			break
		}
		if i != 0 {
			svg.WriteByte(',')
		}
		svg.WriteString(fmt.Sprintf("#_%d", i))
	}
	svg.WriteString(fmt.Sprintf("{animation: %fs linear _k infinite}", float64(len(frames)-1)*frametime+lastFrameTime))
	svg.WriteString(fmt.Sprintf("#_%d{animation: %fs linear _j infinite}", len(frames)-1, float64(len(frames)-1)*frametime+lastFrameTime))
	delay := 0.0
	for i := range frames {
		svg.WriteString(fmt.Sprintf("#_%d {animation-delay: %fs}", i, frametime*delay))
		delay += 1
	}
	frac := 100 / float64(len(frames))
	svg.WriteString(fmt.Sprintf("@keyframes _k {0%%, %f%% {visibility: visible } %f%%, 100%% {visibility: hidden }}", frac, frac))
	frac = 100 * lastFrameTime / (float64(len(frames)-1)*frametime + lastFrameTime)
	svg.WriteString(fmt.Sprintf("@keyframes _j {0%%, %f%% {visibility: visible } %f%%, 100%% {visibility: hidden }}", frac, frac))
	svg.WriteString("</style></defs>")

	svg.WriteString("</svg>")

	if err := os.WriteFile(`graph.svg`, svg.Bytes(), 0o654); err != nil {
		panic(err)
	}
}

func distance(n1, n2 *Node) float64 {
	return math.Sqrt(math.Pow(n1.X-n2.X, 2) + math.Pow(n1.Y-n2.Y, 2))
}
