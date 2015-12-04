package main

import (
	"clearblade/token"
	"fmt"
	cb "github.com/clearblade/Go-SDK"
)

func init() {
	funcMap["setUser"] = &Statement{setUser, setUserHelp}
}

func setUser(ctx map[string]interface{}, args []interface{}) error {
	sysKey := scriptVars["systemKey"].(string)
	sysSec := scriptVars["systemSecret"].(string)
	email := valueOf(ctx, getArg(args, 0)).(string)
	userInfo := scriptVars["users"].(map[string]interface{})[email].(map[string]interface{})
	password := userInfo["password"].(string)
	userClient := cb.NewUserClient(sysKey, sysSec, email, password)
	if err := userClient.Authenticate(); err != nil {
		return err
	}

	sess, _ := token.Token(userClient.UserToken).Uuid()
	fmt.Printf("AUTHENTICATED SESSION: %s\n", sess)

	// Now, might as well set up mqtt
	if err := userClient.InitializeMQTT("", "", 60); err != nil {
		return err
	}
	if err := userClient.ConnectMQTT(nil, nil); err != nil {
		return err
	}

	/*
		triggerChan, err := userClient.Subscribe("/clearblade/internal/trigger", 0)
		if err != nil {
			return fmt.Errorf("Subscribe failed: %s", err.Error())
		}
	*/
	ctx["userClient"] = userClient
	if err := userClient.Publish("/who/am/i", []byte(fmt.Sprintf("%p", userClient.MQTTClient)), 2); err != nil {
		return err
	}
	ctx["email"] = email
	//ctx["triggerChannel"] = triggerChan

	return nil
}

func setUserHelp() string {
	return "setUser help not yet implemented"
}

func getArg(args []interface{}, index int) interface{} {
	if index >= len(args) {
		fatal("Attempt to get non-existent argument")
	}
	return args[index]
}
