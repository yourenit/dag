package dag

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Graph struct {
	nodes map[NodeID]*Node
}

func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[NodeID]*Node),
	}
}

func (g *Graph) CreateNode(id NodeID, data interface{}) *Node {
	if _, exists := g.nodes[id]; exists {
		return nil // Or return error, node ID already exists
	}
	node := &Node{
		NodeIDValue:   id,
		Data: data,
	}
	g.nodes[id] = node
	return node
}

func (g *Graph) AddEdge(parentID NodeID, childID NodeID) error {
	parent, okParent := g.nodes[parentID]
	child, okChild := g.nodes[childID]
	if !okParent || !okChild {
		return errors.New("parent or child node not found")
	}

	if g.hasCycle(child, parent) { // Cycle detection before adding edge
		return errors.New("cycle detected, cannot add edge")
	}

	parent.Children = append(parent.Children, child)
	child.Parents = append(child.Parents, parent)
	return nil
}

// hasCycle uses DFS to detect cycles.
func (g *Graph) hasCycle(startNode *Node, targetNode *Node) bool {
	visited := make(map[*Node]bool)
	recursionStack := make(map[*Node]bool)
	return g.hasCycleUtil(startNode, targetNode, visited, recursionStack)
}

func (g *Graph) hasCycleUtil(currentNode *Node, targetNode *Node, visited map[*Node]bool, recursionStack map[*Node]bool) bool {
	visited[currentNode] = true
	recursionStack[currentNode] = true

	for _, child := range currentNode.Children {
		if child == targetNode { // Check if adding edge to targetNode will create cycle
			return true
		}
		if !visited[child] {
			if g.hasCycleUtil(child, targetNode, visited, recursionStack) {
				return true
			}
		} else if recursionStack[child] {
			return true
		}
	}

	recursionStack[currentNode] = false // remove the node from recursion stack
	return false
}

func (g *Graph) DeleteNode(nodeID NodeID) error {
	node, ok := g.nodes[nodeID]
	if !ok {
		return errors.New("node not found")
	}

	// Remove node from parents' children list
	for _, parent := range node.Parents {
		g.removeEdgeFromParent(parent, node)
	}

	// Remove node from children's parents list
	for _, child := range node.Children {
		g.removeEdgeFromChild(child, node)
	}

	delete(g.nodes, nodeID)
	return nil
}

func (g *Graph) removeEdgeFromParent(parent *Node, child *Node) {
	for i, c := range parent.Children {
		if c == child {
			parent.Children = append(parent.Children[:i], parent.Children[i+1:]...)
			return
		}
	}
}

func (g *Graph) removeEdgeFromChild(child *Node, parent *Node) {
	for i, p := range child.Parents {
		if p == parent {
			child.Parents = append(child.Parents[:i], child.Parents[i+1:]...)
			return
		}
	}
}

func (g *Graph) DeleteEdge(parentID NodeID, childID NodeID) error {
	parent, okParent := g.nodes[parentID]
	child, okChild := g.nodes[childID]
	if !okParent || !okChild {
		return errors.New("parent or child node not found")
	}

	g.removeEdgeFromParent(parent, child)
	g.removeEdgeFromChild(child, parent)
	return nil
}

func (g *Graph) GetParents(nodeID NodeID) ([]*Node, error) {
	node, ok := g.nodes[nodeID]
	if !ok {
		return nil, errors.New("node not found")
	}
	return node.Parents, nil
}

func (g *Graph) GetChildren(nodeID NodeID) ([]*Node, error) {
	node, ok := g.nodes[nodeID]
	if !ok {
		return nil, errors.New("node not found")
	}
	return node.Children, nil
}

func (g *Graph) DFS(startNodeID NodeID) ([]*Node, error) {
	startNode, ok := g.nodes[startNodeID]
	if !ok {
		return nil, errors.New("start node not found")
	}

	visited := make(map[*Node]bool)
	var result []*Node
	stack := []*Node{startNode} // 使用栈代替递归

	visited[startNode] = true // 标记起始节点为已访问

	for len(stack) > 0 {
		currentNode := stack[len(stack)-1] // 获取栈顶节点
		stack = stack[:len(stack)-1]       // 出栈
		result = append(result, currentNode)  // 将当前节点添加到结果列表

		// 遍历当前节点的子节点
		for _, child := range currentNode.Children {
			if !visited[child] { // 如果子节点未被访问过
				visited[child] = true    // 标记子节点为已访问
				stack = append(stack, child) // 入栈子节点
			}
		}
	}
	return result, nil
}

// dfsUtil (不再需要) - 迭代 DFS 版本不再需要辅助递归函数
// func (g *Graph) dfsUtil(node *Node, visited map[*Node]bool, result *[]*Node) {
// 	... (之前的递归实现) ...
// }

func (g *Graph) BFS(startNodeID NodeID) ([]*Node, error) {
	startNode, ok := g.nodes[startNodeID]
	if !ok {
		return nil, errors.New("start node not found")
	}

	visited := make(map[*NodeID]bool)
	var result []*Node
	queue := []*Node{startNode}
	visited[&startNode.NodeIDValue] = true

	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:] // Dequeue
		result = append(result, currentNode)

		for _, child := range currentNode.Children {
			if !visited[&child.NodeIDValue] {
				visited[&child.NodeIDValue] = true
				queue = append(queue, child) // Enqueue
			}
		}
	}
	return result, nil
}

func (g *Graph) TopologicalSort() ([]*Node, error) {
	inDegree := make(map[NodeID]int)
	for _, node := range g.nodes {
		inDegree[node.NodeIDValue] = 0
	}

	for _, node := range g.nodes {
		for _, child := range node.Children {
			inDegree[child.NodeIDValue]++
		}
	}

	var queue []*Node
	for _, node := range g.nodes {
		if inDegree[node.NodeIDValue] == 0 {
			queue = append(queue, node)
		}
	}

	var result []*Node
	count := 0
	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]
		result = append(result, currentNode)
		count++

		for _, child := range currentNode.Children {
			inDegree[child.NodeIDValue]--
			if inDegree[child.NodeIDValue] == 0 {
				queue = append(queue, child)
			}
		}
	}

	if count != len(g.nodes) {
		return nil, errors.New("cycle detected, topological sort not possible")
	}

	return result, nil
}

// jsonGraph 用于 JSON 序列化/反序列化
type jsonGraph struct {
	Nodes []jsonNode `json:"nodes"`
	Edges []jsonEdge `json:"edges"`
}

type jsonNode struct {
	ID   NodeID      `json:"id"`
	Data interface{} `json:"data"`
}

type jsonEdge struct {
	ParentID NodeID `json:"parent"`
	ChildID  NodeID `json:"child"`
}

// MarshalJSON 实现了 json.Marshaler 接口，用于自定义序列化
func (g *Graph) MarshalJSON() ([]byte, error) {
	jsonGraph := jsonGraph{
		Nodes: make([]jsonNode, 0, len(g.nodes)),
		Edges: make([]jsonEdge, 0),
	}

	for _, node := range g.nodes {
		jsonGraph.Nodes = append(jsonGraph.Nodes, jsonNode{
			ID:   node.NodeIDValue,
			Data: node.Data,
		})
		for _, child := range node.Children {
			jsonGraph.Edges = append(jsonGraph.Edges, jsonEdge{
				ParentID: node.NodeIDValue,
				ChildID:  child.NodeIDValue,
			})
		}
	}
	return json.Marshal(jsonGraph)
}

// UnmarshalJSON 实现了 json.Unmarshaler 接口，用于自定义反序列化
func (g *Graph) UnmarshalJSON(data []byte) error {
	jsonGraph := jsonGraph{}
	if err := json.Unmarshal(data, &jsonGraph); err != nil {
		return err
	}

	g.nodes = make(map[NodeID]*Node)
	nodeMap := make(map[NodeID]*Node) // 用于在反序列化时快速查找节点

	// 反序列化节点
	for _, jsonNode := range jsonGraph.Nodes {
		node := g.CreateNode(jsonNode.ID, jsonNode.Data)
		if node == nil {
			return fmt.Errorf("duplicate node ID: %s", jsonNode.ID)
		}
		nodeMap[node.NodeIDValue] = node
	}

	// 反序列化边
	for _, jsonEdge := range jsonGraph.Edges {
		parent, okParent := nodeMap[jsonEdge.ParentID]
		child, okChild := nodeMap[jsonEdge.ChildID]
		if !okParent || !okChild {
			return fmt.Errorf("edge refers to non-existent node: parentID=%s, childID=%s", jsonEdge.ParentID, jsonEdge.ChildID)
		}
		parent.Children = append(parent.Children, child)
		child.Parents = append(child.Parents, parent)
	}

	return nil
}

// SaveToFile 将 DAG 保存到 JSON 文件
func (g *Graph) SaveToFile(filepath string) error {
	data, err := json.MarshalIndent(g, "", "  ") // 使用缩进格式化 JSON
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, data, 0644)
}

// LoadFromFile 从 JSON 文件加载 DAG
func (g *Graph) LoadFromFile(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	return g.UnmarshalJSON(data)
}