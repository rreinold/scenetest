package main

import (
	"fmt"
	cb "github.com/clearblade/Go-SDK"
)

type callStmt struct{}
type createServiceStmt struct{}
type updateServiceStmt struct{}
type deleteServiceStmt struct{}

func init() {
	funcMap["call"] = &callStmt{}
	funcMap["createService"] = &createServiceStmt{}
	funcMap["updateService"] = &updateServiceStmt{}
	funcMap["deleteService"] = &deleteServiceStmt{}
}

func (c *callStmt) run(context map[string]interface{}, args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("Usage: [call, <serviceName>, {param, ...}]")
	}
	if _, ok := args[0].(string); !ok {
		return nil, fmt.Errorf("Service name must be a string")
	}
	if _, ok := args[1].(map[string]interface{}); !ok {
		return nil, fmt.Errorf("Service params must be a map")
	}
	svcName := args[0].(string)
	params := args[1].(map[string]interface{})
	sysKey := scriptVars["systemKey"].(string)
	userClient := context["userClient"].(*cb.UserClient)
	resp, err := userClient.CallService(sysKey, svcName, params)
	if err != nil {
		return nil, err
	}
	//context["returnValue"] = resp["results"]
	return resp["results"], nil
}

func (c *callStmt) help() string {
	return "[\"call\", \"<serviceName>\", {<arg key/value pairs}]"
}

func fixParams(params []interface{}) []string {
	rval := make([]string, len(params))
	for idx, val := range params {
		rval[idx] = val.(string)
	}
	return rval
}

func (c *createServiceStmt) run(context map[string]interface{}, args []interface{}) (interface{}, error) {
	//  Arg parsing nonsense
	if len(args) < 2 {
		return nil, fmt.Errorf(c.help())
	}
	if _, ok := args[0].(string); !ok {
		return nil, fmt.Errorf("Service name must be a string")
	}
	if _, ok := args[1].(string); !ok {
		return nil, fmt.Errorf("Service code must be a string")
	}
	if len(args) == 2 {
		args = append(args, []interface{}{})
	}
	if _, ok := args[2].([]interface{}); !ok {
		return nil, fmt.Errorf("Parameters arg must be a string\n")
	}

	//  Create the service
	svcName, code, params := args[0].(string), args[1].(string), fixParams(args[2].([]interface{}))
	adminClient := context["adminClient"].(*cb.DevClient)
	sysKey := scriptVars["systemKey"].(string)
	if err := adminClient.NewService(sysKey, svcName, code, params); err != nil {
		return nil, err
	}

	//  As a convenience, add the Authenticated role
	if err := adminClient.AddServiceToRole(sysKey, svcName, "Authenticated", int(15)); err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *createServiceStmt) help() string {
	return "[\"createService\", \"<serviceName>\", \"<code>\", [<paramName>...]]"
}

func (u *updateServiceStmt) run(context map[string]interface{}, args []interface{}) (interface{}, error) {
	//  Arg parsing nonsense
	if len(args) < 2 {
		return nil, fmt.Errorf(u.help())
	}
	if _, ok := args[0].(string); !ok {
		return nil, fmt.Errorf("Service name must be a string")
	}
	if _, ok := args[1].(string); !ok {
		return nil, fmt.Errorf("Service code must be a string")
	}
	if len(args) == 2 {
		args = append(args, []interface{}{})
	}
	if _, ok := args[2].([]interface{}); !ok {
		return nil, fmt.Errorf("Parameters arg must be a string\n")
	}

	//  Create the service
	svcName, code, params := args[0].(string), args[1].(string), fixParams(args[2].([]interface{}))
	adminClient := context["adminClient"].(*cb.DevClient)
	sysKey := scriptVars["systemKey"].(string)
	if err := adminClient.UpdateService(sysKey, svcName, code, params); err != nil {
		return nil, err
	}
	return nil, nil
}

func (u *updateServiceStmt) help() string {
	return "[\"updateService\", \"<serviceName>\", \"<code\", [<paramName>,...]]"
}

func (d *deleteServiceStmt) run(context map[string]interface{}, args []interface{}) (interface{}, error) {
	//  Arg parsing nonsense
	if len(args) < 1 {
		return nil, fmt.Errorf(d.help())
	}
	if _, ok := args[0].(string); !ok {
		return nil, fmt.Errorf("Service name must be a string")
	}

	//  Create the service
	svcName := args[0].(string)
	adminClient := context["adminClient"].(*cb.DevClient)
	sysKey := scriptVars["systemKey"].(string)
	if err := adminClient.DeleteService(sysKey, svcName); err != nil {
		return nil, err
	}
	return nil, nil
}

func (d *deleteServiceStmt) help() string {
	return "[\"deleteService\", \"<serviceName>\"]"
}
