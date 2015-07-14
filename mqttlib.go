package main

import ()

func init() {
	funcMap["subscribe"] = subscribe
	funcMap["publish"] = publish
}

func subscribe(context map[string]interface{}, args []interface{}) error {
	return nil
}

func publish(context map[string]interface{}, args []interface{}) error {
	return nil
}
