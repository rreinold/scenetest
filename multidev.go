package main

type addDeveloper struct{}
type removeDeveloper struct{}
type getDevelopersForSystem struct{}
type getSystemsForDeveloper struct{}
type changeOwner struct{}

func init() {
	funcMap["addDeveloper"] = &addDeveloper{}
	funcMap["removeDeveloper"] = &removeDeveloper{}
	funcMap["getDevelopersForSYstem"] = &getDevelopersForSystem{}
	funcMap["getSystemsForDeveloper"] = &getSystemsForDeveloper{}
	funcMap["changeOwner"] = &changeOwner{}
}

func (ad *addDeveloper) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 1, ""); err != nil {
		return nil, err
	}
	devEmail := args[0].(string)
	changes := map[string]interface{}{"add": []string{devEmail}}
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	return adminClient.UpdateDevelopersForSystem(sysKey, changes)
}

func (rd *removeDeveloper) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 1, ""); err != nil {
		return nil, err
	}
	devEmail := args[0].(string)
	changes := map[string]interface{}{"remove": []string{devEmail}}
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	return adminClient.UpdateDevelopersForSystem(sysKey, changes)
}

func (ds *getDevelopersForSystem) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	return adminClient.GetDevelopersForSystem(sysKey)
}

func (sd *getSystemsForDeveloper) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 1, ""); err != nil {
		return nil, err
	}
	devEmail := args[0].(string)
	return adminClient.GetSystemsForDeveloper(devEmail)
}

func (co *changeOwner) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 1, ""); err != nil {
		return nil, err
	}
	devEmail := args[0].(string)
	changes := map[string]interface{}{"owner": devEmail}
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	return adminClient.UpdateDevelopersForSystem(sysKey, changes)
}

func (ad *addDeveloper) help() string {
	return "[\"addDeveloper\", \"developerEmail\"]"
}

func (rd *removeDeveloper) help() string {
	return "[\"removeDeveloper\", \"developerEmail\"]"
}

func (ds *getDevelopersForSystem) help() string {
	return "[\"getDevelopersForSystem\"]"
}

func (sd *getSystemsForDeveloper) help() string {
	return "[\"getSystemsForDeveloper\", \"developerEmail\"]"
}

func (co *changeOwner) help() string {
	return "[\"changeOwner\", \"developerEmail\"]"
}
