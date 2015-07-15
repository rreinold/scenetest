package main

import (
	"fmt"
	cb "github.com/clearblade/Go-SDK"
	mqtt "github.com/clearblade/mqtt_parsing"
	"time"
)

func init() {
	funcMap["subscribe"] = subscribe
	funcMap["publish"] = publish
	funcMap["waitMessage"] = waitMessage
}

func subscribe(context map[string]interface{}, args []interface{}) error {
	if len(args) != 2 {
		return fmt.Errorf("Usage: [subscribe, topic, qos(int)]")
	}
	topic := args[0].(string)
	qos := int(args[1].(float64))
	userClient := context["userClient"].(*cb.UserClient)
	trigChan, err := userClient.Subscribe(topic, qos)
	if err != nil {
		return err
	}
	context[topic] = trigChan
	return nil
}

func waitMessage(context map[string]interface{}, args []interface{}) error {
	if len(args) != 2 {
		fmt.Errorf("Usage: [waitMessage, topic, timeout]")
	}
	topic := args[0].(string)
	timeout := time.Duration(args[1].(float64))
	trigChan := context[topic].(<-chan *mqtt.Publish)
	var stuff *mqtt.Publish
	select {
	case stuff = <-trigChan:
	case <-time.After(time.Second * timeout):
		return fmt.Errorf("Timed out waiting for message to arrive on %s", topic)
	}
	realStuff := string(stuff.Payload)
	fmt.Printf("Got message %s\n", realStuff)
	context["returnValue"] = realStuff
	return nil
}

func publish(context map[string]interface{}, args []interface{}) error {
	if len(args) != 3 {
		fmt.Errorf("Usage: [publish, topic, message_body, qos]")
	}
	userClient := context["userClient"].(*cb.UserClient)
	fmt.Printf("ARGS ARE: %+v\n", args)
	topic := args[0].(string)
	body := []byte(args[1].(string))
	qos := int(args[2].(float64))
	if err := userClient.Publish(topic, body, qos); err != nil {
		return fmt.Errorf("Publish failed: %s", err.Error())
	}
	return nil
}
