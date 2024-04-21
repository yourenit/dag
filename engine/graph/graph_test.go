package graph

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/yourenit/dag_engine/model"
)

func TestGraphModel(t *testing.T) {
	body,err := os.ReadFile("graph_1.json")
	if err != nil {
		t.Fatal(err.Error())
	}
	graph := &model.Graph{}
	err = json.Unmarshal(body,graph)
	if err != nil {
		t.Fatal(err.Error())
	}

	for _,node := range graph.Nodes {
		t.Logf("\nnode: %+v",node)
	}
	for _,edge := range graph.Edges {
		t.Logf("\nedge: %+v",edge)
	}
	// t.Logf("\nnodes: %+v\nedges: %+v",content.Nodes[0], content.Edges[0])
}