package algorithm

import (
	"fmt"
	"golina/container/graph"
	"testing"
)

func TestBFS(t *testing.T) {
	g, err := graph.NewGraphFromJSON("../test.json", "graph")
	if err != nil {
		panic(err)
	}
	// fmt.Println(g)
	nodes, err := BFS(g, graph.StringID("A"))

	if err != nil || len(nodes) != 8 {
		fmt.Println(err, len(nodes))
		t.Fail()
	}

	// start vertex
	if nodes[0].ID().String() != "A" {
		t.Fail()
	}

	// nodes with one depths between
	res1 := map[string]struct{}{
		"B": struct{}{},
		"D": struct{}{},
		"G": struct{}{},
		"H": struct{}{},
	}
	for _, n := range nodes[1:5] {
		if _, found := res1[n.ID().String()]; !found {
			t.Fail()
		}
	}

	// nodes with two depths between
	res2 := map[string]struct{}{
		"E": struct{}{},
		"F": struct{}{},
		"C": struct{}{},
	}
	for _, n := range nodes[5:] {
		if _, found := res2[n.ID().String()]; !found {
			t.Fail()
		}
	}
}
