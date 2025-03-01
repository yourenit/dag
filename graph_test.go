package dag

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestCreateNode(t *testing.T) {
	g := NewGraph()
	nodeID := NodeID("node1")
	data := "Node 1 Data"
	node := g.CreateNode(nodeID, data)

	if node == nil {
		t.Errorf("CreateNode returned nil")
		t.FailNow()
	}
	if node.NodeIDValue != nodeID {
		t.Errorf("CreateNode returned node with incorrect ID: got %v, want %v", node.NodeIDValue, nodeID)
		t.FailNow()
	}
	if node.Data != data {
		t.Errorf("CreateNode returned node with incorrect Data: got %v, want %v", node.Data, data)
		t.FailNow()
	}
	if _, exists := g.nodes[nodeID]; !exists {
		t.Errorf("Node not added to graph's nodes map")
		t.FailNow()
	}
}

func TestAddEdge(t *testing.T) {
	g := NewGraph()
	node1 := g.CreateNode("node1", "Node 1")
	node2 := g.CreateNode("node2", "Node 2")

	err := g.AddEdge("node1", "node2")
	if err != nil {
		t.Errorf("AddEdge returned error: %v", err)
		t.FailNow()
	}

	if len(node1.Children) != 1 || node1.Children[0] != node2 {
		t.Errorf("Edge not added to parent node's children")
		t.FailNow()
	}
	if len(node2.Parents) != 1 || node2.Parents[0] != node1 {
		t.Errorf("Edge not added to child node's parents")
		t.FailNow()
	}
}

func TestDeleteNode(t *testing.T) {
	g := NewGraph()
	g.CreateNode("node1", "Node 1")
	node2 := g.CreateNode("node2", "Node 2")
	g.AddEdge("node1", "node2")

	err := g.DeleteNode("node1")
	if err != nil {
		t.Errorf("DeleteNode returned error: %v", err)
		t.FailNow()
	}

	if _, exists := g.nodes["node1"]; exists {
		t.Errorf("Node not deleted from graph's nodes map")
		t.FailNow()
	}
	if len(node2.Parents) != 0 {
		t.Errorf("Parent edge not removed from child node")
		t.FailNow()
	}
}

func TestDeleteEdge(t *testing.T) {
	g := NewGraph()
	//nolint:unused
	node1 := g.CreateNode("node1", "Node 1")
	node2 := g.CreateNode("node2", "Node 2")
	g.AddEdge("node1", "node2")

	err := g.DeleteEdge("node1", "node2")
	if err != nil {
		t.Errorf("DeleteEdge returned error: %v", err)
		t.FailNow()
	}

	if len(node1.Children) != 0 {
		t.Errorf("Edge not removed from parent node's children")
		t.FailNow()
	}
	if len(node2.Parents) != 0 {
		t.Errorf("Edge not removed from child node's parents")
		t.FailNow()
	}
}

func TestGetParentsChildren(t *testing.T) {
	g := NewGraph()
	//nolint:unused
	node1 := g.CreateNode("node1", "Node 1")
	//nolint:unused
	node2 := g.CreateNode("node2", "Node 2")
	g.AddEdge("node1", "node2")

	parents, _ := g.GetParents("node2")
	if len(parents) != 1 || parents[0] != node1 {
		t.Errorf("GetParents returned incorrect parents")
		t.FailNow()
	}

	children, _ := g.GetChildren("node1")
	if len(children) != 1 || children[0] != node2 {
		t.Errorf("GetChildren returned incorrect children")
		t.FailNow()
	}
}

func TestDFS(t *testing.T) {
	g := NewGraph()
	node1 := g.CreateNode("node1", "Node 1")
	node2 := g.CreateNode("node2", "Node 2")
	node3 := g.CreateNode("node3", "Node 3")
	g.AddEdge("node1", "node2")
	g.AddEdge("node2", "node3")

	result, _ := g.DFS("node1")
	if len(result) != 3 || result[0] != node1 || result[1] != node2 || result[2] != node3 {
		t.Errorf("DFS returned incorrect order: got %v, want %v", result, []*Node{node1, node2, node3})
		t.FailNow()
	}
}

func TestBFS(t *testing.T) {
	g := NewGraph()
	node1 := g.CreateNode("node1", "Node 1")
	node2 := g.CreateNode("node2", "Node 2")
	node3 := g.CreateNode("node3", "Node 3")
	g.AddEdge("node1", "node2")
	g.AddEdge("node1", "node3")

	result, _ := g.BFS("node1")
	if len(result) != 3 || result[0] != node1 || !(result[1] == node2 && result[2] == node3) && !(result[1] == node3 && result[2] == node2) {
		t.Errorf("BFS returned incorrect order: got %v, want [node1, node2, node3] or [node1, node3, node2], got %v", result, []*Node{node1, node2, node3})
		t.FailNow()
	}
}

func TestTopologicalSort(t *testing.T) {
	g := NewGraph()
	node1 := g.CreateNode("node1", "Node 1")
	node2 := g.CreateNode("node2", "Node 2")
	node3 := g.CreateNode("node3", "Node 3")
	g.AddEdge("node1", "node2")
	g.AddEdge("node2", "node3")

	result, _ := g.TopologicalSort()
	if len(result) != 3 || result[0] != node1 || result[1] != node2 || result[2] != node3 {
		t.Errorf("TopologicalSort returned incorrect order: got %v, want %v", result, []*Node{node1, node2, node3})
		t.FailNow()
	}
}

func TestCycleDetection(t *testing.T) {
	g := NewGraph()
	//nolint:unused
	g.CreateNode("node1", "Node 1")
	//nolint:unused
	g.CreateNode("node2", "Node 2")
	g.AddEdge("node1", "node2")

	err := g.AddEdge("node2", "node1") // 添加环
	if err == nil {
		t.Errorf("Cycle detection failed, AddEdge should return error for cycle")
		t.FailNow()
	}
}

func TestJSONSerialization(t *testing.T) {
	g := NewGraph()
	g.CreateNode("node1", "Node 1")
	g.CreateNode("node2", "Node 2")
	g.AddEdge("node1", "node2")

	jsonData, err := json.Marshal(g)
	if err != nil {
		t.Errorf("JSON Marshaling failed: %v", err)
		t.FailNow()
	}

	g2 := NewGraph()
	err = json.Unmarshal(jsonData, g2)
	if err != nil {
		t.Errorf("JSON Unmarshaling failed: %v", err)
		t.FailNow()
	}

	if !reflect.DeepEqual(g, g2) {
		t.Errorf("Graphs before and after JSON serialization are not equal")
		t.FailNow()
	}
}

func TestSaveLoadToFile(t *testing.T) {
	g := NewGraph()
	g.CreateNode("node1", "Node 1")
	g.CreateNode("node2", "Node 2")
	g.AddEdge("node1", "node2")

	filepath := "test_dag.json"
	err := g.SaveToFile(filepath)
	if err != nil {
		t.Errorf("SaveToFile failed: %v", err)
		t.FailNow()
	}
	defer os.Remove(filepath)

	g2 := NewGraph()
	err = g2.LoadFromFile(filepath)
	if err != nil {
		t.Errorf("LoadFromFile failed: %v", err)
		t.FailNow()
	}

	if !reflect.DeepEqual(g, g2) {
		t.Errorf("Graphs before and after SaveToFile/LoadFromFile are not equal")
		t.FailNow()
	}
}