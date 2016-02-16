package main

import (
	//"bytes"
	"cbjson"
	//"encoding/json"
	"flag"
	"fmt"
	cb "github.com/clearblade/Go-SDK"
	//"io/ioutil"
	"os"
	"sync"
)

type Stmt interface {
	run(ctx map[string]interface{}, args []interface{}) (interface{}, error)
	help() string
}

const (
	SceneTestEnvVar = "SCENETEST_PATH"
)

var (
	MsgAddr        string
	PlatformAddr   string
	ScriptFile     string
	SetupFile      string
	TeardownFile   string
	InfoFile       string
	JustParse      bool
	GetSomeHelp    bool
	SceneRoot      string
	FileSearchPath []string
	Login          bool
	ShutUp         bool
)

var (
	funcMap    = map[string]Stmt{}
	scriptVars = map[string]interface{}{}
	globals    = map[string]interface{}{}
	globalLock = sync.RWMutex{}
	printLock  = sync.Mutex{}
)

func init() {
	flag.StringVar(&MsgAddr, "messaging-url", "undefined", "Msg service location")
	flag.StringVar(&PlatformAddr, "platform-url", "undefined", "Platform location")
	flag.StringVar(&ScriptFile, "run", "Do Not Run Script", "Script file to execute")
	flag.StringVar(&SetupFile, "setup", "Do Not Setup", "File to setup system(s)")
	flag.StringVar(&TeardownFile, "teardown", "Do Not Teardown", "File that tears you up")
	flag.StringVar(&InfoFile, "info", "Do Not Get Info", "File generated by setup to be used by the tests")
	flag.BoolVar(&GetSomeHelp, "help", false, "Print out a help string for all statements")
	flag.BoolVar(&JustParse, "parse", false, "Just parse everything and don't execute")
	flag.BoolVar(&Login, "login", true, "login to the clearblade system as specified by cmd line args to TestSystemInfo.json")
	flag.BoolVar(&ShutUp, "silent", false, "Shut Up!j")
	scriptVars["roles"] = map[string]interface{}{}
	scriptVars["users"] = map[string]interface{}{}
	scriptVars["collections"] = map[string]interface{}{}
	scriptVars["items"] = map[string]interface{}{}
	scriptVars["codeServices"] = map[string]interface{}{}
	scriptVars["codeLibraries"] = map[string]interface{}{}
	scriptVars["triggers"] = map[string]interface{}{}
	scriptVars["timers"] = map[string]interface{}{}
}

func main() {
	flag.Parse()
	cb.CB_ADDR = PlatformAddr
	cb.CB_MSG_ADDR = MsgAddr
	setupSceneRoot()
	setupFileSearchPath()
	if GetSomeHelp {
		showHelp()
		return
	}
	if JustParse {
		parseProvidedFiles()
		return
	}
	if SetupFile != "Do Not Setup" && InfoFile != "Do Not Get Info" {
		fatal("Can't have both a setup file and an info file. I know. Confusing.")
	}
	if SetupFile != "Do Not Setup" {
		performSetup(getJSON(SetupFile))
	}
	if InfoFile != "Do Not Get Info" {
		scriptVars = getJSON(InfoFile)
		if PlatformAddr == "undefined" {
			cb.CB_ADDR = scriptVars["platformUrl"].(string)
		}
		if MsgAddr == "undefined" {
			cb.CB_MSG_ADDR = scriptVars["messagingUrl"].(string)
		}
	}
	if ScriptFile != "Do Not Run Script" {
		executeTestScript(getJSON(ScriptFile))
	}
	if TeardownFile != "Do Not Teardown" {
		performTeardown(getJSON(TeardownFile))
	}
}

func getJSON(filename string) map[string]interface{} {
	theStuff, _, err := cbjson.GetJSONFile(filename)
	if err != nil {
		goodbye(err)
	}
	return theStuff
}

func goodbye(err error) {
	myPrintf("%s\n", err.Error())
	os.Exit(1)
}

func parseProvidedFiles() {
	if SetupFile != "Do Not Setup" {
		myPrintf("Parsing %s... ", SetupFile)
		getJSON(SetupFile)
		myPrintf("ok\n")
	}
	if TeardownFile != "Do Not Teardown" {
		myPrintf("Parsing %s... ", TeardownFile)
		getJSON(TeardownFile)
		myPrintf("ok\n")
	}
	if ScriptFile != "Do Not Run Script" {
		myPrintf("Parsing %s... ", ScriptFile)
		getJSON(ScriptFile)
		myPrintf("ok\n")
	}
}

func weInTheHouse() {
	globalLock.Lock()
}

func weOuttaTheHouse() {
	globalLock.Unlock()
}

func getGlobal(name string) interface{} {
	if val, ok := globals[name]; ok {
		return val
	}
	return nil
}

func setGlobal(name string, val interface{}) {
	globals[name] = val
}

func fatalf(theFmt string, args ...interface{}) {
	myPrintf(theFmt, args...)
	os.Exit(1)
}

func myPrintf(theFmt string, args ...interface{}) {
	printLock.Lock()
	defer printLock.Unlock()
	fmt.Printf(theFmt, args...)
}

func myNestingPrintf(ctx map[string]interface{}, theFmt string, args ...interface{}) {
	duhFmt := ""
	lvl := ctx["__nestingLevel"].(int)
	for i := 0; i < lvl; i++ {
		duhFmt += "    "
	}
	myPrintf(duhFmt+theFmt, args...)
}
