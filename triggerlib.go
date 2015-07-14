package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/clearblade/mqtt_parsing"
	"time"
)

func init() {
	funcMap["createTrigger"] = createTrigger
	funcMap["waitTrigger"] = waitTrigger
}

func createTrigger(ctx map[string]interface{}, args []interface{}) error {
	fmt.Printf("In createTrigger: %+v\n", args)
	return nil
}

func waitTrigger(ctx map[string]interface{}, args []interface{}) error {
	if len(args) != 3 {
		return fmt.Errorf("Wrong number of arguments to waitTrigger: %d", len(args))
	}
	eClass := args[0].(string)
	eType := args[1].(string)
	timeout := time.Duration(args[2].(float64))

	if _, ok := ctx["triggerChannel"]; !ok {
		return fmt.Errorf("No trigger channel to wait on...")
	}
	trigChan := ctx["triggerChannel"].(<-chan *mqtt.Publish)

	var stuff *mqtt.Publish
	select {
	case stuff = <-trigChan:
	case <-time.After(time.Second * timeout):
		return fmt.Errorf("Timed out waiting for trigger notification")
	}
	byts := stuff.Payload

	realStuff := map[string]interface{}{}
	if err := json.Unmarshal(byts, &realStuff); err != nil {
		return fmt.Errorf("Unmarshal of trigger message payload failed: %s", err.Error())
	}

	if err := validateTrigger(eClass, eType, realStuff); err != nil {
		return err
	}
	fmt.Printf("Son bitch\n")
	return nil
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
