package dag

import "fmt"

type NodeID string

type Node struct {
	NodeIDValue NodeID      // 重命名 ID 字段为 NodeIDValue
	Data        interface{}
	Parents     []*Node
	Children    []*Node
}

// 实现 graph.Node 接口
func (n *Node) ID() string {
	return string(n.NodeIDValue)
}

func (n *Node) Attributes() map[string]string {
	return map[string]string{
		"label": fmt.Sprintf("%v", n.Data), // 使用 Data 作为节点标签
	}
}
