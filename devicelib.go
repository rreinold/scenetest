package main

import (
	//"fmt"
	cb "github.com/clearblade/Go-SDK"
)

type createDevice struct{}
type updateDevice struct{}
type deleteDevice struct{}
type getDevice struct{}

func init() {
	funcMap["createDevice"] = &createDevice{}
	//funcMap["updateDevice"] = &updateDevice{}
	funcMap["deleteDevice"] = &deleteDevice{}
	//funcMap["getDevice"] = &getDevice{}
}

func (ct *createDevice) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 2, "", map[string]interface{}{}); err != nil {
		return nil, err
	}
	deviceName := args[0].(string)
	deviceInput := args[1].(map[string]interface{})
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	devClient := ctx["adminClient"].(*cb.DevClient)
	return devClient.CreateDevice(sysKey, deviceName, deviceInput)
}

func (ct *deleteDevice) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 1, ""); err != nil {
		return nil, err
	}
	deviceName := args[0].(string)
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	devClient := ctx["adminClient"].(*cb.DevClient)
	return nil, devClient.DeleteDevice(sysKey, deviceName)
}

func (ct *deleteDevice) help() string {
	return "[\"deleteDevice\", \"deviceName\"]"
}

func (ct *createDevice) help() string {
	return "[\"createDevice\", \"deviceName\", {<device meta>}]"
}
