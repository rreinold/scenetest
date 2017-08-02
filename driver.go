package main

import (
	"fmt"
	cb "github.com/clearblade/Go-SDK"
	"os"
	"sync"
	"time"
)

var script map[string]interface{}
var adminClient *cb.DevClient
var adminClients = map[string]*cb.DevClient{}
var adminClientsLock = sync.RWMutex{}

func executeTestScript(theScript map[string]interface{}) {
	script = theScript
	if !NoLogin { // nice double negative
		adminClient = authDevForScriptRun()
	}

	//  Get the scenarios and dereference any variables appearing
	//  in the scenarios. These are typically scenario counts where
	//  the number of scenarios to run is set in globals
	sequencing := getVar("sequencing", script, "Parallel").(string)
	scenarios := getVar("scenarios", script, []string{}).([]interface{})
	if glbs, ok := script["globals"].(map[string]interface{}); ok {
		globals = glbs
	}
	if sequencing == "Serial" {
		panic("Cannot user Serial as a run mode -- discontinued. :-)")
	}
	runParallel(scenarios)
}

func authDevForScriptRun() *cb.DevClient {
	return authDevWithAddrs(cb.CB_ADDR, cb.CB_MSG_ADDR)
}

func authDevWithAddrs(httpAddr, mqttAddr string) *cb.DevClient {
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	adminClientsLock.Lock()
	defer adminClientsLock.Unlock()
	if cli, ok := adminClients[httpAddr]; ok {
		return cli
	}
	theDev := scriptVars["developer"].(map[string]interface{})
	email, password := theDev["email"].(string), theDev["password"].(string)
	ac := cb.NewDevClientWithAddrs(httpAddr, mqttAddr, email, password)
	if err := ac.Authenticate(); err != nil {
		fatal("Could not authenticate developer: " + err.Error())
	}
	adminClients[httpAddr] = ac
	return ac
}

//
//  This little gem just allows us to have two different models for a scenario: The old one
//  I don't like:
//
//		"scenarioName": {
//			"name": "scenarioName",
//			"steps": [
//				<steps>
//			]
//		}
//
//  and a stripped down new one:
//
//  	"scenarioName": [
//  		<steps>
//  	]
//
//  It just "maps" the stripped down scenario into what we're used to. Maintain backward
//  compatibility and provide a streamlined new syntax.
//
func getNameAndSteps(outerName string, scenario interface{}) (string, map[string]interface{}) {
	switch scenario.(type) {
	case map[string]interface{}:
		mapScen := scenario.(map[string]interface{})
		return mapScen["name"].(string), mapScen
	case []interface{}:
		return outerName, map[string]interface{}{"name": outerName, "steps": scenario}
	default:
		fatal(fmt.Sprintf("Poorly defined scenario: %#v", scenario))
		return "", nil
	}
}

func runSerial(scenarios []interface{}) {
	for _, scenarioSpec := range scenarios {
		name, count, err := parseScenario(scenarioSpec)
		if err != nil {
			fatal(err.Error())
		}
		duh, ok := script[name]
		if !ok {
			panic(fmt.Errorf("Scenario %s not found", name))
		}
		name, scenario := getNameAndSteps(name, duh)
		//scenario := duh.(map[string]interface{})
		for i := 0; i < count; i++ {
			nm := fmt.Sprintf("%s(%d)", name, i+1)
			runOneScenario(nm, scenario, nil, nil, i)
		}
	}
}

func runParallel(scenarios []interface{}) {
	TotalCount = 0
	startChan := make(chan bool)
	doneChan := make(chan bool)
	for _, scenarioSpec := range scenarios {
		name, count, err := parseScenario(scenarioSpec)
		if err != nil {
			fatal(err.Error())
		}
		duh, ok := script[name]
		if !ok {
			fatal(fmt.Sprintf("Couldn't find script %s", name))
		}
		name, scenario := getNameAndSteps(name, duh)
		TotalCount += count
		for i := 0; i < count; i++ {
			nm := fmt.Sprintf("%s(%d)", name, i+1)
			go runOneScenario(nm, scenario, startChan, doneChan, i)
		}

	}

	for i := 0; i < TotalCount; i++ {
		startChan <- true
	}
	for ; TotalCount > 0; TotalCount-- {
		<-doneChan
	}
	close(startChan)
	close(doneChan)
}

func putVarsInContext(context map[string]interface{}) map[string]interface{} {
	locals, ok := script["locals"]
	if !ok {
		return context
	}
	for key, val := range locals.(map[string]interface{}) {
		context[key] = copyVariable(val)
	}
	return context
}

func runOneScenario(name string, scenario map[string]interface{}, startChan <-chan bool, doneChan chan<- bool, myInstance int) {

	context := map[string]interface{}{}
	context["scenario_name"] = name
	context["adminClient"] = adminClient
	context["__nestingLevel"] = int(0)
	context["myInstance"] = myInstance
	context = putVarsInContext(context)

	steps := getVar("steps", scenario, [][]interface{}{}).([]interface{})

	// We wait to be notified on startChan because we want all scenarios created
	// before any of them start. This is primarily for the ["syncall", <syncName>]
	// command. One of the biggest mistakes I (swm) make is miscalculating the
	// number of scenarios in a test. A common practice is to wait until all scenarios
	// get to a certain point before proceeding.
	if startChan != nil {
		<-startChan
	}

	for _, step := range steps {
		sliceStep := step.([]interface{})
		if len(sliceStep) == 0 {
			fmt.Printf("Skipping empty step\n")
			continue
		}
		_, err := runOneStep(context, sliceStep)
		if err != nil {
			if err.Error() == "break" {
				err = fmt.Errorf("Encountered break statement outside of loop")
			}
			fatal(fmt.Sprintf("Exiting because of error: %s", err.Error()))
		}
	}

	//  If we're in parallel, need to tell parent we're done.
	if doneChan != nil {
		doneChan <- true
	}
}

func runOneStep(context map[string]interface{}, step []interface{}) (interface{}, error) {
	myName := context["scenario_name"].(string)
	if len(step) == 0 {
		return nil, nil
	}
	method := step[0].(string)
	args := dereferenceVariables(context, step[1:])
	if theStmt, ok := funcMap[method]; ok {
		rval, err := theStmt.run(context, args)
		timeStr := time.Now().Format(time.UnixDate)
		if err == nil {
			context["returnValue"] = rval
			if !ShutUp {
				myNestingPrintf(context, "%s:\t%s:\t%s succeeded\n", timeStr, myName, method)
			}
			if method == "exit" {
				os.Exit(rval.(int))
			}
			return rval, nil
		} else if err.Error() == "break" {
			if !ShutUp {
				myNestingPrintf(context, "%s:\t%s:\t%s succeeded\n", timeStr, myName, method)
			}
			return nil, err
		} else {
			myNestingPrintf(context, "%s(%s):%s failed: %s\n", myName, timeStr, method, err.Error())
			return nil, err
		}
	}
	return nil, fmt.Errorf("Unknown statement: %s\n", method)
}

func getVar(name string, script map[string]interface{}, defaultVal interface{}) interface{} {
	if val, ok := script[name]; ok {
		return val
	}
	return defaultVal
}

func dereferenceVariables(context map[string]interface{}, args []interface{}) []interface{} {
	rval := []interface{}{}
	for _, arg := range args {
		rval = append(rval, valueOf(context, arg))
	}
	return rval
}

func parseScenario(scenario interface{}) (string, int, error) {
	switch scenario.(type) {
	case []interface{}:
		theScenario := scenario.([]interface{})
		if len(theScenario) != 2 {
			return "", 0, fmt.Errorf("A scenario must be of the form [\"name\", <count>]")
		}
		parsedScenario := dereferenceVariables(map[string]interface{}{}, theScenario)
		return parsedScenario[0].(string), int(parsedScenario[1].(float64)), nil
	default:
		return "", 0, fmt.Errorf("Bad scenario type")
	}
}

func copyVariable(variable interface{}) interface{} {
	switch variable.(type) {
	case map[string]interface{}:
		mapCopy := map[string]interface{}{}
		for key, val := range variable.(map[string]interface{}) {
			mapCopy[key] = val
		}
		return mapCopy
	case []interface{}:
		sliceCopy := make([]interface{}, len(variable.([]interface{})))
		for idx, val := range variable.([]interface{}) {
			sliceCopy[idx] = val
		}
		return sliceCopy
	default:
		return variable
	}
}
