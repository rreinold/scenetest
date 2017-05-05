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

type exitStmt struct{}
type sleepStmt struct{}
type repeatStmt struct{}
type forStmt struct{}
type ifStmt struct{}
type ifElseStmt struct{}
type printStmt struct{}
type setStmt struct{}
type setGlobalStmt struct{}
type assertStmt struct{}
type syncStmt struct{}
type syncAllStmt struct{}
type failStmt struct{}
type ignoreStmt struct{}
type breakStmt struct{}
type elemOfStmt struct{}
type setElemStmt struct{}
type concatStmt struct{}

func init() {
	funcMap["exit"] = &exitStmt{}
	funcMap["set"] = &setStmt{}
	funcMap["setGlobal"] = &setGlobalStmt{}
	funcMap["print"] = &printStmt{}
	funcMap["assert"] = &assertStmt{}
	funcMap["sleep"] = &sleepStmt{}
	funcMap["fail"] = &failStmt{}
	funcMap["ignore"] = &ignoreStmt{}
	funcMap["sync"] = &syncStmt{}
	funcMap["syncAll"] = &syncAllStmt{}
	funcMap["repeat"] = &repeatStmt{}
	funcMap["for"] = &forStmt{}
	funcMap["break"] = &breakStmt{}
	funcMap["if"] = &ifStmt{}
	funcMap["ifElse"] = &ifElseStmt{}
	funcMap["elemOf"] = &elemOfStmt{}
	funcMap["setElem"] = &setElemStmt{}
	funcMap["concat"] = &concatStmt{}
	syncLock = new(sync.Mutex)
}

func (s *syncStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 2, "", float64(1.11)); err != nil {
		return nil, err
	}
	syncKey := args[0].(string)
	syncCount := int(args[1].(float64))
	if err := doTheSync(syncKey, syncCount); err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *syncStmt) help() string {
	return "[\"sync\", <lockName>, count]"
}

func (s *syncAllStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 1, ""); err != nil {
		return nil, err
	}
	syncKey := args[0].(string)
	syncCount := TotalCount
	if err := doTheSync(syncKey, syncCount); err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *syncAllStmt) help() string {
	return "[\"syncAll\", <lockName>]"
}

func doTheSync(syncKey string, syncCount int) error {
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
		mySyncStuff.count = 0 // so we can use this sync point again in the test
	}
	syncLock.Unlock()
	<-mySyncStuff.c // wait for everybody to sync up
	return nil
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
		myPrintf("%+v", arg)
	}
	fmt.Println("")
	return nil, nil
}

func (p *printStmt) help() string {
	return "[\"print\", <arg>, ...]"
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
	if _, ok := args[0].(float64); !ok {
		return nil, fmt.Errorf("Repeat count must be a number")
	}
	if _, ok := args[1].([]interface{}); !ok {
		return nil, fmt.Errorf("Repeat statements must be an array")
	}
	stmts := args[1].([]interface{})
	iterCount := int(0)
	for count := int(args[0].(float64)); count > 0; count-- {
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
	incrNestingLevel(ctx)
	defer decrNestingLevel(ctx)
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
	iterCount := 0
	for {
		// Eval test condition, break out of loop if false
		val = valueOf(ctx, testCond)
		if !findTheTruth(val) {
			break
		}
		iterCount++
		if !ShutUp {
			myPrintf("For: iteration %d\n", iterCount)
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
	varName := args[0].(string)
	value := args[1]
	if isGlobal {
		setGlobal(varName, value)
	} else {
		ctx[varName] = value
	}
	return value, nil
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
		return nil, fmt.Errorf("fail statement takes one arg: %s", f.help())
	}
	switch args[0].(type) {
	case error:
		return args[0].(error).Error(), nil
	default:
		theTruth := findTheTruth(args[0])
		if !theTruth {
			return theTruth, nil
		}
		return nil, fmt.Errorf("fail statement succeeded -- very bad: %v", args[0])
	}
}

func (f *failStmt) help() string {
	return "[\"fail\", <stmt>]"
}

func (f *ignoreStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ignore statement takes one arg: %s", f.help())
	}
	return true, nil
}

func (f *ignoreStmt) help() string {
	return "[\"ignore\", <stmt>]"
}

func (e *elemOfStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("elemOf statement takes two args: %s", e.help())
	}
	switch args[0].(type) {
	case map[string]interface{}:
		theMap := args[0].(map[string]interface{})
		if key, ok := args[1].(string); ok {
			if val, ok := theMap[key]; ok {
				return val, nil
			} else {
				return nil, fmt.Errorf("Bad key to elemOf: %s", key)
			}
		} else {
			return nil, fmt.Errorf("elemOf: Map key must be a string")
		}
	case []interface{}:
		theSlice := args[0].([]interface{})
		if key, err := numberTypeOrFail(args[1]); err == nil {
			intKey := int(key)
			if intKey < 0 || intKey >= len(theSlice) {
				return nil, fmt.Errorf("Bad index in elemOf: %d", key)
			} else {
				return theSlice[intKey], nil
			}
		} else {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("elemOf statement takes either a slice or a map as the first arg\n")
	}
}

func (e *elemOfStmt) help() string {
	return "[\"elemOf\", <object>, <key>]"
}

func (e *setElemStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("setElemStmt statement takes three args: %s", e.help())
	}
	if err := argCheck(args, 3, "", nil, nil); err != nil {
		return nil, err
	}
	variable := args[0].(string)
	element := args[1]
	value := args[2]
	theVar := lookupVar(ctx, variable)
	if theVar == nil {
		return nil, fmt.Errorf("The variable '%s' does not exist", variable)
	}
	switch theVar.(type) {
	case map[string]interface{}:
		if _, ok := element.(string); !ok {
			return nil, fmt.Errorf("setElem: map key must be a string")
		}
		theMap := theVar.(map[string]interface{})
		index := element.(string)
		theMap[index] = value
	case []interface{}:
		if _, ok := element.(float64); !ok {
			return nil, fmt.Errorf("setElem: slice key must be a number")
		}
		theSlice := theVar.([]interface{})
		index := int(element.(float64))
		theSlice[index] = value
	default:
		return nil, fmt.Errorf("Type for variable '%s' must be a map or a slice", variable)
	}
	return theVar, nil
}

func (e *setElemStmt) help() string {
	return "[\"setElem\", \"<objName>\", \"<elemName>\", <value>]"
}

func (r *concatStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 2, "", float64(0)); err != nil {
		return nil, err
	}
	baseString := args[0].(string)
	baseInt := int(args[1].(float64))
	return fmt.Sprintf("%s%d", baseString, baseInt), nil
}

func (r *concatStmt) help() string {
	return "[\"concat\", \"<stringArg>\", <intArg>]"
}

func (e *exitStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 1, float64(0)); err != nil {
		return nil, err
	}
	exitStatus := int(args[0].(float64))
	return exitStatus, nil
}

func (e *exitStmt) help() string {
	return "[\"exit\", <exitStatusInt>]"
}
