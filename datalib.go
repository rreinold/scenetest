package main

import (
	"fmt"
	cb "github.com/clearblade/Go-SDK"
)

func init() {
	funcMap["query"] = query
	funcMap["createItem"] = createItem
}

func createItem(context map[string]interface{}, args []interface{}) error {
	userClient := context["userClient"].(*cb.UserClient)
	if len(args) != 2 {
		return fmt.Errorf("Usage: [createItem, <colName>, <argMap>]")
	}
	colName := args[0].(string)
	rowInfo := args[1].(map[string]interface{})
	colId, err := collectionNameToId(colName)
	if err != nil {
		return err
	}
	if err = userClient.InsertData(colId, rowInfo); err != nil {
		return err
	}
	return nil
}

func query(context map[string]interface{}, args []interface{}) error {
	return nil
}

func collectionNameToId(colName string) (string, error) {
	cols := scriptVars["collections"].(map[string]interface{})
	if colId, ok := cols[colName]; ok {
		return colId.(string), nil
	}
	return "", fmt.Errorf("Could not find collection %s", colName)
}
