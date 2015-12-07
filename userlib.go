package main

import (
	"clearblade/token"
	"fmt"
	cb "github.com/clearblade/Go-SDK"
)

func init() {
	funcMap["setUser"] = &Statement{setUser, setUserHelp}
	funcMap["createUser"] = &Statement{createUser, createUserHelp}
	funcMap["updateUser"] = &Statement{updateUser, updateUserHelp}
	funcMap["deleteUser"] = &Statement{deleteUser, deleteUserHelp}
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

func createUser(ctx map[string]interface{}, args []interface{}) error {
	devCli := ctx["adminClient"].(*cb.DevClient)
	if len(args) != 2 {
		return fmt.Errorf(createUserHelp())
	}
	if _, ok := args[0].(string); !ok {
		return fmt.Errorf("First arg to createUser must be a string")
	}
	if _, ok := args[1].(string); !ok {
		return fmt.Errorf("Second arg to createUser must be a string")
	}
	email := args[0].(string)
	password := args[1].(string)
	sysKey := scriptVars["systemKey"].(string)
	sysSec := scriptVars["systemSecret"].(string)
	newUser, err := devCli.RegisterNewUser(email, password, sysKey, sysSec)
	if err != nil {
		return fmt.Errorf("Could not create user: %s", err.Error())
	}
	fmt.Printf("GOT NEW USER: %+v\n", newUser)
	return nil
}

func createUserHelp() string {
	return "createUser help not yet implemented"
}

func updateUser(ctx map[string]interface{}, args []interface{}) error {
	return nil
}

func updateUserHelp() string {
	return "updateUser help not yet implemented"
}

func deleteUser(ctx map[string]interface{}, args []interface{}) error {
	return nil
}

func deleteUserHelp() string {
	return "deleteUser help not yet implemented"
}

func getArg(args []interface{}, index int) interface{} {
	if index >= len(args) {
		fatal("Attempt to get non-existent argument")
	}
	return args[index]
}
