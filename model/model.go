package model

import (
	"encoding/json"
)

type Edge struct {
	ID string `json:"id"`
	SourceId string `json:"sourceId"`
	TargetId string `json:"targetId"`
	SourceHandler string `json:"sourceHandle"`
}

type NodeBase struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Kind string `json:"type"`
	Content json.RawMessage `json:"content"`	// 定义为string，不同的node内容不同，自己按需转换
}

type FunctionNodeContent struct {
	FuncName string	`json:"funcName"`
	Args []*Args `json:"args"`
}

type Args struct {
	Value string `json:"value"`
	Kind string `json:"type"`
}

type SwitchNodeContent struct {
	Statements []*SwitchStatement `json:"statements"`
}

type SwitchStatement struct {
	ID string `json:"id"`
	Condition string `json:"condition"`
}

type DecisionTableNodeContent struct {
	Rules []map[string]string `json:"rules"`
	Inputs []*DecisionTableInputField `json:"inputs"`
	Outputs []*DecisionTableOnputField `json:"outputs"`
}

type DecisionTableInputField struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Field string `json:"field"`
}

type DecisionTableOnputField struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Field string `json:"field"`
}

type ExpressionNodeContent struct {
	Expressions []struct{
		ID string `json:"id"`
		Key string `json:"key"`
		Value string `json:"value"`
	} `json:"expressions"`
}

type Graph struct {
	Nodes []*NodeBase `json:"nodes"`
	Edges []*Edge `json:"edges"`
}

type GraphResponse struct {
	Performance string `json:"performance"`
	Result map[string]interface{} `json:"result"`
	Trace map[string]*GraphTrace `json:"trace"`
}

type GraphTrace struct {
	Input map[string]interface{} `json:"input"`
	Output map[string]interface{} `json:"output"`
	Name string `json:"name"`
	ID string `json:"id"`
	Performance string  `json:"performance"`
	TraceData map[string]interface{} `json:"traceData"`
}

type NodeResult struct {
	Result map[string]interface{} `json:"result"`
	Trace map[string]interface{} `json:"trace"`
}