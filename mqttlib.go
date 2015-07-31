package main

import (
	"fmt"
	cb "github.com/clearblade/Go-SDK"
	mqtt "github.com/clearblade/mqtt_parsing"
	"time"
)

func init() {
	funcMap["subscribe"] = &Statement{subscribe, subscribeHelp}
	funcMap["publish"] = &Statement{publish, publishHelp}
	funcMap["waitMessage"] = &Statement{waitMessage, waitMessageHelp}
}

func subscribe(context map[string]interface{}, args []interface{}) error {
	if len(args) != 2 {
		return fmt.Errorf("Usage: [subscribe, topic, qos(int)]")
	}
	topic := valueOf(context, args[0]).(string)
	qos := int(valueOf(context, args[1]).(float64))
	userClient := context["userClient"].(*cb.UserClient)
	trigChan, err := userClient.Subscribe(topic, qos)
	if err != nil {
		return err
	}
	context[topic] = trigChan
	return nil
}

func subscribeHelp() string {
	return "subscribe help not yet implemented"
}

func waitMessage(context map[string]interface{}, args []interface{}) error {
	if len(args) != 2 {
		fmt.Errorf("Usage: [waitMessage, topic, timeout]")
	}
	topic := valueOf(context, args[0]).(string)
	timeout := time.Duration(args[1].(float64))
	trigChan := context[topic].(<-chan *mqtt.Publish)
	var stuff *mqtt.Publish
	select {
	case stuff = <-trigChan:
	case <-time.After(time.Second * timeout):
		return fmt.Errorf("Timed out waiting for message to arrive on %s", topic)
	}
	realStuff := string(stuff.Payload)
	context["returnValue"] = realStuff
	return nil
}

func waitMessageHelp() string {
	return "[\"waitMessage\", \"topic\", timeout]"
}

func publish(context map[string]interface{}, args []interface{}) error {
	if len(args) != 3 {
		fmt.Errorf("Usage: [publish, topic, message_body, qos]")
	}
	userClient := context["userClient"].(*cb.UserClient)
	topic := valueOf(context, args[0]).(string)
	body := []byte(valueOf(context, args[1]).(string))
	qos := int(valueOf(context, args[2]).(float64))
	if err := userClient.Publish(topic, body, qos); err != nil {
		return fmt.Errorf("Publish failed: %s", err.Error())
	}
	return nil
}

func publishHelp() string {
	return "publish help not yet implemented"
}
