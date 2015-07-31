package main

import (
	"fmt"
	"reflect"
	"time"
)

func init() {
	funcMap["sleep"] = &Statement{sleep, sleepHelp}
	funcMap["repeat"] = &Statement{repeat, doRepeatHelp}
	funcMap["print"] = &Statement{doPrint, doPrintHelp}
	funcMap["set"] = &Statement{doSet, doSetHelp}
	funcMap["assert"] = &Statement{doAssert, doAssertHelp}
}

func sleep(ctx map[string]interface{}, args []interface{}) error {
	secs := time.Duration(getArg(args, 0).(float64))
	time.Sleep(secs * time.Millisecond)
	return nil
}

func sleepHelp() string {
	return "[\"sleep\", <milliseconds>]"
}

func doPrint(ctx map[string]interface{}, args []interface{}) error {
	for idx, arg := range args {
		if idx > 0 {
			fmt.Printf(" ")
		}
		fmt.Printf("%+v", valueOf(ctx, arg))
	}
	fmt.Println("")
	return nil
}

func doPrintHelp() string {
	return "[\"print\", <arg>, ...]"
}

func decrNestingLevel() {
	nestingLevel--
}

func repeat(ctx map[string]interface{}, args []interface{}) error {
	nestingLevel++
	defer decrNestingLevel()
	if len(args) != 2 {
		return fmt.Errorf("Usage: [repeat, <num>, [<stmt>...]]")
	}
	if _, ok := valueOf(ctx, args[0]).(float64); !ok {
		return fmt.Errorf("Repeat count must be a number")
	}
	if _, ok := valueOf(ctx, args[1]).([]interface{}); !ok {
		return fmt.Errorf("Repeat statements must be an array")
	}
	stmts := valueOf(ctx, args[1]).([]interface{})
	iterCount := int(0)
	for count := int(valueOf(ctx, args[0]).(float64)); count > 0; count-- {
		iterCount++
		fmt.Printf("Repeat: iteration %d\n", iterCount)
		for _, stmt := range stmts {
			runOneStep(ctx, stmt.([]interface{}))
		}
	}
	return nil
}

func doRepeatHelp() string {
	return "[\"repeat\", <count>, [<statement>, ...]]"
}

func doSet(ctx map[string]interface{}, args []interface{}) error {
	if err := argCheck(args, 2, "", nil); err != nil {
		return fmt.Errorf("Call to set failed: %s", err.Error())
	}

	varName := valueOf(ctx, args[0]).(string)
	value := valueOf(ctx, args[1])
	ctx[varName] = value
	ctx["returnValue"] = value
	return nil
}

func doSetHelp() string {
	return "[\"set\", \"<varName>\", <value> ]"
}

func doAssert(ctx map[string]interface{}, args []interface{}) error {
	if len(args) != 2 {
		return fmt.Errorf("Usage: [assert, <val1>, <val2> ]")
	}
	result := reflect.DeepEqual(args[0], args[1])
	ctx["returnValue"] = result
	if !result {
		return fmt.Errorf("Assertion failed: %+v and %+v are not equal", args[0], args[1])
	}
	return nil
}

var assertHelpStr = `["assert", <val1>, <val2>] or ["assert", <val1>, "<operator>", <val2>]`

func doAssertHelp() string {
	return assertHelpStr
}
