package main

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"
)

type SyncStuff struct {
	count int
	c     chan bool
}

var syncMap = map[string]*SyncStuff{}
var syncLock *sync.Mutex

func init() {
	funcMap["sleep"] = &Statement{sleep, sleepHelp}
	funcMap["repeat"] = &Statement{repeat, doRepeatHelp}
	funcMap["while"] = &Statement{doWhile, doWhileHelp}
	funcMap["print"] = &Statement{doPrint, doPrintHelp}
	funcMap["set"] = &Statement{doSet, doSetHelp}
	funcMap["setGlobal"] = &Statement{doSetGlobal, doSetGlobalHelp}
	funcMap["incrGlobal"] = &Statement{incrGlobal, incrGlobalHelp}
	funcMap["decrGlobal"] = &Statement{decrGlobal, decrGlobalHelp}
	funcMap["assert"] = &Statement{doAssert, doAssertHelp}
	funcMap["sync"] = &Statement{doSync, doSyncHelp}
	syncLock = new(sync.Mutex)
}

func doSync(ctx map[string]interface{}, args []interface{}) error {
	if len(args) != 2 {
		return fmt.Errorf("Usage: %s", doSyncHelp())
	}
	syncKey := args[0].(string)
	syncCount := int(args[1].(float64))
	syncLock.Lock()
	fmt.Printf("IN: %v, %v\n", syncKey, syncCount)
	mySyncStuff, ok := syncMap[syncKey]
	if !ok {
		fmt.Printf("Making new one\n")
		mySyncStuff = &SyncStuff{0, make(chan bool, syncCount)}
		syncMap[syncKey] = mySyncStuff
	}
	mySyncStuff.count++
	if mySyncStuff.count >= syncCount {
		fmt.Printf("Unsyncing\n")
		for i := 0; i < syncCount; i++ {
			mySyncStuff.c <- true
		}
		mySyncStuff.count = 0
	}
	syncLock.Unlock()
	fmt.Printf("Waiting nicely\n")
	<-mySyncStuff.c
	return nil
}

func doSyncHelp() string {
	return "[\"sync\", <lockName>, count]"
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

func doWhile(ctx map[string]interface{}, args []interface{}) error {
	nestingLevel++
	defer decrNestingLevel()
	if len(args) != 4 {
		return fmt.Errorf("Usage: [while, <var>, <op>, <val>, [<stmt>...]]")
	}
	if _, ok := args[0].(string); !ok {
		return fmt.Errorf("First arg must be a variable name")
	}
	if _, ok := args[1].(string); !ok {
		return fmt.Errorf("Second arg must be a string value")
	}
	if _, ok := args[3].([]interface{}); !ok {
		return fmt.Errorf("While statements must be an array")
	}
	theVar := args[0].(string)
	theOp := args[1].(string)
	theVal := args[2]
	theStmts := args[3].([]interface{})
	iterCount := int(0)
	for {
		//  Test the condition and exit the loop if necessary
		theVarsVal := getGlobal(theVar)
		if !evaluateExpression(theOp, theVarsVal, theVal) {
			break
		}
		iterCount++
		fmt.Printf("While: iteration %d\n", iterCount)
		for _, stmt := range theStmts {
			runOneStep(ctx, stmt.([]interface{}))
		}
	}
	return nil
}

func evaluateExpression(op string, val1, val2 interface{}) bool {
	fmt.Printf("EVALUATE: %s, %+v, %+v\n", op, val1, val2)
	rval := false
	switch op {
	case "==":
		rval = makeNum(val1) == makeNum(val2)
	case "!=":
		rval = makeNum(val1) != makeNum(val2)
	case ">":
		rval = makeNum(val1) > makeNum(val2)
	case "<":
		rval = makeNum(val1) < makeNum(val2)
	case ">=":
		rval = makeNum(val1) >= makeNum(val2)
	case "<=":
		rval = makeNum(val1) <= makeNum(val2)
	}
	fmt.Printf("EVAL: %+v\n", rval)
	return rval
}

func makeNum(val interface{}) int {
	switch val.(type) {
	case string:
		theVal, err := strconv.Atoi(val.(string))
		if err != nil {
			fatal(err.Error())
		}
		return theVal
	case int:
		return val.(int)
	case uint8:
		return int(val.(uint8))
	case uint16:
		return int(val.(uint16))
	case uint32:
		return int(val.(uint32))
	case uint64:
		return int(val.(uint64))
	case int8:
		return int(val.(int8))
	case int16:
		return int(val.(int16))
	case int32:
		return int(val.(int32))
	case int64:
		return int(val.(int64))
	case float32:
		return int(val.(float32))
	case float64:
		return int(val.(float64))
	default:
		panic(fmt.Sprintf("Super-duper awful type passed to makeNum(): %+v", reflect.TypeOf(val)))
	}
	return -1 // Not reached
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

func doWhileHelp() string {
	return "[\"while\", <varName>, <varVal>, [<statement>, ...]]"
}

func doSet(ctx map[string]interface{}, args []interface{}) error {
	return doTheSet(ctx, args, false)
}

func doSetGlobal(ctx map[string]interface{}, args []interface{}) error {
	return doTheSet(ctx, args, true)
}

func doTheSet(ctx map[string]interface{}, args []interface{}, isGlobal bool) error {
	if err := argCheck(args, 2, "", nil); err != nil {
		return fmt.Errorf("Call to set failed: %s", err.Error())
	}
	varName := valueOf(ctx, args[0]).(string)
	value := valueOf(ctx, args[1])
	if isGlobal {
		setGlobal(varName, value)
	} else {
		ctx[varName] = value
	}
	ctx["returnValue"] = value
	return nil
}

func incrGlobal(ctx map[string]interface{}, args []interface{}) error {
	weInTheHouse()
	defer weOuttaTheHouse()
	if err := argCheck(args, 1, ""); err != nil {
		return fmt.Errorf("Call to incrGlobal failed: %s", err.Error())
	}
	varName := args[0].(string)
	setGlobal(varName, makeNum(getGlobal(varName))+1)
	return nil
}

func incrGlobalHelp() string {
	return "[\"incrGlobal\", <varName>]"
}

func decrGlobal(ctx map[string]interface{}, args []interface{}) error {
	weInTheHouse()
	defer weOuttaTheHouse()
	if err := argCheck(args, 1, ""); err != nil {
		return fmt.Errorf("Call to decrGlobal failed: %s", err.Error())
	}
	varName := args[0].(string)
	setGlobal(varName, makeNum(getGlobal(varName))-1)
	return nil
}

func decrGlobalHelp() string {
	return "[\"decrGlobal\", <varName>]"
}

func doSetHelp() string {
	return "[\"set\", \"<varName>\", <value> ]"
}

func doSetGlobalHelp() string {
	return "[\"setGlobal\", \"<varName>\", <value> ]"
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
