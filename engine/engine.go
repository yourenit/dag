package engine

import (
	"encoding/json"

	"github.com/yourenit/dag_engine/engine/graph"
	"github.com/yourenit/dag_engine/model"
)


type Engine struct {
}

func New() *Engine {
	return &Engine{}
}



func (e *Engine) Evaluate(content string) (*model.GraphResponse, error) {
	type Args struct {
		Content *model.Graph `json:"content"`
		Context map[string]interface{} `json:"context"`
	}
	args := &Args{}
	err := json.Unmarshal([]byte(content), args)
	if err != nil {
		return nil, err
	}
	return graph.Evaluate(args.Content, args.Context)
}



