package main

import (
	"encoding/json"
	"fmt"
	cb "github.com/clearblade/Go-SDK"
	mqtt "github.com/clearblade/mqtt_parsing"
	"time"
)

type subscribeStmt struct{}
type publishStmt struct{}
type waitMessageStmt struct{}

func init() {
	funcMap["subscribe"] = &subscribeStmt{}
	funcMap["publish"] = &publishStmt{}
	funcMap["waitMessage"] = &waitMessageStmt{}
}

func (s *subscribeStmt) run(context map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 2, "", float64(0)); err != nil {
		return nil, err
	}
	topic := args[0].(string)
	qos := int(args[1].(float64))
	userClient := context["userClient"].(*cb.UserClient)
	trigChan, err := userClient.Subscribe(topic, qos)
	if err != nil {
		return nil, err
	}
	context[topic] = trigChan
	return trigChan, nil
}

func (s *subscribeStmt) help() string {
	return "[\"subscribe\", <topicString>, <QoSLevel>]"
}

func (w *waitMessageStmt) run(context map[string]interface{}, args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("Usage: [waitMessage, topic, timeout]")
	}
	topic := args[0].(string)
	timeout := time.Duration(args[1].(float64))
	trigChan := context[topic].(<-chan *mqtt.Publish)
	var stuff *mqtt.Publish
	select {
	case stuff = <-trigChan:
	case <-time.After(time.Second * timeout):
		return nil, fmt.Errorf("Timed out waiting for message to arrive on %s", topic)
	}
	realStuff := string(stuff.Payload)
	return realStuff, nil
}

func (w *waitMessageStmt) help() string {
	return "[\"waitMessage\", \"topic\", timeout]"
}

func (p *publishStmt) run(context map[string]interface{}, args []interface{}) (retVal interface{}, err error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("Usage: [publish, topic, message_body, qos]")
	}
	userClient := context["userClient"].(*cb.UserClient)
	topic := args[0].(string)
	var body []byte
	switch args[1].(type) {
	case string:
		body = []byte(args[1].(string))
	case map[string]interface{}:
		if body, err = json.Marshal(args[1].(map[string]interface{})); err != nil {
			return nil, fmt.Errorf("Can't serialize message body %s", err.Error())
		}
	default:
		return nil, fmt.Errorf("Unsupported message body")
	}
	qos := int(args[2].(float64))
	if err := userClient.Publish(topic, body, qos); err != nil {
		return nil, fmt.Errorf("Publish failed: %s", err.Error())
	}
	return nil, nil
}

func (p *publishStmt) help() string {
	return "publish help not yet implemented"
}
