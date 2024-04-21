package traversal

import (
	"container/list"
	"errors"
	"strconv"

	"github.com/yourenit/dag_engine/engine/node"
	"github.com/yourenit/dag_engine/model"

	"github.com/heimdalr/dag"
)

type GraphWalker struct {
	ToVisit list.List

	// 存储每个节点的输出
	NodeData map[string]map[string]interface{}
}

func New(g *dag.DAG) *GraphWalker {
	w := &GraphWalker{
		ToVisit: list.List{},
		NodeData: make(map[string]map[string]interface{}),
	}
	roots := g.GetRoots()
	for rootId := range roots {
		w.ToVisit.PushBack(rootId)
	}
	return w
}



func (w *GraphWalker) Next(g *dag.DAG, graph *model.Graph) (string, map[string]interface{}, error) {
	for nodeId := w.ToVisit.Front(); nodeId != nil; {
		nid := nodeId.Value.(string)
		vertex,err := g.GetVertex(nid)
		if err != nil {
			return "",nil, err
		}

		traceData := make(map[string]interface{})
		mergeInputs := w.MergeIncoming(g, nid)
		// 如果当前节点是switch节点，判断有效的条件
		if n := vertex.(*model.NodeBase); n != nil && n.Kind == "switchNode" {
			validStatement := node.GetValidStatement(n.Content, mergeInputs)
			if len(validStatement) == 0 {
				return "",nil, nil
			}
			traceData["statements"] = validStatement
			for _,edge := range graph.Edges {
				if edge.SourceId != nid {
					continue
				}
				if _, ok := validStatement[edge.SourceHandler]; edge.SourceHandler!= "" && !ok {
					err := g.DeleteEdge(edge.SourceId, edge.TargetId)
					if err != nil {
						return "",nil, err
					}
				}
			}
			
		}

		childrens,err := g.GetChildren(nid)
		if err != nil {
			return "", nil,err
		}

		
		for childrenId := range childrens {
			w.ToVisit.PushBack(childrenId)
		}
		w.ToVisit.Remove(nodeId)

		
		
		return nid,traceData, nil
	}
	return "",nil,errors.New("no nodes to visit")
}

func (w *GraphWalker) MergeIncoming(g *dag.DAG, nid string) map[string]interface{} {
	// 获取当前节点的所有父节点并合并输出
	parents,err := g.GetParents(nid)
	if err != nil {
		return nil
	}

	// 合并待优化，打平嵌套结构
	mergeInputs := make(map[string]interface{})
	for parentNodeId := range parents {
		parentOutput,ok := w.NodeData[parentNodeId];
		if !ok {
			continue
		}
		for key,value := range parentOutput {
			mergeInputs[key] = value
		}
	}
	
	return Flatten(mergeInputs)
}

func Flatten(m map[string]interface{}) map[string]interface{} {
    o := map[string]interface{}{}
    for k, v := range m {
        switch child := v.(type) {
        case map[string]interface{}:
            nm := Flatten(child)
            for nk, nv := range nm {
                o[k+"."+nk] = nv
            }
        case []interface{}:
            for i, val := range child {
                o[k+"."+strconv.Itoa(i)] = val
            }
        default:
            o[k] = v
        }
    }
    return o
}