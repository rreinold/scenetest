package main

import (
	"fmt"
	cb "github.com/clearblade/Go-SDK"
)

type createEdge struct{}
type updateEdge struct{}
type deleteEdge struct{}
type getEdge struct{}

func init() {
	funcMap["createEdge"] = &createEdge{}
	//funcMap["updateEdge"] = &updateEdge{}
	funcMap["deleteEdge"] = &deleteEdge{}
	//funcMap["getEdge"] = &getEdge{}
}

func (ct *createEdge) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 2, "", map[string]interface{}{}); err != nil {
		return nil, err
	}
	edgeName := args[0].(string)
	edgeInput := args[1].(map[string]interface{})
	delete(edgeInput, "name")
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	devClient := ctx["adminClient"].(*cb.DevClient)
	fmt.Printf("CREATE EDGE '%s': %+#v\n", edgeName, edgeInput)
	return devClient.CreateEdge(sysKey, edgeName, edgeInput)
}

func (ct *deleteEdge) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 1, ""); err != nil {
		return nil, err
	}
	edgeName := args[0].(string)
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	devClient := ctx["adminClient"].(*cb.DevClient)
	return nil, devClient.DeleteEdge(sysKey, edgeName)
}

func (ct *deleteEdge) help() string {
	return "[\"deleteEdge\", \"edgeName\"]"
}

func (ct *createEdge) help() string {
	return "[\"createEdge\", \"edgeName\", {<edge meta>}]"
}
