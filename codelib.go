package main

import (
	"fmt"
	cb "github.com/clearblade/Go-SDK"
)

func init() {
	funcMap["call"] = &Statement{call, callHelp}
}

func call(context map[string]interface{}, args []interface{}) error {
	if len(args) != 2 {
		return fmt.Errorf("Usage: [call, <serviceName>, {param, ...}]")
	}
	if _, ok := args[0].(string); !ok {
		return fmt.Errorf("Service name must be a string")
	}
	if _, ok := args[1].(map[string]interface{}); !ok {
		return fmt.Errorf("Service params must be a map")
	}
	svcName := valueOf(context, args[0]).(string)
	params := valueOf(context, args[1]).(map[string]interface{})
	sysKey := scriptVars["systemKey"].(string)
	userClient := context["userClient"].(*cb.UserClient)
	resp, err := userClient.CallService(sysKey, svcName, params)
	if err != nil {
		return err
	}
	context["returnValue"] = resp
	return nil
}

func callHelp() string {
	return "[\"call\", \"<serviceName>\", {<arg key/value pairs}]"
}
