package main

import (
	"fmt"
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
	for _, scenarioName := range scenarios {
		if scenario, ok := script[scenarioName.(string)]; ok {
			runOneScenario(scenario.(map[string]interface{}), nil)
		} else {
			panic(fmt.Errorf("Scenario %s not found", scenarioName.(string)))
		}
	}
}

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

func runOneScenario(scenario map[string]interface{}, doneChan chan<- bool) {

	context := map[string]interface{}{}
	context["scenario_name"] = scenario["name"].(string)

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
	args := step[1:]
	if theFunc, ok := funcMap[method]; ok {
		err := theFunc(context, args)
		if err == nil {
			fmt.Printf("%s:%s succeeded!\n", myName, method)
		} else {
			fmt.Printf("%s:%s failed!: %s\n", myName, method, err.Error())
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
