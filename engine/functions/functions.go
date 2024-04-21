package functions

import (
	"encoding/json"
	"errors"

	"gopkg.in/Knetic/govaluate.v2"
)

var LocalFunctions = map[string]govaluate.ExpressionFunction{
	"sum": Sum,
	"mul": Mul,
	"test": Test,
}

func GetLocalfunctions() map[string]govaluate.ExpressionFunction { return LocalFunctions }

func RegisterLocalFunction(key string,f govaluate.ExpressionFunction) {
	LocalFunctions[key] = f
}

func Sum(arguments ...interface{}) (interface{}, error) {
	if len(arguments) < 2 {
		return 0, errors.New("invalid number of arguments")
	}
	args1 := arguments[0].(float64)
	args2 := arguments[1].(float64)
	return genericSum([]float64{args1, args2}), nil
}

func genericSum[N int64 | float64](nums []N) N {
    var sum N
    for _, num := range nums {
        sum += num
    }
    return sum
}

func Mul(arguments ...interface{}) (interface{}, error) {
	if len(arguments) < 2 {
		return 0, errors.New("invalid number of arguments")
	}
	args1 := arguments[0].(float64)
	args2 := arguments[1].(float64)
	return genericMul([]float64{args1, args2}), nil
}

func genericMul[N int64 | float64](nums []N) N {
    var res N = 1
    for _, num := range nums {
        res *= num
    }
    return res
}

type T struct {
	A int
	B string
	C float64
	D []int
}

func Test(arguments ...interface{}) (interface{}, error) {
	
	return structToMap(T{
		A: 1,
		B: "2",
		C: 3,
		D: []int{4,5},
	})
}

func structToMap(obj interface{}) (map[string]interface{}, error) {
    var result map[string]interface{}

    jsonBytes, err := json.Marshal(obj)
    if err != nil {
        return nil, err
    }

    err = json.Unmarshal(jsonBytes, &result)
    if err != nil {
        return nil, err
    }

    return result, nil
}
