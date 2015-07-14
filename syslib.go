package main

import (
	//"fmt"
	"time"
)

func init() {
	funcMap["sleep"] = sleep
}

func sleep(ctx map[string]interface{}, args []interface{}) error {
	secs := time.Duration(getArg(args, 0).(float64))
	time.Sleep(secs * time.Second)
	return nil
}
