package main

import (
	"fmt"
	"time"
)

func init() {
	funcMap["sleep"] = sleep
	funcMap["repeat"] = repeat
}

func sleep(ctx map[string]interface{}, args []interface{}) error {
	secs := time.Duration(getArg(args, 0).(float64))
	time.Sleep(secs * time.Second)
	return nil
}

func repeat(ctx map[string]interface{}, args []interface{}) error {
	if len(args) != 2 {
		return fmt.Errorf("Usage: [repeat, <num>, [<stmt>...]]")
	}
	if _, ok := args[0].(float64); !ok {
		return fmt.Errorf("Repeat count must be a number")
	}
	if _, ok := args[1].([]interface{}); !ok {
		return fmt.Errorf("Repeat statements must be an array")
	}
	stmts := args[1].([]interface{})
	iterCount := int(0)
	for count := int(args[0].(float64)); count > 0; count-- {
		iterCount++
		fmt.Printf("Repeat: Iteration %d\n", iterCount)
		for _, stmt := range stmts {
			runOneStep(ctx, stmt.([]interface{}))
		}
	}
	return nil
}
