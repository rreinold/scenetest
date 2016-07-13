package main

import (
	"flag"
	"fmt"
	cb "github.com/clearblade/Go-SDK"
	"github.com/clearblade/cbjson"
	"os"
	"strings"
	"sync"
)

type Stmt interface {
	run(ctx map[string]interface{}, args []interface{}) (interface{}, error)
	help() string
}

type Edge struct {
	name     string
	ipAddr   string
	httpPort string
	mqttPort string
}

const (
	SceneTestEnvVar = "SCENETEST_PATH"
)

var (
	MsgAddr      string
	PlatformAddr string
	ScriptFile   string
	SetupFile    string
	TeardownFile string
	InfoFile     string
	//JustParse      bool
	GetSomeHelp    bool
	SceneRoot      string
	FileSearchPath []string
	NoLogin        bool
	ShutUp         bool
	EdgeInfo       string
)

var (
	funcMap        = map[string]Stmt{}
	scriptVars     = map[string]interface{}{}
	globals        = map[string]interface{}{}
	edgeMap        = map[string]*Edge{}
	globalLock     = sync.RWMutex{}
	printLock      = sync.Mutex{}
	scriptVarsLock = sync.RWMutex{}
)

func init() {
	flag.StringVar(&MsgAddr, "messaging-url", "", "Msg service location")
	flag.StringVar(&PlatformAddr, "platform-url", "", "Platform location")
	flag.StringVar(&InfoFile, "info", "", "File generated by setup to be used by the tests")
	flag.BoolVar(&GetSomeHelp, "help", false, "Print out a help string for all statements")
	flag.BoolVar(&NoLogin, "nologin", false, "login to the clearblade system as specified by cmd line args to TestSystemInfo.json")
	flag.BoolVar(&ShutUp, "silent", false, "Shut Up!")
	flag.StringVar(&EdgeInfo, "edge-info", "", "name|ip|httpport|mqttport,name|ip|httpport|mqttport -- edge config")

	scriptVars["roles"] = map[string]interface{}{}
	scriptVars["users"] = map[string]interface{}{}
	scriptVars["collections"] = map[string]interface{}{}
	scriptVars["items"] = map[string]interface{}{}
	scriptVars["codeServices"] = map[string]interface{}{}
	scriptVars["codeLibraries"] = map[string]interface{}{}
	scriptVars["triggers"] = map[string]interface{}{}
	scriptVars["timers"] = map[string]interface{}{}
	scriptVars["devices"] = map[string]interface{}{}
	scriptVars["edges"] = map[string]interface{}{}
}

func extractCommand() (string, error) {
	if len(os.Args) < 2 {
		return "", fmt.Errorf("Missing command")
	}
	rval := os.Args[1]
	os.Args = append(os.Args[:1], os.Args[2:]...)
	return rval, nil
}

func getFileOrDie() string {
	args := flag.Args()
	if len(args) != 1 {
		fmt.Printf("Missing file argument\n")
		os.Exit(1)
	}
	return args[0]
}

func mustHaveOne(stuff ...string) {
	if len(stuff)%2 != 0 {
		fatal("Internal error: mustHave takes an even number of args\n")
	}
	argsPassed := ""
	for i := 0; i < len(stuff); i += 2 {
		key, val := stuff[i], stuff[i+1]
		argsPassed = argsPassed + key + " "
		if val != "" {
			return
		}
	}
	fatal(fmt.Sprintf("You need to pass one of these arguments: %s\n", argsPassed))
}

func mustHaveAll(stuff ...string) {
	if len(stuff)%2 != 0 {
		fatal("Internal error: mustHave takes an even number of args\n")
	}
	for i := 0; i < len(stuff); i += 2 {
		key, val := stuff[i], stuff[i+1]
		if val == "" {
			fatal(fmt.Sprintf("Missing '%s' argument\n", key))
		}
	}
}

func main() {
	warnScenetestPath()
	theCommand, err := extractCommand()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	flag.Parse()
	cb.CB_ADDR = PlatformAddr
	cb.CB_MSG_ADDR = MsgAddr
	setupSceneRoot()
	setupFileSearchPath()
	parseAndProcessEdgeInfo()

	if theCommand == "setup" {

		mustHaveAll("platformUrl", PlatformAddr, "messagingUrl", MsgAddr, "info", InfoFile)
		SetupFile = getFileOrDie()
		performSetup(getJSON(SetupFile))

	} else if theCommand == "run" {

		mustHaveAll("info", InfoFile)
		ScriptFile = getFileOrDie()
		scriptVars = getJSON(InfoFile)
		cb.CB_ADDR = scriptVars["platformUrl"].(string)
		cb.CB_MSG_ADDR = scriptVars["messagingUrl"].(string)
		executeTestScript(getJSON(ScriptFile))

	} else if theCommand == "teardown" {

		mustHaveAll("info", InfoFile)
		scriptVars = getJSON(InfoFile)
		cb.CB_ADDR = scriptVars["platformUrl"].(string)
		cb.CB_MSG_ADDR = scriptVars["messagingUrl"].(string)
		performTeardown()

	} else if theCommand == "help" {

		showHelp(flag.Args())

	} else {
		fmt.Printf("Unknown Command '%s'\n", theCommand)
		os.Exit(1)
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

func warnScenetestPath() {
	SceneRoot = os.Getenv(SceneTestEnvVar)
	if SceneRoot == "" {
		myPrintf("Warning: SCENETEST_PATH environment variable is not set. Help is disabled\n")
	}
}

func parseAndProcessEdgeInfo() {
	if EdgeInfo == "" {
		return
	}
	edgeStrings := strings.Split(EdgeInfo, ",")
	for _, oneEdgeString := range edgeStrings {
		splitEdge := strings.Split(oneEdgeString, "|")
		if len(splitEdge) != 4 {
			fatalf("Invalid edge specification: %s\n", oneEdgeString)
		}
		newEdge := &Edge{
			name:     splitEdge[0],
			ipAddr:   splitEdge[1],
			httpPort: splitEdge[2],
			mqttPort: splitEdge[3],
		}
		edgeMap[splitEdge[0]] = newEdge
	}
}

func getEdgeInfo(edgeName string) (*Edge, error) {
	if rval, ok := edgeMap[edgeName]; ok {
		return rval, nil
	}
	return nil, fmt.Errorf("Edge '%s' is unknown to scenetest", edgeName)
}

func (e *Edge) makeNiceAddrs() (string, string) {
	httpAddr := e.ipAddr + ":" + e.httpPort
	mqttThing := strings.TrimPrefix(e.ipAddr, "http://")
	mqttThing = strings.TrimPrefix(mqttThing, "https://")
	mqttAddr := mqttThing + ":" + e.mqttPort
	return httpAddr, mqttAddr
}
