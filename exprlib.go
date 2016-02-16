package main

// TODO -- implement all the commented out stuff

import (
	"fmt"
	"reflect"
)

type BinaryExprStmt interface {
	compare(left, right interface{}) bool
	Stmt
}

type UnaryExprStmt interface {
	compare(operand interface{}) bool
	Stmt
}

type equalsOp struct{}
type notEqualsOp struct{}
type greaterThanOp struct{}
type greaterEqualOp struct{}
type lessThanOp struct{}
type lessEqualOp struct{}
type andOp struct{}
type orOp struct{}
type notOp struct{}
type truthOp struct{}

func init() {
	funcMap["=="] = &equalsOp{}
	funcMap["!="] = &notEqualsOp{}
	funcMap[">"] = &greaterThanOp{}
	funcMap[">="] = &greaterEqualOp{}
	funcMap["<"] = &lessThanOp{}
	funcMap["<="] = &lessEqualOp{}
	funcMap["&&"] = &andOp{}
	funcMap["||"] = &orOp{}
	funcMap["!"] = &notOp{}
	funcMap["()"] = &truthOp{}
}

func (e *equalsOp) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	return binaryCompareOp(e, ctx, args)
}

func (e *equalsOp) help() string {
	return "[\"==\", <leftOperand>, <rightOperand>]"
}

func (e *equalsOp) compare(left, right interface{}) bool {
	return left == right
}

func (n *notEqualsOp) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	return binaryCompareOp(n, ctx, args)
}

func (n *notEqualsOp) help() string {
	return "[\"==\", <leftOperand>, <rightOperand>]"
}

func (n *notEqualsOp) compare(left, right interface{}) bool {
	return left != right
}

func (g *greaterThanOp) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	return binaryCompareOp(g, ctx, args)
}

func (g *greaterThanOp) help() string {
	return "[\"==\", <leftOperand>, <rightOperand>]"
}

func (g *greaterThanOp) compare(left, right interface{}) bool {
	leftNum, rightNum, err := numberTypesOrFail(left, right)
	if err != nil {
		fatalf("> bad operand(s): %s\n", err.Error())
	}
	return leftNum > rightNum
}

func (g *greaterEqualOp) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	return binaryCompareOp(g, ctx, args)
}

func (g *greaterEqualOp) help() string {
	return "[\"==\", <leftOperand>, <rightOperand>]"
}

func (g *greaterEqualOp) compare(left, right interface{}) bool {
	leftNum, rightNum, err := numberTypesOrFail(left, right)
	if err != nil {
		fatalf(">= bad operand(s): %s\n", err.Error())
	}
	return leftNum >= rightNum
}

func (l *lessThanOp) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	return binaryCompareOp(l, ctx, args)
}

func (l *lessThanOp) help() string {
	return "[\"<\", <leftOperand>, <rightOperand>]"
}

func (l *lessThanOp) compare(left, right interface{}) bool {
	leftNum, rightNum, err := numberTypesOrFail(left, right)
	if err != nil {
		fatalf("< bad operand(s): %s\n", err.Error())
	}
	return leftNum < rightNum
}

func (l *lessEqualOp) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	return binaryCompareOp(l, ctx, args)
}

func (l *lessEqualOp) help() string {
	return "[\"<=\", <leftOperand>, <rightOperand>]"
}

func (l *lessEqualOp) compare(left, right interface{}) bool {
	leftNum, rightNum, err := numberTypesOrFail(left, right)
	if err != nil {
		fatalf("< bad operand(s): %s\n", err.Error())
	}
	return leftNum <= rightNum
}

func (a *andOp) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	return binaryLogicalOp(a, ctx, args)
}

func (a *andOp) help() string {
	return "[\"&&\", <leftOperand>, <rightOperand>]"
}

func (a *andOp) compare(left, right interface{}) bool {
	return left.(bool) && right.(bool)
}

func (o *orOp) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	return binaryLogicalOp(o, ctx, args)
}

func (o *orOp) help() string {
	return "[\"||\", <leftOperand>, <rightOperand>]"
}

func (o *orOp) compare(left, right interface{}) bool {
	return left.(bool) || right.(bool)
}

func (n *notOp) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	return unaryLogicalOp(n, ctx, args)
}

func (n *notOp) help() string {
	return "[\"!\", <operand>]"
}

func (n *notOp) compare(operand interface{}) bool {
	return !operand.(bool)
}

func (t *truthOp) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	return unaryLogicalOp(t, ctx, args)
}

func (t *truthOp) help() string {
	return "[\"()\", <operand>]"
}

func (t *truthOp) compare(operand interface{}) bool {
	return operand.(bool)
}

//
//  Non-statements functions -- support for the above.
//
func evalSubStmt(ctx map[string]interface{}, subStmt interface{}) (interface{}, error) {
	var leftRes interface{}
	var err error
	if isSlice(subStmt) {
		leftRes, err = runOneStep(ctx, subStmt.([]interface{}))
		if err != nil {
			return nil, err
		}
	} else {
		leftRes = subStmt
	}
	return leftRes, nil
}

func binaryCompareOp(e BinaryExprStmt, ctx map[string]interface{}, args []interface{}) (bool, error) {
	if len(args) != 2 {
		return false, fmt.Errorf("Wrong number of args to binary comparison operator\n")
	}
	leftRes, err := evalSubStmt(ctx, args[0])
	if err != nil {
		return false, err
	}
	rightRes, err := evalSubStmt(ctx, args[1])
	if err != nil {
		return false, err
	}
	return e.compare(leftRes, rightRes), nil
}

func binaryLogicalOp(e BinaryExprStmt, ctx map[string]interface{}, args []interface{}) (bool, error) {
	if len(args) != 2 {
		return false, fmt.Errorf("Wrong number of args (%d) to binary logical operator\n", len(args))
	}
	leftRes, err := evalSubStmt(ctx, args[0])
	if err != nil {
		return false, err
	}
	rightRes, err := evalSubStmt(ctx, args[1])
	if err != nil {
		return false, err
	}

	return e.compare(findTheTruth(leftRes), findTheTruth(rightRes)), nil
}

func unaryLogicalOp(e UnaryExprStmt, ctx map[string]interface{}, args []interface{}) (bool, error) {
	if len(args) != 1 {
		return false, fmt.Errorf("Wrong number of args (%d) to unary logical operator\n", len(args))
	}
	res, err := evalSubStmt(ctx, args[0])
	if err != nil {
		return false, err
	}
	return e.compare(res), nil
}

func outerType(arg interface{}) string {
	return reflect.ValueOf(arg).Kind().String()
}

func isSlice(arg interface{}) bool {
	return outerType(arg) == "slice"
}

func isMap(arg interface{}) bool {
	return outerType(arg) == "map"
}

func numberTypesOrFail(left, right interface{}) (float64, float64, error) {
	leftNum, err := numberTypeOrFail(left)
	if err != nil {
		return 0, 0, err
	}
	rightNum, err := numberTypeOrFail(right)
	if err != nil {
		return 0, 0, err
	}
	return leftNum, rightNum, nil
}

func numberTypeOrFail(arg interface{}) (float64, error) {

	switch arg.(type) {
	case int:
		return float64(arg.(int)), nil
	case int8:
		return float64(arg.(int8)), nil
	case int16:
		return float64(arg.(int16)), nil
	case int32:
		return float64(arg.(int32)), nil
	case int64:
		return float64(arg.(int64)), nil
	case uint:
		return float64(arg.(uint)), nil
	case uint8:
		return float64(arg.(uint8)), nil
	case uint16:
		return float64(arg.(uint16)), nil
	case uint32:
		return float64(arg.(uint32)), nil
	case uint64:
		return float64(arg.(uint64)), nil
	case float32:
		return float64(arg.(float32)), nil
	case float64:
		return arg.(float64), nil
	default:
		return 0, fmt.Errorf("Argument %+v is not a number", arg)
	}
}

func findTheTruth(arg interface{}) bool {
	switch arg.(type) {
	case bool:
		return arg.(bool)
	case uint8, uint16, uint32, uint64, int8, int16, int32, int64, uint, int, float32, float64:
		return arg != 0
	case string:
		return arg != ""
	case map[string]interface{}:
		mapArg := arg.(map[string]interface{})
		return mapArg != nil && len(mapArg) > 0
	case []interface{}:
		sliceArg := arg.([]interface{})
		return sliceArg != nil && len(sliceArg) > 0
	default:
		return false
	}
}

func evalExprStmt(ctx map[string]interface{}, stmt []interface{}) (bool, error) {
	res, err := runOneStep(ctx, stmt)
	if err != nil {
		return false, err
	}
	if result, ok := res.(bool); ok {
		return result, nil
	}
	return false, fmt.Errorf("Expression evaluation yield non-boolean result")
}
