package main

import (
	//"clearblade/token"
	"fmt"
	cb "github.com/clearblade/Go-SDK"
)

type setUserStmt struct{}
type createUserStmt struct{}
type updateUserStmt struct{}
type deleteUserStmt struct{}
type createUserColumnStmt struct{}
type getUserColumnsStmt struct{}

func init() {
	funcMap["setUser"] = &setUserStmt{}
	funcMap["createUser"] = &createUserStmt{}
	funcMap["updateUser"] = &updateUserStmt{}
	funcMap["deleteUser"] = &deleteUserStmt{}
	funcMap["createUserColumn"] = &createUserColumnStmt{}
	funcMap["getUserColumns"] = &getUserColumnsStmt{}
}

func (s *setUserStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	sysKey := scriptVars["systemKey"].(string)
	sysSec := scriptVars["systemSecret"].(string)
	email := valueOf(ctx, getArg(args, 0)).(string)
	userInfo := scriptVars["users"].(map[string]interface{})[email].(map[string]interface{})
	password := userInfo["password"].(string)
	userClient := cb.NewUserClient(sysKey, sysSec, email, password)
	if err := userClient.Authenticate(); err != nil {
		return nil, err
	}

	//sess, _ := token.Token(userClient.UserToken).Uuid()

	// Now, might as well set up mqtt
	if err := userClient.InitializeMQTT("", "", 60); err != nil {
		return nil, err
	}
	if err := userClient.ConnectMQTT(nil, nil); err != nil {
		return nil, err
	}

	ctx["userClient"] = userClient
	if err := userClient.Publish("/who/am/i", []byte(fmt.Sprintf("%p", userClient.MQTTClient)), 2); err != nil {
		return nil, err
	}
	ctx["email"] = email
	//ctx["triggerChannel"] = triggerChan

	return nil, nil
}

func (s *setUserStmt) help() string {
	return "setUser help not yet implemented"
}

func (c *createUserStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	devCli := ctx["adminClient"].(*cb.DevClient)
	if len(args) != 2 {
		return nil, fmt.Errorf(c.help())
	}
	if _, ok := args[0].(string); !ok {
		return nil, fmt.Errorf("First arg to createUser must be a string")
	}
	if _, ok := args[1].(string); !ok {
		return nil, fmt.Errorf("Second arg to createUser must be a string")
	}
	email := args[0].(string)
	password := args[1].(string)
	sysKey := scriptVars["systemKey"].(string)
	sysSec := scriptVars["systemSecret"].(string)
	newUser, err := devCli.RegisterNewUser(email, password, sysKey, sysSec)
	if err != nil {
		return nil, fmt.Errorf("Could not create user: %s", err.Error())
	}
	return newUser["user_id"], nil
}

func (c *createUserStmt) help() string {
	return "createUser help not yet implemented"
}

func (u *updateUserStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	return nil, nil
}

func (u *updateUserStmt) help() string {
	return "updateUser help not yet implemented"
}

func (d *deleteUserStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	devCli := ctx["adminClient"].(*cb.DevClient)
	if len(args) != 1 {
		return nil, fmt.Errorf(d.help())
	}
	if _, ok := args[0].(string); !ok {
		return nil, fmt.Errorf("First arg to deleteUser must be a string")
	}
	userId := args[0].(string)
	sysKey := scriptVars["systemKey"].(string)
	if err := devCli.DeleteUser(sysKey, userId); err != nil {
		return nil, fmt.Errorf("Could not delete user: %s", err.Error())
	}
	return nil, nil
}

func (d *deleteUserStmt) help() string {
	return "deleteUser help not yet implemented"
}

func (c *createUserColumnStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	devCli := ctx["adminClient"].(*cb.DevClient)
	if len(args) != 2 {
		return nil, fmt.Errorf(c.help())
	}
	if _, ok := args[0].(string); !ok {
		return nil, fmt.Errorf("First arg to createUserColumn must be a string column name")
	}
	if _, ok := args[1].(string); !ok {
		return nil, fmt.Errorf("Second arg to createUserColumn must be a string column type")
	}
	colName := args[0].(string)
	colType := args[1].(string)
	sysKey := scriptVars["systemKey"].(string)
	if err := devCli.CreateUserColumn(sysKey, colName, colType); err != nil {
		return nil, fmt.Errorf("Could not create users table column: %s", err.Error())
	}
	return nil, nil
}

func (c *createUserColumnStmt) help() string {
	return "[\"createUserColumn\", <colName>, <colType>]"
}

func (g *getUserColumnsStmt) run(ctx map[string]interface{}, args []interface{}) (interface{}, error) {
	devCli := ctx["adminClient"].(*cb.DevClient)
	if len(args) != 0 {
		return nil, fmt.Errorf(g.help())
	}
	sysKey := scriptVars["systemKey"].(string)
	cols, err := devCli.GetUserColumns(sysKey)
	if err != nil {
		return nil, fmt.Errorf("Could not get users table columns: %s", err.Error())
	}
	return cols, nil
}

func (g *getUserColumnsStmt) help() string {
	return "[\"getUserColumns\"]"
}

func getArg(args []interface{}, index int) interface{} {
	if index >= len(args) {
		fatal("Attempt to get non-existent argument")
	}
	return args[index]
}
