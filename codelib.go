package main

func init() {
	funcMap["runService"] = runService
}

func runService(context map[string]interface{}, args []interface{}) error {
	return nil
}
