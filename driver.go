package main

import (
	"fmt"
	"time"
)

var script map[string]interface{}

func init() {
}

func executeTestScript(theScript map[string]interface{}) {
	script = theScript
	sequencing := getVar("sequencing", script, "Parallel").(string)
	scenarios := getVar("scenarios", script, []string{}).([]interface{})
	if sequencing == "Serial" {
		runSerial(scenarios)
	} else if sequencing == "Parallel" {
		runParallel(scenarios)
	} else {
		panic(fmt.Errorf("Bad sequencing: %s\n", sequencing))
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
		scenario := duh.(map[string]interface{})
		for i := 0; i < count; i++ {
			nm := fmt.Sprintf("%s(%d)", scenario["name"].(string), i)
			runOneScenario(nm, scenario, nil)
		}
	}
}

/*
func runSerial(scenarios []interface{}) {
	for _, scenarioName := range scenarios {
		if scenario, ok := script[scenarioName.(string)]; ok {
			runOneScenario(scenario.(map[string]interface{}), nil)
		} else {
			panic(fmt.Errorf("Scenario %s not found", scenarioName.(string)))
		}
	}
}
*/

/*
func runParallel(scenarios []interface{}) {
	doneChan := make(chan bool)
	for _, scenarioName := range scenarios {
		if scenario, ok := script[scenarioName.(string)]; ok {
			go runOneScenario(scenario.(map[string]interface{}), doneChan)
		} else {
			panic(fmt.Errorf("Scenario %s not found", scenarioName.(string)))
		}
	}

	for i := 0; i < len(scenarios); i++ {
		<-doneChan
	}
}
*/

func runParallel(scenarios []interface{}) {
	totalCount := 0
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
		scenario := duh.(map[string]interface{})
		totalCount += count
		for i := 0; i < count; i++ {
			nm := fmt.Sprintf("%s(%d)", scenario["name"].(string), i)
			go runOneScenario(nm, scenario, doneChan)
		}

	}
	for ; totalCount > 0; totalCount-- {
		<-doneChan
	}
}

func runOneScenario(name string, scenario map[string]interface{}, doneChan chan<- bool) {

	context := map[string]interface{}{}
	context["scenario_name"] = name

	steps := getVar("steps", scenario, [][]interface{}{}).([]interface{})

	for _, step := range steps {
		runOneStep(context, step.([]interface{}))
	}

	//  If we're in parallel, need to tell parent we're done.
	if doneChan != nil {
		doneChan <- true
	}
}

func runOneStep(context map[string]interface{}, step []interface{}) {
	myName := context["scenario_name"].(string)
	if len(step) == 0 {
		return
	}
	method := step[0].(string)
	args := dereferenceVariables(context, step[1:])
	if theStmt, ok := funcMap[method]; ok {
		err := theStmt.RunFunc(context, args)
		timeStr := time.Now().Format(time.UnixDate)
		if err == nil {
			fmt.Printf("%s(%s):%s succeeded!\n", myName, timeStr, method)
		} else {
			fmt.Printf("%s(%s):%s failed!: %s\n", myName, timeStr, method, err.Error())
			fatal("Exiting because of error")
		}
	} else {
		fmt.Printf("Unknown function %s\n", method)
	}
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
	case string:
		return scenario.(string), 1, nil
	case []interface{}:
		theScenario := scenario.([]interface{})
		if len(theScenario) != 2 {
			return "", 0, fmt.Errorf("An array scenario must be of the form [\"name\", <count>]")
		}
		return theScenario[0].(string), int(theScenario[1].(float64)), nil
	default:
		return "", 0, fmt.Errorf("Bad scenario type")
	}
}
