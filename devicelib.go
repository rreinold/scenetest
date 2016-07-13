package main

import (
	"fmt"
	cb "github.com/clearblade/Go-SDK"
)

type createDevice struct{}
type updateDevice struct{}
type deleteDevice struct{}
type getDevice struct{}
type deviceConnectNovi struct{}

func init() {
	funcMap["createDevice"] = &createDevice{}
	funcMap["updateDevice"] = &updateDevice{}
	funcMap["deleteDevice"] = &deleteDevice{}
	funcMap["getDevice"] = &getDevice{}
	funcMap["deviceConnectNovi"] = &deviceConnectNovi{}
}

func (gd *getDevice) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 1, ""); err != nil {
		return nil, err
	}
	deviceName := args[0].(string)
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	devClient := ctx["adminClient"].(*cb.DevClient)
	return devClient.GetDevice(sysKey, deviceName)
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

func (ct *updateDevice) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 2, "", map[string]interface{}{}); err != nil {
		return nil, err
	}
	deviceName := args[0].(string)
	deviceChanges := args[1].(map[string]interface{})
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	client, err := getCurrentClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.UpdateDevice(sysKey, deviceName, deviceChanges)
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
func (dcn *deviceConnectNovi) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 1, ""); err != nil {
		return nil, err
	}
	deviceName := args[0].(string)
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	sysSec := scriptVars["systemSecret"].(string)
	deviceInfo := scriptVars["devices"].(map[string]interface{})[deviceName].(map[string]interface{})
	activeKey := deviceInfo["active_key"].(string)
	deviceClient := cb.NewDeviceClient(sysKey, sysSec, deviceName, activeKey)
	if _, err := deviceClient.AuthenticateDeviceWithKey(sysKey, deviceName, activeKey); err != nil {
		return nil, err
	}

	if err := deviceClient.InitializeMQTT("", "", 60); err != nil {
		return nil, err
	}
	if err := deviceClient.ConnectMQTT(nil, nil); err != nil {
		return nil, err
	}
	if err := deviceClient.Publish("/who/am/i", []byte(fmt.Sprintf("%p", deviceClient.MQTTClient)), 2); err != nil {
		return nil, err
	}

	ctx["deviceClient"] = deviceClient
	ctx["userClient"] = deviceClient // keeps existing code that uses userClient working
	ctx["deviceName"] = deviceName
	return deviceClient, nil
}

func (gd *getDevice) help() string {
	return "[\"getDevice\", \"deviceName\"]"
}

func (ct *createDevice) help() string {
	return "[\"createDevice\", \"deviceName\", {<device meta>}]"
}

func (ct *updateDevice) help() string {
	return "[\"updateDevice\", \"deviceName\", {<device changes>}]"
}

func (ct *deleteDevice) help() string {
	return "[\"deleteDevice\", \"deviceName\"]"
}

func (ct *deviceConnectNovi) help() string {
	return "[\"deviceConnectNovi\", \"deviceName\"]"
}
