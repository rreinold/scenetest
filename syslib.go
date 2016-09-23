package main

import (
	//"fmt"
	cb "github.com/clearblade/Go-SDK"
)

type createSystemStmt struct{}
type deleteSystemStmt struct{}
type getSystemsForDeveloperStmt struct{}
type getDevelopersForSystemStmt struct{}
type updateDevelopersForSystemStmt struct{}

var (
	sysOwnerMap = map[string]string{} // TODO USE THIS LATER
)

func init() {
	funcMap["createSystem"] = &createSystemStmt{}
	funcMap["deleteSystem"] = &deleteSystemStmt{}
	funcMap["getSystemsForDeveloper"] = &getSystemsForDeveloperStmt{}
	funcMap["getDevelopersForSystem"] = &getDevelopersForSystemStmt{}
	funcMap["updateDevelopersForSystem"] = &updateDevelopersForSystemStmt{}
}

func (cs *createSystemStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 2, "", ""); err != nil {
		return nil, err
	}
	sysName := args[0].(string)
	sysDescr := args[1].(string)
	client := ctx["developerClient"].(*cb.DevClient)
	sysKey, err := client.NewSystem(sysName, sysDescr, true)
	if err != nil {
		return nil, err
	}
	return sysKey, nil
}

func (ds *deleteSystemStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 1, "", map[string]interface{}{}); err != nil {
		return nil, err
	}
	sysKey := args[0].(string)
	client := ctx["developerClient"].(*cb.DevClient)
	return nil, client.DeleteSystem(sysKey)
}

func (gs *getSystemsForDeveloperStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 1, ""); err != nil {
		return nil, err
	}
	client := ctx["developerClient"].(*cb.DevClient)
	developerEmail := args[0].(string)
	stuff, err := client.GetSystemsForDeveloper(developerEmail)
	if err != nil {
		return nil, err
	}
	return stuff, nil
}

func (gd *getDevelopersForSystemStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 1, ""); err != nil {
		return nil, err
	}
	client := ctx["developerClient"].(*cb.DevClient)
	systemKey := args[0].(string)
	stuff, err := client.GetDevelopersForSystem(systemKey)
	if err != nil {
		return nil, err
	}
	return stuff, nil
}

func (gs *updateDevelopersForSystemStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 2, "", map[string]interface{}{}); err != nil {
		return nil, err
	}
	systemKey := args[0].(string)
	changes := args[1].(map[string]interface{})
	// NOTE: This needs to change as we have to be the right developer for the system. Odds
	//  are "adminClient" will be wrong some of the time.
	client := ctx["developerClient"].(*cb.DevClient)
	return client.UpdateDevelopersForSystem(systemKey, changes)
}

func (gd *createSystemStmt) help() string {
	return "[\"createSystem\", \"<System Name>\", \"<System Description>\"]"
}

func (ct *deleteSystemStmt) help() string {
	return "[\"deleteSystem\", \"<systemKey>\"]"
}

func (ct *getSystemsForDeveloperStmt) help() string {
	return "[\"getSystemsForDeveloper\", \"<developerEmail>\"]"
}

func (ct *getDevelopersForSystemStmt) help() string {
	return "[\"getDevelopersForSystem\", \"<systemKey>\"]"
}

func (ct *updateDevelopersForSystemStmt) help() string {
	return "[\"updateDevelopersForSystem\", \"<systemKey>\", { \"add\": [...], \"remove\": [...], \"owner\": \"<ownerEmail\" }]"
}
