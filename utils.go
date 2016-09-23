package main

import (
	"encoding/json"
	"fmt"
	cb "github.com/clearblade/Go-SDK"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strings"
)

func getVarOrFile(from map[string]interface{}, key string) string {
	ival, ok := from[key]
	if !ok {
		fatal(fmt.Sprintf("Undefined key \"%s\n", key))
	}
	val := ival.(string)
	if strings.HasPrefix(val, "@") {
		byts, err := ioutil.ReadFile(strings.TrimPrefix(val, "@"))
		if err != nil {
			fatal(fmt.Sprintf("Could not read file %s: %s\n", val, err.Error()))
		}
		return string(byts)
	}
	return val
}

func readFileFromEnvironment(filename string) (string, error) {
	filename = strings.TrimPrefix(filename, "@")
	thisSearchPath := []string{""}
	if !strings.HasPrefix(filename, "/") {
		thisSearchPath = FileSearchPath
	}
	for _, prefix := range thisSearchPath {
		if bytes, err := ioutil.ReadFile(prefix + filename); err == nil {
			return string(bytes), nil
		}
	}
	return "", fmt.Errorf("File '%s' not found in environment", filename)
}

func saveSetupState(originalJson map[string]interface{}) {
	setupState["platformUrl"] = PlatformAddr
	setupState["messagingUrl"] = MsgAddr
	marshalled, err := json.MarshalIndent(setupState, "", "    ")
	if err != nil {
		fatal(fmt.Sprintf("MarshalIndent failed: %s\n", err.Error()))
	}

	scriptVars["platformUrl"] = PlatformAddr
	scriptVars["messagingUrl"] = MsgAddr
	scriptVars["teardown"] = setupState

	marshalled, err = json.MarshalIndent(scriptVars, "", "    ")
	if err != nil {
		fatal(fmt.Sprintf("MarshalIndent failed: %s\n", err.Error()))
	}

	fmt.Printf("THE INFO FILE IS: %s\n", InfoFile)
	err = ioutil.WriteFile(InfoFile, marshalled, os.ModePerm)
	if err != nil {
		fatal(fmt.Sprintf("Could not save setup information to '%s':%s\n", InfoFile, err.Error()))
	}
}

func argCheck(args []interface{}, mandatory int, argTypes ...interface{}) error {
	if len(args) < mandatory {
		return fmt.Errorf("Not enough arguments")
	}
	if len(args) > len(argTypes) {
		return fmt.Errorf("Too many arguments")
	}
	for i, actualArg := range args {
		argType := argTypes[i]
		if argType == nil {
			continue // nil means interface{}
		}
		if reflect.TypeOf(actualArg) != reflect.TypeOf(argType) {
			return fmt.Errorf("Argument #%d has type mismatch: %v != %v", i, reflect.TypeOf(actualArg), reflect.TypeOf(argType))
		}
	}
	return nil
}

func valueOf(context map[string]interface{}, thing interface{}) interface{} {
	if isAStatement(thing) {
		stmt := thing.([]interface{})
		res, err := runOneStep(context, stmt)
		if err != nil {
			return err
			//fatal(fmt.Sprintf("Substatement execution failed: %s", err.Error()))
		}
		thing = res
	}
	switch thing.(type) {
	case string:
		thingStr := thing.(string)
		if strings.HasPrefix(thingStr, "@") {
			varName := strings.TrimPrefix(thingStr, "@")

			//  first check locals (context)
			if val, ok := context[varName]; ok {
				return val
			}

			//  now check globals
			if val, ok := globals[varName]; ok {
				return val
			}

			//  Var doesn't exist -- see if there's a file to read from...
			if contents, err := readFileFromEnvironment(varName); err == nil {
				return contents
			}
			myPrintf("DEAD __ CONTEXT: %+v\n", context)
			fatal(fmt.Sprintf("Undefined variable: %s", varName))
		}
	}
	return thing
}

func showHelp(args []string) {
	if SceneTestEnvVar == "" {
		return
	}
	helpStuff := ""
	if len(args) == 0 {
		helpStuff = showHelpArgs()
	} else {
		for _, arg := range args {
			helpStuff += getHelpContents(arg)
		}
	}

	myPrintf("%s", helpStuff)
}

func getHelpContents(name string) string {
	fileName := name + ".txt"
	fileNameFullPath := SceneRoot + "/help/" + fileName
	byts, err := ioutil.ReadFile(fileNameFullPath)
	if err != nil {
		fatal(fmt.Sprintf("Could not read help file '%s'", fileNameFullPath))
	}
	return string(byts)
}

func showHelpArgs() string {
	rval := ""
	rval += "\nUsage: scenetest help <commandOrStatementName> ...\n"
	rval += "\nEnter either one or more of the command names or statement names below\n"
	rval += "(ie \"scenetest help setup\", \"scenetest help createItem\", ...)\n\n"
	rval += "Commands:\n"
	rval += "\tsetup\n\trun\n\tteardown\n"
	exprStrings := []string{}
	stmtStrings := []string{}
	for key, theStmt := range funcMap {
		if isExprStmt(theStmt) {
			exprStrings = append(exprStrings, key)
		} else {
			stmtStrings = append(stmtStrings, key)
		}
	}
	sort.Strings(exprStrings)
	sort.Strings(stmtStrings)
	rval += "\nExpressions:\n"
	for _, exprString := range exprStrings {
		rval += fmt.Sprintf("\t%s\n", exprString)
	}
	rval += "\nStatements:\n"
	for _, stmtString := range stmtStrings {
		rval += fmt.Sprintf("\t%s\n", stmtString)
	}
	return rval
}

func setupSceneRoot() {
	SceneRoot = os.Getenv(SceneTestEnvVar)
	if SceneRoot == "" {
		SceneRoot = "./"
	}
}

func setupFileSearchPath() {
	FileSearchPath = []string{}
	if SceneRoot != "./" {
		FileSearchPath = append(FileSearchPath, SceneRoot)
	}
	FileSearchPath = append(FileSearchPath, "./js/")
	FileSearchPath = append(FileSearchPath, "./")
}

func incrNestingLevel(ctx map[string]interface{}) {
	ctx["__nestingLevel"] = ctx["__nestingLevel"].(int) + 1
}

func decrNestingLevel(ctx map[string]interface{}) {
	ctx["__nestingLevel"] = ctx["__nestingLevel"].(int) - 1
}

func isAStatement(arg interface{}) bool {
	slice, ok := arg.([]interface{})
	if !ok || len(slice) <= 0 {
		return false
	}

	name, ok := slice[0].(string)
	if !ok {
		return false
	}

	if _, ok := funcMap[name]; !ok {
		return false
	}
	return true
}

func lookupVar(context map[string]interface{}, varName string) interface{} {
	//  first check locals (context)
	if val, ok := context[varName]; ok {
		return val
	}

	//  now check globals
	if val, ok := globals[varName]; ok {
		return val
	}
	return nil
}

func getCurrentClient(ctx map[string]interface{}) (cb.Client, error) {
	list := []string{"userClient", "deviceClient", "developerClient", "adminClient"}
	for _, clientName := range list {
		if stuff, ok := ctx[clientName]; ok {
			return stuff.(cb.Client), nil
		}
	}
	return nil, fmt.Errorf("No clients have been created for this scenario yet")
}
