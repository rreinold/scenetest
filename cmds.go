package main

import (
	"encoding/json"
	"flag"
	"fmt"
	cb "github.com/clearblade/Go-SDK"
	"github.com/clearblade/cbjson"
	procs "github.com/clearblade/scenetest/processes"
	"time"
)

type setupCommand struct{}
type runCommand struct{}
type teardownCommand struct{}
type helpCommand struct{}

var (
	commandMap = map[string]ScenetestCmd{}
)

func init() {
	commandMap["setup"] = &setupCommand{}
	commandMap["run"] = &runCommand{}
	commandMap["teardown"] = &teardownCommand{}
	commandMap["help"] = &helpCommand{}
}

func getCommand(name string) ScenetestCmd {
	theCmd, ok := commandMap[name]
	if !ok {
		fatalf("Could not find command '%s'\n", name)
	}
	return theCmd
}

func (s *setupCommand) Run() {
	mustHaveAll("platformUrl", PlatformAddr, "messagingUrl", MsgAddr)

	args := flag.Args()

	switch(len(args)){
		case 0:
			SetupFile = DEFAULT_SETUP_FILE
			break
		case 1:
			SetupFile = args[0]
			break
		default:
			goodbye(fmt.Errorf("Unexpected number of file arguments"))
	}
	if StartPlatform {
		startPlatform()
	}
	performSetup(getJSON(SetupFile))
	if StartEdges {
		startEdges()
	}
}

func (s *runCommand) Run() {

	args := flag.Args()

	switch(len(args)){
		case 0:
			ScriptFile = DEFAULT_RUN_FILE
			break
		case 1:
			ScriptFile = args[0]
		default:
			goodbye(fmt.Errorf("Unexpected number of file arguments"))
	}

	scriptVars = getInfoFile(InfoFile)
	theScript := overrideGlobalsAndLocals(getJSON(ScriptFile))
	cb.CB_ADDR = scriptVars["platformUrl"].(string)
	cb.CB_MSG_ADDR = scriptVars["messagingUrl"].(string)
	if PlatformAddr != "" {
		cb.CB_ADDR = PlatformAddr
	}
	if MsgAddr != "" {
		cb.CB_MSG_ADDR = MsgAddr
	}
	if infoFileAndScriptFileDoNotMatch(scriptVars, theScript) {
	}
	executeTestScript(theScript)
}

func overrideGlobalsAndLocals(testScript map[string]interface{}) map[string]interface{} {
	//
	// there are three flags involved -- overrides, globals, and locals. overrides
	// overrides globals and locals. If that makes any sense.
	//
	varDict := map[string]interface{}{}
	var ok bool
	if Overrides != "" {

		overrides, _, err := cbjson.GetJSONFile(Overrides)
		if err != nil {
			fatalf("Could not process '-overrides' file: %s\n", err)
		}

		if varDict, ok = overrides["globals"].(map[string]interface{}); ok {
			testScript = putVarsIntoDict("globals", testScript, varDict)
		}
		if varDict, ok = overrides["locals"].(map[string]interface{}); ok {
			testScript = putVarsIntoDict("locals", testScript, varDict)
		}
		return testScript
	}
	if Globals != "" {
		if err := json.Unmarshal([]byte(Globals), &varDict); err != nil {
			fatalf("Malformed 'globals' argument: %s\n", err)
		}
		testScript = putVarsIntoDict("globals", testScript, varDict)
	}
	if Locals != "" {
		if err := json.Unmarshal([]byte(Locals), &varDict); err != nil {
			fatalf("Malformed 'locals' argument: %s\n", err)
		}
		testScript = putVarsIntoDict("locals", testScript, varDict)
	}
	return testScript
}

func putVarsIntoDict(varName string, target, source map[string]interface{}) map[string]interface{} {
	if _, has := target[varName]; !has {
		target[varName] = map[string]interface{}{}
	}
	targetVars := target[varName].(map[string]interface{})
	for key, val := range source {
		targetVars[key] = val
	}
	target[varName] = targetVars
	return target
}

func infoFileAndScriptFileDoNotMatch(testInfo, testScript map[string]interface{}) bool {
	/*
		infoFileName, ok := testInfo["name"].(string)
		if !ok {
			fatal("info file does not have the \"name\" key/value pair")
		}
		scriptName, ok := testScript["systemName"].(string)
		if !ok {
			fatal("test script file does not have the \"systemName\" key/value pair")
		}
		return infoFileName == scriptName
	*/
	return false
}

func (s *teardownCommand) Run() {
	mustHaveAll("info", InfoFile)
	scriptVars = getJSON(InfoFile)
	cb.CB_ADDR = scriptVars["platformUrl"].(string)
	cb.CB_MSG_ADDR = scriptVars["messagingUrl"].(string)
	performTeardown()
}

func (s *helpCommand) Run() {
	showHelp(flag.Args())
}

func startPlatform() {
	name := "clearblade"
	args := []string{"-tkey", "Oz49P0NCPD46Ojo6Og"}
	platform := procs.GetProcessManager(PlatformAddr, name, args)
	if platform == nil {
		fatal("Could not start platform process")
	}
	platform.Start()
	//log.Printf("platform PID: %d\n", platform.GetPid())
	setupState["platformPid"] = fmt.Sprintf("%d", platform.GetPid())
	scriptVars["platformPid"] = fmt.Sprintf("%d", platform.GetPid())
	time.Sleep(2 * time.Second)
}

func startEdges() {
}
