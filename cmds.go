package main

import (
	"flag"
	"fmt"
	cb "github.com/clearblade/Go-SDK"
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
	mustHaveAll("platformUrl", PlatformAddr, "messagingUrl", MsgAddr, "info", InfoFile)
	SetupFile = getFileOrDie()
	if StartNovi {
		startNovi()
	}
	performSetup(getJSON(SetupFile))
	if StartEdges {
		startEdges()
	}
}

func (s *runCommand) Run() {
	mustHaveAll("info", InfoFile)
	ScriptFile = getFileOrDie()
	scriptVars = getJSON(InfoFile)
	theScript := getJSON(ScriptFile)
	cb.CB_ADDR = scriptVars["platformUrl"].(string)
	cb.CB_MSG_ADDR = scriptVars["messagingUrl"].(string)
	if infoFileAndScriptFileDoNotMatch(scriptVars, theScript) {
	}
	executeTestScript(theScript)
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

func startNovi() {
	name := "clearblade"
	args := []string{"-tkey", "Oz49P0NCPD46Ojo6Og"}
	novi := procs.GetProcessManager(PlatformAddr, name, args)
	if novi == nil {
		fatal("Could not start novi process")
	}
	novi.Start()
	//log.Printf("NOVI PID: %d\n", novi.GetPid())
	setupState["noviPid"] = fmt.Sprintf("%d", novi.GetPid())
	scriptVars["noviPid"] = fmt.Sprintf("%d", novi.GetPid())
	time.Sleep(2 * time.Second)
}

func startEdges() {
}
