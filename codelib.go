package main

import (
	"fmt"
	cb "github.com/clearblade/Go-SDK"
)

func init() {
	funcMap["call"] = &Statement{call, callHelp}
	funcMap["createService"] = &Statement{createService, createServiceHelp}
	funcMap["updateService"] = &Statement{updateService, updateServiceHelp}
	funcMap["deleteService"] = &Statement{deleteService, deleteServiceHelp}
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
	context["returnValue"] = resp["results"]
	return nil
}

func callHelp() string {
	return "[\"call\", \"<serviceName>\", {<arg key/value pairs}]"
}

func fixParams(params []interface{}) []string {
	rval := make([]string, len(params))
	for idx, val := range params {
		rval[idx] = val.(string)
	}
	return rval
}

func createService(context map[string]interface{}, args []interface{}) error {
	//  Arg parsing nonsense
	if len(args) < 2 {
		return fmt.Errorf(createServiceHelp())
	}
	if _, ok := args[0].(string); !ok {
		return fmt.Errorf("Service name must be a string")
	}
	if _, ok := args[1].(string); !ok {
		return fmt.Errorf("Service code must be a string")
	}
	if len(args) == 2 {
		args = append(args, []interface{}{})
	}
	if _, ok := args[2].([]interface{}); !ok {
		return fmt.Errorf("Parameters arg must be a string\n")
	}

	//  Create the service
	svcName, code, params := args[0].(string), args[1].(string), fixParams(args[2].([]interface{}))
	adminClient := context["adminClient"].(*cb.DevClient)
	sysKey := scriptVars["systemKey"].(string)
	if err := adminClient.NewService(sysKey, svcName, code, params); err != nil {
		return err
	}

	//  As a convenience, add the Authenticated role
	if err := adminClient.AddServiceToRole(sysKey, svcName, "Authenticated", int(15)); err != nil {
		return err
	}

	return nil
}

func createServiceHelp() string {
	return "[\"createService\", \"<serviceName>\", \"<code>\", [<paramName>...]]"
}

func updateService(context map[string]interface{}, args []interface{}) error {
	//  Arg parsing nonsense
	if len(args) < 2 {
		return fmt.Errorf(updateServiceHelp())
	}
	if _, ok := args[0].(string); !ok {
		return fmt.Errorf("Service name must be a string")
	}
	if _, ok := args[1].(string); !ok {
		return fmt.Errorf("Service code must be a string")
	}
	if len(args) == 2 {
		args = append(args, []interface{}{})
	}
	if _, ok := args[2].([]interface{}); !ok {
		return fmt.Errorf("Parameters arg must be a string\n")
	}

	//  Create the service
	svcName, code, params := args[0].(string), args[1].(string), fixParams(args[2].([]interface{}))
	adminClient := context["adminClient"].(*cb.DevClient)
	sysKey := scriptVars["systemKey"].(string)
	if err := adminClient.UpdateService(sysKey, svcName, code, params); err != nil {
		return err
	}
	return nil
}

func updateServiceHelp() string {
	return "[\"updateService\", \"<serviceName>\", \"<code\", [<paramName>,...]]"
}

func deleteService(context map[string]interface{}, args []interface{}) error {
	//  Arg parsing nonsense
	if len(args) < 1 {
		return fmt.Errorf(deleteServiceHelp())
	}
	if _, ok := args[0].(string); !ok {
		return fmt.Errorf("Service name must be a string")
	}

	//  Create the service
	svcName := args[0].(string)
	adminClient := context["adminClient"].(*cb.DevClient)
	sysKey := scriptVars["systemKey"].(string)
	if err := adminClient.DeleteService(sysKey, svcName); err != nil {
		return err
	}
	return nil
}

func deleteServiceHelp() string {
	return "[\"deleteService\", \"<serviceName>\"]"
}
