package main

type createAdaptor struct{}
type updateAdaptor struct{}
type deleteAdaptor struct{}
type getAdaptor struct{}
type deployAdaptor struct{}
type controlAdaptor struct{}

type createAdaptorFile struct{}
type updateAdaptorFile struct{}
type deleteAdaptorFile struct{}
type getAdaptorFile struct{}

func init() {
	funcMap["createAdaptor"] = &createAdaptor{}
	funcMap["updateAdaptor"] = &updateAdaptor{}
	funcMap["deleteAdaptor"] = &deleteAdaptor{}
	funcMap["getAdaptor"] = &getAdaptor{}
	funcMap["deployAdaptor"] = &deployAdaptor{}
	funcMap["controlAdaptor"] = &controlAdaptor{}
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
	return adminClient.GetAdaptor(sysKey, adaptorName)
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
	return adminClient.UpdateAdaptor(sysKey, adaptorName, adaptorChanges)
}

func (ct *deleteAdaptor) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 1, ""); err != nil {
		return nil, err
	}
	adaptorName := args[0].(string)
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	return nil, adminClient.DeleteAdaptor(sysKey, adaptorName)
}

func (gd *getAdaptorFile) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 2, "", ""); err != nil {
		return nil, err
	}
	adaptorName := args[0].(string)
	adaptorFileName := args[1].(string)
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	return adminClient.GetAdaptorFile(sysKey, adaptorName, adaptorFileName)
}

func (ct *createAdaptorFile) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 4, "", "", []byte{}, map[string]interface{}{}); err != nil {
		return nil, err
	}
	adaptorName := args[0].(string)
	adaptorFileName := args[1].(string)
	adaptorInput := args[3].(map[string]interface{})
	adaptorInput["file"] = args[2].([]byte)
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	return adminClient.CreateAdaptorFile(sysKey, adaptorName, adaptorFileName, adaptorInput)
}

func (ct *updateAdaptorFile) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 3, "", "", map[string]interface{}{}); err != nil {
		return nil, err
	}
	adaptorName := args[0].(string)
	adaptorFileName := args[1].(string)
	adaptorChanges := args[2].(map[string]interface{})
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	return adminClient.UpdateAdaptorFile(sysKey, adaptorName, adaptorFileName, adaptorChanges)
}

func (ct *deleteAdaptorFile) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 2, "", ""); err != nil {
		return nil, err
	}
	adaptorName := args[0].(string)
	adaptorFileName := args[1].(string)
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	return nil, adminClient.DeleteAdaptorFile(sysKey, adaptorName, adaptorFileName)
}

func (da *deployAdaptor) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 2, "", map[string]interface{}{}); err != nil {
		return nil, err
	}
	adaptorName := args[0].(string)
	deploySpec := args[1].(map[string]interface{})
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	return adminClient.DeployAdaptor(sysKey, adaptorName, deploySpec)
}

func (da *controlAdaptor) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 2, "", map[string]interface{}{}); err != nil {
		return nil, err
	}
	adaptorName := args[0].(string)
	controlSpec := args[1].(map[string]interface{})
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
	sysKey := scriptVars["systemKey"].(string)
	return adminClient.ControlAdaptor(sysKey, adaptorName, controlSpec)
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

func (da *deployAdaptor) help() string {
	return "[\"deployAdaptor\", \"adaptorName\", {deploySpec}]"
}

func (da *controlAdaptor) help() string {
	return "[\"controlAdaptor\", \"adaptorName\", {controlSpec}]"
}
