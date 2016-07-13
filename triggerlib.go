package main

import (
	"encoding/json"
	"fmt"
	cb "github.com/clearblade/Go-SDK"
	mqtt "github.com/clearblade/mqtt_parsing"
	"time"
)

type createTrigger struct{}
type deleteTrigger struct{}
type createTimer struct{}
type deleteTimer struct{}
type waitTrigger struct{}
type subscribeTriggers struct{}

func init() {
	funcMap["createTrigger"] = &createTrigger{}
	funcMap["deleteTrigger"] = &deleteTrigger{}
	funcMap["createTimer"] = &createTimer{}
	funcMap["deleteTimer"] = &deleteTimer{}
	funcMap["waitTrigger"] = &waitTrigger{}
	funcMap["subscribeTriggers"] = &subscribeTriggers{}
}

func (ct *createTrigger) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 1, map[string]interface{}{}); err != nil {
		return nil, err
	}
	triggerInput := args[0].(map[string]interface{})
	triggerName := triggerInput["name"].(string)
	delete(triggerInput, triggerName)
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	devClient := ctx["adminClient"].(*cb.DevClient)
	return devClient.CreateEventHandler(sysKey, triggerName, triggerInput)
}

func (ct *createTrigger) help() string {
	return "[\"createTrigger\", {<trigger meta>}]"
}

func (ct *deleteTrigger) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 1, ""); err != nil {
		return nil, err
	}
	triggerName := args[0].(string)
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	devClient := ctx["adminClient"].(*cb.DevClient)
	return nil, devClient.DeleteEventHandler(sysKey, triggerName)
}

func (ct *deleteTrigger) help() string {
	return "[\"deleteTrigger\", \"triggerName\"]"
}

func (ct *createTimer) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 1, map[string]interface{}{}); err != nil {
		return nil, err
	}
	timerInput := args[0].(map[string]interface{})
	timerName := timerInput["name"].(string)
	delete(timerInput, timerName)
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	devClient := ctx["adminClient"].(*cb.DevClient)
	startTime := timerInput["start_time"].(string)
	timerInput["start_time"] = startTime
	return devClient.CreateTimer(sysKey, timerName, timerInput)
}

func (ct *createTimer) help() string {
	return "[\"createTimer\", {<timer meta>}]"
}

func (ct *deleteTimer) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 1, ""); err != nil {
		return nil, err
	}
	timerName := args[0].(string)
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	devClient := ctx["adminClient"].(*cb.DevClient)
	return nil, devClient.DeleteTimer(sysKey, timerName)
}

func (ct *deleteTimer) help() string {
	return "[\"deleteTimer\", \"timerName\"]"
}

func (st *subscribeTriggers) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	// No args
	if len(args) != 0 {
		return nil, fmt.Errorf("subscribeTriggers takes no arguments")
	}
	userClient := ctx["userClient"].(*cb.UserClient)
	triggerChan, err := userClient.Subscribe("/clearblade/internal/trigger", 0)
	if err != nil {
		return nil, err
	}
	ctx["triggerChannel"] = triggerChan
	return triggerChan, nil
}

func (st *subscribeTriggers) help() string {
	return "[\"subscribeTriggers\"]"
}

func (wt *waitTrigger) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 3, "", "", float64(0)); err != nil {
		return nil, err
	}
	eClass := args[0].(string)
	eType := args[1].(string)
	timeout := time.Duration(args[2].(float64))

	if _, ok := ctx["triggerChannel"]; !ok {
		return nil, fmt.Errorf("No trigger channel to wait on...")
	}
	trigChan := ctx["triggerChannel"].(<-chan *mqtt.Publish)

	for {
		var stuff *mqtt.Publish
		select {
		case stuff = <-trigChan:
		case <-time.After(time.Second * timeout):
			return nil, fmt.Errorf("Timed out waiting for trigger notification")
		}
		byts := stuff.Payload

		realStuff := map[string]interface{}{}
		if err := json.Unmarshal(byts, &realStuff); err != nil {
			return nil, fmt.Errorf("Unmarshal of trigger message payload failed: %s", err.Error())
		}

		if validateTrigger(eClass, eType, realStuff) {
			return nil, nil
		}
	}
}

func (wt *waitTrigger) help() string {
	return "waitTrigger help not yet implemented"
}

func validateTrigger(trigClass, trigType string, msgBody map[string]interface{}) bool {
	realClass := msgBody["msgClass"].(string)
	realType := msgBody["msgType"].(string)
	return (trigClass == realClass) && (trigType == realType)
}
