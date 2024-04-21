package node

import (
	"encoding/json"

	"github.com/yourenit/dag_engine/engine/functions"
	"github.com/yourenit/dag_engine/model"

	"gopkg.in/Knetic/govaluate.v2"
)


func GetValidStatement(content json.RawMessage, inputs map[string]interface{}) map[string]*model.SwitchStatement {
	sContent := model.SwitchNodeContent{}
	err := json.Unmarshal(content,&sContent)
	if err != nil {
		return nil
	}
	validStatements := make(map[string]*model.SwitchStatement)
	for _,statement := range sContent.Statements {
		expr, err := govaluate.NewEvaluableExpression(statement.Condition)
		if err != nil {
			continue
		}
		result,err := expr.Evaluate(inputs)
		if err != nil {
			continue
		}
		if result.(bool) {
			validStatements[statement.ID] = statement
		}
	}
	return validStatements
}


func DecisionTableHandle(content json.RawMessage, inputs map[string]interface{}) (*model.NodeResult) {
	traceData := make(map[string]interface{})
	dContent := model.DecisionTableNodeContent{}
	err := json.Unmarshal(content,&dContent)
	if err != nil {
		return &model.NodeResult{
			Trace: map[string]interface{}{
				"err": err.Error(),
			},
		}
	}

	hitRules := make(map[string]interface{},0)
	for _,input := range dContent.Inputs {
		// 从rule中找input对应的id
		for _,rule := range dContent.Rules {
			if ruleValue,ok := rule[input.ID];ok {
				expr,err := govaluate.NewEvaluableExpression(ruleValue)
				if err != nil {
					continue
				}
				result, err := expr.Evaluate(inputs)
				if err != nil {
					continue
				}
				if result.(bool) {
					hitRules[rule["_id"]] = rule
				}
			}
		}
	}

	outputs := make(map[string]string)
	for _, output := range dContent.Outputs {
		for _,rule := range dContent.Rules {
			if _,ok := hitRules[rule["_id"]]; !ok {
				continue
			}
			outputs[output.ID] = output.Field
		}
	}

	results := make(map[string]interface{})
	
	for _,rule := range dContent.Rules {
		for outputId,field := range outputs {
			if fieldValue,ok := rule[outputId]; ok {
				results[field] = fieldValue
			}
		}
	}
	
	return &model.NodeResult{
		Result: results,
		Trace: traceData,
	}
}

func FunctionNodeHandle(content json.RawMessage, inputs map[string]interface{}) map[string]interface{} {
	funcContent := model.FunctionNodeContent{}
	err := json.Unmarshal(content,&funcContent)
	if err != nil {
		return nil
	}

	return nil
}

func ExpressionNodeHandle(content json.RawMessage, inputs map[string]interface{}) (*model.NodeResult) {
	exprContent := model.ExpressionNodeContent{}
	err := json.Unmarshal(content,&exprContent)
	if err != nil {
		return &model.NodeResult{
			Trace: map[string]interface{}{
				"err": err.Error(),
			},
		}
	}
	results := make(map[string]interface{})
	traceData := make(map[string]interface{})
	for _, expression := range exprContent.Expressions {
		expr, err := govaluate.NewEvaluableExpressionWithFunctions(expression.Value, functions.GetLocalfunctions())
		if err != nil {
			return &model.NodeResult{
				Trace: map[string]interface{}{
					"err": err.Error(),
				},
			}
		}
		result, err := expr.Evaluate(inputs)
		if err != nil {
			return &model.NodeResult{
				Trace: map[string]interface{}{
					"err": err.Error(),
				},
			}
		}
		results[expression.Key] = result
		traceData[expression.Key] = map[string]interface{}{
			"expression": expression.Value,
			"result": result,
		}
	}

	return &model.NodeResult{
		Result: results,
		Trace: traceData,
	}
}
