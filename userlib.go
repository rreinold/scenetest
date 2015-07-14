package main

import (
	"fmt"
	cb "github.com/clearblade/Go-SDK"
)

func init() {
	funcMap["setUser"] = setUser
}

func setUser(ctx map[string]interface{}, args []interface{}) error {
	sysKey := scriptVars["systemKey"].(string)
	sysSec := scriptVars["systemSecret"].(string)
	email := getArg(args, 0).(string)
	userInfo := scriptVars["users"].(map[string]interface{})[email].(map[string]interface{})
	password := userInfo["password"].(string)
	fmt.Printf("K/S: %s/%s\n", sysKey, sysSec)
	userClient := cb.NewUserClient(sysKey, sysSec, email, password)
	if err := userClient.Authenticate(); err != nil {
		return err
	}

	// Now, might as well set up mqtt
	if err := userClient.InitializeMQTT(email, "", 60); err != nil {
		return err
	}
	if err := userClient.ConnectMQTT(nil, nil); err != nil {
		return err
	}

	triggerChan, err := userClient.Subscribe("/clearblade/internal/trigger", 0)
	if err != nil {
		return fmt.Errorf("Subscribe failed: %s", err.Error())
	}
	ctx["userClient"] = userClient
	ctx["email"] = email
	ctx["triggerChannel"] = triggerChan

	return nil
}

func getArg(args []interface{}, index int) interface{} {
	if index >= len(args) {
		fatal("Attempt to get non-existent argument")
	}
	return args[index]
}
