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

type sleepStmt struct{}
type repeatStmt struct{}
type forStmt struct{}
type whileStmt struct{}
type ifStmt struct{}
type ifElseStmt struct{}
type printStmt struct{}
type setStmt struct{}
type setGlobalStmt struct{}
type incrGlobalStmt struct{}
type decrGlobalStmt struct{}
type assertStmt struct{}
type syncStmt struct{}
type failStmt struct{}
type breakStmt struct{}

func init() {
	funcMap["sleep"] = &sleepStmt{}
	funcMap["repeat"] = &repeatStmt{}
	funcMap["for"] = &forStmt{}
	funcMap["while"] = &whileStmt{}
	funcMap["if"] = &ifStmt{}
	funcMap["if-else"] = &ifElseStmt{}
	funcMap["ifelse"] = &ifElseStmt{}
	funcMap["if else"] = &ifElseStmt{}
	funcMap["print"] = &printStmt{}
	funcMap["set"] = &setStmt{}
	funcMap["setGlobal"] = &setGlobalStmt{}
	funcMap["incrGlobal"] = &incrGlobalStmt{}
	funcMap["decrGlobal"] = &decrGlobalStmt{}
	funcMap["assert"] = &assertStmt{}
	funcMap["sync"] = &syncStmt{}
	funcMap["fail"] = &failStmt{}
	funcMap["break"] = &breakStmt{}
	syncLock = new(sync.Mutex)
}

func (s *syncStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("Usage: %s", s.help())
	}
	syncKey := args[0].(string)
	syncCount := int(args[1].(float64))
	syncLock.Lock()
	mySyncStuff, ok := syncMap[syncKey]
	if !ok {
		mySyncStuff = &SyncStuff{0, make(chan bool, syncCount)}
		syncMap[syncKey] = mySyncStuff
	}
	mySyncStuff.count++
	if mySyncStuff.count >= syncCount {
		for i := 0; i < syncCount; i++ {
			mySyncStuff.c <- true
		}
		mySyncStuff.count = 0
	}
	syncLock.Unlock()
	<-mySyncStuff.c
	return nil, nil
}

func (s *syncStmt) help() string {
	return "[\"sync\", <lockName>, count]"
}

func (s *sleepStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	secs := time.Duration(getArg(args, 0).(float64))
	time.Sleep(secs * time.Millisecond)
	return nil, nil
}

func (s *sleepStmt) help() string {
	return "[\"sleep\", <milliseconds>]"
}

func (p *printStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	for idx, arg := range args {
		if idx > 0 {
			myPrintf(" ")
		}
		myPrintf("%+v", valueOf(ctx, arg))
	}
	fmt.Println("")
	return nil, nil
}

func (p *printStmt) help() string {
	return "[\"print\", <arg>, ...]"
}

func (w *whileStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	incrNestingLevel(ctx)
	defer decrNestingLevel(ctx)
	if len(args) != 4 {
		return nil, fmt.Errorf("Usage: [while, <var>, <op>, <val>, [<stmt>...]]")
	}
	if _, ok := args[0].(string); !ok {
		return nil, fmt.Errorf("First arg must be a variable name")
	}
	if _, ok := args[1].(string); !ok {
		return nil, fmt.Errorf("Second arg must be a string value")
	}
	if _, ok := args[3].([]interface{}); !ok {
		return nil, fmt.Errorf("While statements must be an array")
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
		myPrintf("While: iteration %d\n", iterCount)
		for _, stmt := range theStmts {
			if _, err := runOneStep(ctx, stmt.([]interface{})); err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}

func evaluateExpression(op string, val1, val2 interface{}) bool {
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

func (r *repeatStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	incrNestingLevel(ctx)
	defer decrNestingLevel(ctx)
	if len(args) != 2 {
		return nil, fmt.Errorf("Usage: [repeat, <num>, [<stmt>...]]")
	}
	if _, ok := valueOf(ctx, args[0]).(float64); !ok {
		return nil, fmt.Errorf("Repeat count must be a number")
	}
	if _, ok := valueOf(ctx, args[1]).([]interface{}); !ok {
		return nil, fmt.Errorf("Repeat statements must be an array")
	}
	stmts := valueOf(ctx, args[1]).([]interface{})
	iterCount := int(0)
	for count := int(valueOf(ctx, args[0]).(float64)); count > 0; count-- {
		iterCount++
		if !ShutUp {
			myPrintf("Repeat: iteration %d\n", iterCount)
		}
		for _, stmt := range stmts {
			if _, err := runOneStep(ctx, stmt.([]interface{})); err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}

func (r *repeatStmt) help() string {
	return "[\"repeat\", <count>, [<statement>, ...]]"
}

///////////////////////////////// "for" statement

func (f *forStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 2, []interface{}{}, []interface{}{}); err != nil {
		return nil, fmt.Errorf("Bad arguments to for statement: %s", err.Error())
	}
	condList := args[0].([]interface{})
	stmtList := args[1].([]interface{})
	if len(condList) != 3 {
		return nil, fmt.Errorf("Must be three conditions in for stmt condition list")
	}
	preCond, testCond, postCond := condList[0], condList[1], condList[2]

	//  Do the pre condition just once. We don't care about the returned
	//  value; just that it succeeded
	valueOf(ctx, preCond)

	var val interface{}
	var err error
	for {
		// Eval test condition, break out of loop if false
		val = valueOf(ctx, testCond)
		if !findTheTruth(val) {
			break
		}

		//  true test condition -- now execute each sub statement
		for _, stmt := range stmtList {
			if val, err = runOneStep(ctx, stmt.([]interface{})); err != nil {
				// If there's a break here, catch it here and return success
				// note that we're not executing the post statement/condition
				if err.Error() == "break" {
					return nil, nil
				}
				return nil, err
			}
		}

		//  Finally, do the post condition
		val = valueOf(ctx, postCond)
	}
	return val, nil
}

func (f *forStmt) help() string {
	return "[\"for\", [<initStmt>, <condStmt>, <endLoopStmt>] , [<statement>, ...]]"
}

///////////////////////////////// "if" statement

func (i *ifStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("if statement: wrong number of args: %s", i.help())
	}
	condition := findTheTruth(args[0])
	stmtList := args[1].([]interface{})

	// eval the condition expression and return if error or false
	/*
		res, err := evalExprStmt(ctx, condition)
		if err != nil {
			return nil, err
		}
		if res == false {
			return nil, nil
		}
	*/
	if condition == false {
		return nil, nil
	}

	// The condition returned true. Execute the list of statements.
	for _, stmt := range stmtList {
		_, err := runOneStep(ctx, stmt.([]interface{}))
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i *ifStmt) help() string {
	return "[\"if\" [<exprStnt>], [<stmt>, <stmt>...]]"
}

func (i *ifElseStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("if statement: wrong number of args: %s", i.help())
	}
	condition := args[0].([]interface{})
	ifList := args[1].([]interface{})
	elseList := args[2].([]interface{})

	// eval the condition expression and return if error or false
	res, err := evalExprStmt(ctx, condition)
	if err != nil {
		return nil, err
	}
	var stmtList []interface{}
	if res {
		stmtList = ifList
	} else {
		stmtList = elseList
	}

	// The condition returned true. Execute the list of statements.
	for _, stmt := range stmtList {
		_, err := runOneStep(ctx, stmt.([]interface{}))
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i *ifElseStmt) help() string {
	return "[\"if else\" [<exprStnt>], [<if-stmt>, <if stmt>...], [<else-stmt>, <else-stmt>...]"
}

func (w *whileStmt) help() string {
	return "[\"while\", <varName>, <varVal>, [<statement>, ...]]"
}

func (s *setStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	return doTheSet(ctx, args, false)
}

func (s *setGlobalStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	return doTheSet(ctx, args, true)
}

func doTheSet(ctx map[string]interface{}, args []interface{}, isGlobal bool) (interface{}, error) {
	if err := argCheck(args, 2, "", nil); err != nil {
		return nil, fmt.Errorf("Call to set failed: %s", err.Error())
	}
	varName := valueOf(ctx, args[0]).(string)
	value := valueOf(ctx, args[1])
	if isGlobal {
		setGlobal(varName, value)
	} else {
		ctx[varName] = value
	}
	return value, nil
}

func (i *incrGlobalStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	weInTheHouse()
	defer weOuttaTheHouse()
	if err := argCheck(args, 1, ""); err != nil {
		return nil, fmt.Errorf("Call to incrGlobal failed: %s", err.Error())
	}
	varName := args[0].(string)
	setGlobal(varName, makeNum(getGlobal(varName))+1)
	return nil, nil
}

func (i *incrGlobalStmt) help() string {
	return "[\"incrGlobal\", <varName>]"
}

func (d *decrGlobalStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	weInTheHouse()
	defer weOuttaTheHouse()
	if err := argCheck(args, 1, ""); err != nil {
		return nil, fmt.Errorf("Call to decrGlobal failed: %s", err.Error())
	}
	varName := args[0].(string)
	setGlobal(varName, makeNum(getGlobal(varName))-1)
	return nil, nil
}

func (d *decrGlobalStmt) help() string {
	return "[\"decrGlobal\", <varName>]"
}

func (s *setStmt) help() string {
	return "[\"set\", \"<varName>\", <value> ]"
}

func (s *setGlobalStmt) help() string {
	return "[\"setGlobal\", \"<varName>\", <value> ]"
}

func (a *assertStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Usage: [assert, <expr> ]")
	}
	myPrintf("ARGS0 == %v\n", args[0])
	theTruth := findTheTruth(args[0])
	if !theTruth {
		return nil, fmt.Errorf("Assertion failed: %v", args[0])
	}
	return args[0], nil
}

var assertHelpStr = `["assert", <expr>]`

func (a *assertStmt) help() string {
	return assertHelpStr
}

func (b *breakStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("break statement takes no args")
	}
	//  This is a bit of a hack, but return the break error up the stack -- it gets "caught"
	//  by all the looping statements. If it gets all the way out to the outer scope/nesting level
	//  then the error is reported.
	return nil, fmt.Errorf("break")
}

func (b *breakStmt) help() string {
	return "[\"break\"]"
}

func (f *failStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("fail statement takes one statement arg: %s", f.help())
	}
	if _, ok := args[0].([]interface{}); !ok {
		return nil, fmt.Errorf("fail statement expects one statement (slice) arg")
	}
	rval, err := runOneStep(ctx, args[0].([]interface{}))
	if err == nil {
		return nil, fmt.Errorf("fail statement succeeded -- very bad: %v", rval)
	}
	return err, nil
}

func (f *failStmt) help() string {
	return "[\"fail\", <stmt>]"
}
