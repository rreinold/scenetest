package main

import (
	"encoding/json"
	"fmt"
	cb "github.com/clearblade/Go-SDK"
	mqtt "github.com/clearblade/mqtt_parsing"
	"time"
)

type createTrigger struct{}
type createTimer struct{}
type waitTrigger struct{}
type subscribeTriggers struct{}

func init() {
	funcMap["createTrigger"] = &createTrigger{}
	funcMap["createTimer"] = &createTimer{}
	funcMap["waitTrigger"] = &waitTrigger{}
	funcMap["subscribeTriggers"] = &subscribeTriggers{}
	funcMap["subscribeTrigger"] = &subscribeTriggers{}
}

func (ct *createTrigger) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	return nil, nil
}

func (ct *createTrigger) help() string {
	return "createTrigger help not yet implemented"
}

func (ct *createTimer) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Wrong number of args to create timer: %d", len(args))
	}
	timerInput, ok := args[0].(map[string]interface{})
	timerName := timerInput["name"].(string)
	delete(timerInput, timerName)
	if !ok {
		return nil, fmt.Errorf("Argument to createTimer must be a map of attributes")
	}
	sysKey := scriptVars["systemKey"].(string)
	devClient := ctx["adminClient"].(*cb.DevClient)
	startTime := timerInput["start_time"].(string)
	timerInput["start_time"] = startTime
	newTimer, err := devClient.CreateTimer(sysKey, timerName, timerInput)
	if err != nil {
		return nil, err
	}
	return newTimer, nil
}

func (ct *createTimer) help() string {
	return "[\"createTimer\", {<timer meta>}]"
}

func (st *subscribeTriggers) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	// No args
	if len(args) != 0 {
		return nil, fmt.Errorf("subscribeTriggers takes no arguments")
	}
	userClient := ctx["userClient"].(*cb.UserClient)
	myPrintf("Doing the subscribe\n")
	triggerChan, err := userClient.Subscribe("/clearblade/internal/trigger", 0)
	if err != nil {
		return nil, err
	}
	myPrintf("Done subscribing\n")
	ctx["triggerChannel"] = triggerChan
	return triggerChan, nil
}

func (st *subscribeTriggers) help() string {
	return "[\"subscribeTriggers\"]"
}

func (wt *waitTrigger) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("Wrong number of arguments to waitTrigger: %d", len(args))
	}
	eClass := valueOf(ctx, args[0]).(string)
	eType := valueOf(ctx, args[1]).(string)
	timeout := time.Duration(valueOf(ctx, args[2]).(float64))

	if _, ok := ctx["triggerChannel"]; !ok {
		return nil, fmt.Errorf("No trigger channel to wait on...")
	}
	trigChan := ctx["triggerChannel"].(<-chan *mqtt.Publish)

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

	if err := validateTrigger(eClass, eType, realStuff); err != nil {
		return nil, err
	}
	return nil, nil
}

func (wt *waitTrigger) help() string {
	return "waitTrigger help not yet implemented"
}

func validateTrigger(trigClass, trigType string, msgBody map[string]interface{}) error {
	realClass := msgBody["msgClass"].(string)
	realType := msgBody["msgType"].(string)
	if trigClass != realClass {
		return fmt.Errorf("Bad message class: %s; expected %s", realClass, trigClass)
	}
	if trigType != realType {
		return fmt.Errorf("Bad message type: %s; expected %s", realType, trigType)
	}
	return nil
}
