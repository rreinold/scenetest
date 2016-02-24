package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	//"sort"
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
	/*
		fileToWrite := "setupState.json"
		if _, ok := originalJson["teardownFile"]; ok {
			fileToWrite = originalJson["teardownFile"].(string)
		}
		err = ioutil.WriteFile(fileToWrite, marshalled, os.ModePerm)
		if err != nil {
			fatal(fmt.Sprintf("Could not save setup state: %s\n", err.Error()))
		}
	*/

	scriptVars["platformUrl"] = PlatformAddr
	scriptVars["messagingUrl"] = MsgAddr
	scriptVars["teardown"] = setupState

	marshalled, err = json.MarshalIndent(scriptVars, "", "    ")
	if err != nil {
		fatal(fmt.Sprintf("MarshalIndent failed: %s\n", err.Error()))
	}
	/*
		fileToWrite = "info.json"
		if _, ok := originalJson["infoFile"]; ok {
			fileToWrite = originalJson["infoFile"].(string)
		}
	*/
	err = ioutil.WriteFile(InfoFile, marshalled, os.ModePerm)
	if err != nil {
		fatal(fmt.Sprintf("Could not save setup information to '%s': %s\n", InfoFile, err.Error()))
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
			return fmt.Errorf("Argument #%d has type mismatch: %v != %v", reflect.TypeOf(actualArg), reflect.TypeOf(argType))
		}
	}
	return nil
}

func valueOf(context map[string]interface{}, thing interface{}) interface{} {
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

func showHelp() {
	/*
		keys := make([]string, len(funcMap))
		i := 0
		for k := range funcMap {
			keys[i] = k
			i += 1
		}
		sort.Strings(keys)
		for _, funcName := range keys {
			stmt := funcMap[funcName]
			myPrintf("%s\n", stmt.HelpFunc())
		}
	*/
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
