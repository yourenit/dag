package graph

import (
	"fmt"
	"time"

	"github.com/yourenit/dag_engine/engine/node"
	"github.com/yourenit/dag_engine/engine/traversal"

	"github.com/yourenit/dag_engine/model"

	"github.com/heimdalr/dag"
)



func Evaluate(g *model.Graph, inputs map[string]interface{}) (*model.GraphResponse,error) {
	start := time.Now().Local().UnixMicro()
	d := dag.NewDAG()
	for _, node := range g.Nodes {
		if err := d.AddVertexByID(node.ID, node); err != nil {
			return nil, err
		}
		
	}
	for _, edge := range g.Edges {
		if err := d.AddEdge(edge.SourceId, edge.TargetId);err != nil {
			return nil, err
		}
	}

	nodeTrace := make(map[string]*model.GraphTrace)
	walker := traversal.New(d)
	for {
		nid,traceData,err:=walker.Next(d, g)
		if err != nil || nid == "" {
			break
		}
		vertex, err := d.GetVertex(nid)
		if err != nil {
			continue
		}

		mergeInputs := walker.MergeIncoming(d, nid)
		switch vertex.(*model.NodeBase).Kind {
		case "inputNode":
			walker.NodeData[nid] = inputs
			nodeTrace[nid] = &model.GraphTrace{
				Name: vertex.(*model.NodeBase).Name,
				ID: vertex.(*model.NodeBase).ID,
				Performance: fmt.Sprintf("%.2fms", float64(time.Now().Local().UnixMicro() - start)/1000.0),
			}
		case "outputNode":
			walker.NodeData[nid] = mergeInputs
			nodeTrace[nid] = &model.GraphTrace{
				Name: vertex.(*model.NodeBase).Name,
				ID: vertex.(*model.NodeBase).ID,
			}
			return &model.GraphResponse{
				Result: mergeInputs,
				Performance: fmt.Sprintf("%.2fms", float64(time.Now().Local().UnixMicro() - start)/1000.0),
				Trace: nodeTrace,
			}, nil
		case "switchNode":
			walker.NodeData[nid] = mergeInputs
			nodeTrace[nid] = &model.GraphTrace{
				Input: mergeInputs,
				Output: mergeInputs,
				Name: vertex.(*model.NodeBase).Name,
				ID: vertex.(*model.NodeBase).ID,
				TraceData: traceData,
				Performance: fmt.Sprintf("%.2fms", float64(time.Now().Local().UnixMicro() - start)/1000.0),
			}
		case "decisionTableNode":
			result :=node.DecisionTableHandle(vertex.(*model.NodeBase).Content, mergeInputs)
			walker.NodeData[nid] = result.Result
			nodeTrace[nid] = &model.GraphTrace{
				Input: mergeInputs,
				Output: result.Result,
				Name: vertex.(*model.NodeBase).Name,
				ID: vertex.(*model.NodeBase).ID,
				TraceData: result.Trace,
				Performance: fmt.Sprintf("%.2fms", float64(time.Now().Local().UnixMicro() - start)/1000.0),
			}
		case "functionNode":
			walker.NodeData[nid] = mergeInputs
		case "expressionNode":
			result := node.ExpressionNodeHandle(vertex.(*model.NodeBase).Content, mergeInputs)
			walker.NodeData[nid] = result.Result
			nodeTrace[nid] = &model.GraphTrace{
				Input: mergeInputs,
				Output: result.Result,
				Name: vertex.(*model.NodeBase).Name,
				ID: vertex.(*model.NodeBase).ID,
				TraceData: result.Trace,
				Performance: fmt.Sprintf("%.2fms", float64(time.Now().Local().UnixMicro() - start)/1000.0),
			}
		}
	}

	return nil, fmt.Errorf("Graph did not halt. Missing output node")
	
}
