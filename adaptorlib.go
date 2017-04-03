package main

type createAdaptor struct{}
type updateAdaptor struct{}
type deleteAdaptor struct{}
type getAdaptor struct{}

type createAdaptorFile struct{}
type updateAdaptorFile struct{}
type deleteAdaptorFile struct{}
type getAdaptorFile struct{}

func init() {
	funcMap["createAdaptor"] = &createAdaptor{}
	funcMap["updateAdaptor"] = &updateAdaptor{}
	funcMap["deleteAdaptor"] = &deleteAdaptor{}
	funcMap["getAdaptor"] = &getAdaptor{}
	funcMap["createAdaptorFile"] = &createAdaptorFile{}
	funcMap["updateAdaptorFile"] = &updateAdaptorFile{}
	funcMap["deleteAdaptorFile"] = &deleteAdaptorFile{}
	funcMap["getAdaptorFile"] = &getAdaptorFile{}
}

func (gd *getAdaptor) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 1, ""); err != nil {
		return nil, err
	}
	adaptorName := args[0].(string)
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	client, err := getCurrentClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetAdaptor(sysKey, adaptorName)
}

func (ct *createAdaptor) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 2, "", map[string]interface{}{}); err != nil {
		return nil, err
	}
	adaptorName := args[0].(string)
	adaptorInput := args[1].(map[string]interface{})
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	return adminClient.CreateAdaptor(sysKey, adaptorName, adaptorInput)
}

func (ct *updateAdaptor) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 2, "", map[string]interface{}{}); err != nil {
		return nil, err
	}
	adaptorName := args[0].(string)
	adaptorChanges := args[1].(map[string]interface{})
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	client, err := getCurrentClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.UpdateAdaptor(sysKey, adaptorName, adaptorChanges)
}

func (ct *deleteAdaptor) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 1, ""); err != nil {
		return nil, err
	}
	adaptorName := args[0].(string)
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	client, err := getCurrentClient(ctx)
	if err != nil {
		return nil, err
	}
	return nil, client.DeleteAdaptor(sysKey, adaptorName)
}

func (gd *getAdaptorFile) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 1, ""); err != nil {
		return nil, err
	}
	adaptorName := args[0].(string)
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	client, err := getCurrentClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetAdaptorFile(sysKey, adaptorName)
}

func (ct *createAdaptorFile) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 3, "", "", map[string]interface{}{}); err != nil {
		return nil, err
	}
	adaptorName := args[0].(string)
	adaptorFileName := args[1].(string)
	adaptorInput := args[2].(map[string]interface{})
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	return adminClient.CreateAdaptorFile(sysKey, adaptorName, adaptorFileName, adaptorInput)
}

func (ct *updateAdaptorFile) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 2, "", map[string]interface{}{}); err != nil {
		return nil, err
	}
	adaptorName := args[0].(string)
	adaptorChanges := args[1].(map[string]interface{})
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	client, err := getCurrentClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.UpdateAdaptorFile(sysKey, adaptorName, adaptorChanges)
}

func (ct *deleteAdaptorFile) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 1, ""); err != nil {
		return nil, err
	}
	adaptorName := args[0].(string)
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	client, err := getCurrentClient(ctx)
	if err != nil {
		return nil, err
	}
	return nil, client.DeleteAdaptorFile(sysKey, adaptorName)
}

func (gd *getAdaptor) help() string {
	return "[\"getAdaptor\", \"adaptorName\"]"
}

func (ct *createAdaptor) help() string {
	return "[\"createAdaptor\", \"adaptorName\", {<adaptor meta>}]"
}

func (ct *updateAdaptor) help() string {
	return "[\"updateAdaptor\", \"adaptorName\", {<adaptor changes>}]"
}

func (ct *deleteAdaptor) help() string {
	return "[\"deleteAdaptor\", \"adaptorName\"]"
}

func (gd *getAdaptorFile) help() string {
	return "[\"getAdaptorFile\", \"adaptorName\", \"adaptorFileName\"]"
}

func (ct *createAdaptorFile) help() string {
	return "[\"createAdaptorFile\", \"adaptorName\", \"adaptorFileName\",  {<adaptor meta>}]"
}

func (ct *updateAdaptorFile) help() string {
	return "[\"updateAdaptorFile\", \"adaptorName\", \"adaptorFileName\", {<adaptor changes>}]"
}

func (ct *deleteAdaptorFile) help() string {
	return "[\"deleteAdaptorFile\", \"adaptorName\", \"adaptorFileName\"]"
}
