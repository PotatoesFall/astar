package astar

import (
	"git.ultraware.nl/martin/temp/minheap"
)

type Frame[T comparable] struct {
	Path      []T
	Traversed map[[2]T]struct{}
}

func Solve[T comparable](start, goal T, neighbors func(T) []T, h func(T) float64, d func(T, T) float64) ([]Frame[T], bool) {
	openSet := minheap.New[T]()
	openSet.Insert(float64(0), start)

	gScore := map[T]float64{}
	fScore := map[T]float64{}

	gScore[start] = 0
	fScore[start] = h(start)

	cameFrom := map[T]T{}

	frames := []Frame[T]{{
		Path:      []T{start},
		Traversed: make(map[[2]T]struct{}),
	}}

	traversed := make(map[[2]T]struct{})

	for openSet.Len() > 0 {
		node, _ := openSet.Extract()

		frames = append(frames, Frame[T]{
			Path:      getPath(cameFrom, node),
			Traversed: copyMap(traversed),
		})

		if node == goal {
			return frames, true
		}

		for _, neighbor := range neighbors(node) {
			dist := d(node, neighbor)
			newG := gScore[node] + dist
			traverse(traversed, node, neighbor)
			if oldG, ok := gScore[neighbor]; !ok || oldG > newG {
				cameFrom[neighbor] = node
				gScore[neighbor] = newG

				newF := newG + h(neighbor)
				fScore[neighbor] = newF
				openSet.Insert(newF, neighbor)
			}
		}
	}

	return frames, false
}

func traverse[T comparable](m map[[2]T]struct{}, n1, n2 T) {
	_, ok1 := m[[2]T{n1, n2}]
	_, ok2 := m[[2]T{n2, n1}]

	if !ok1 && !ok2 {
		m[[2]T{n1, n2}] = struct{}{}
	}
}

func getPath[T comparable](cameFrom map[T]T, node T) []T {
	path := []T{node}
	node, ok := cameFrom[node]
	for ok {
		path = append(path, node)
		node, ok = cameFrom[node]
	}
	return path
}

func copyMap[K comparable, V any](a map[K]V) map[K]V {
	b := make(map[K]V, len(a))
	for k, v := range a {
		b[k] = v
	}
	return b
}
