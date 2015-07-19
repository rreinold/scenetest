package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
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

func saveSetupState(originalJson map[string]interface{}) {
	marshalled, err := json.MarshalIndent(setupState, "", "    ")
	if err != nil {
		fatal(fmt.Sprintf("MarshalIndent failed: %s\n", err.Error()))
	}
	fileToWrite := "setupState.json"
	if _, ok := originalJson["teardownFile"]; ok {
		fileToWrite = originalJson["teardownFile"].(string)
	}
	err = ioutil.WriteFile(fileToWrite, marshalled, os.ModePerm)
	if err != nil {
		fatal(fmt.Sprintf("Could not save setup state: %s\n", err.Error()))
	}

	scriptVars["platformUrl"] = PlatformAddr
	scriptVars["messagingUrl"] = MsgAddr

	marshalled, err = json.MarshalIndent(scriptVars, "", "    ")
	if err != nil {
		fatal(fmt.Sprintf("MarshalIndent failed: %s\n", err.Error()))
	}
	fileToWrite = "info.json"
	if _, ok := originalJson["infoFile"]; ok {
		fileToWrite = originalJson["infoFile"].(string)
	}
	err = ioutil.WriteFile(fileToWrite, marshalled, os.ModePerm)
	if err != nil {
		fatal(fmt.Sprintf("Could not save setup state: %s\n", err.Error()))
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
		if reflect.TypeOf(actualArg) != reflect.TypeOf(argType) {
			return fmt.Errorf("Argument #%d has type mismatch: %v != %v", reflect.TypeOf(actualArg), reflect.TypeOf(argType))
		}
	}
	return nil
}
