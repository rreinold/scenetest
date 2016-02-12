package main

import (
	//"clearblade/token"
	"fmt"
	cb "github.com/clearblade/Go-SDK"
)

func init() {
	funcMap["query"] = &Statement{query, queryHelp}
	funcMap["select"] = &Statement{query, queryHelp} // just a synonym
	funcMap["createItem"] = &Statement{createItem, createItemHelp}
	funcMap["deleteItem"] = &Statement{deleteItem, deleteItemHelp}
	funcMap["deleteAllItems"] = &Statement{deleteAllItems, deleteAllItemsHelp}
}

func createItem(context map[string]interface{}, args []interface{}) error {
	userClient := context["userClient"].(*cb.UserClient)
	if len(args) != 2 {
		return fmt.Errorf("Usage: [createItem, <colName>, <argMap>]")
	}
	colName := valueOf(context, args[0]).(string)
	rowInfo := valueOf(context, args[1]).(map[string]interface{})
	colId, err := collectionNameToId(colName)
	if err != nil {
		return err
	}
	//sess, _ := token.Token(userClient.UserToken).Uuid()
	resp, err := userClient.CreateData(colId, rowInfo)
	if err != nil {
		return err
	}
	if len(resp) != 1 {
		return fmt.Errorf("Wrong number if items returned by createItem: %d", len(resp))
	}
	val := resp[0].(map[string]interface{})
	context["returnValue"] = val["item_id"]
	return nil
}

func deleteItem(context map[string]interface{}, args []interface{}) error {
	if err := argCheck(args, 2, "", ""); err != nil {
		return fmt.Errorf("query: Bad argument(s): %s", err.Error())
	}
	colName := args[0].(string)
	colId, err := collectionNameToId(colName)
	if err != nil {
		return err
	}
	itemId := args[1].(string)
	userClient := context["userClient"].(*cb.UserClient)

	query := &cb.Query{
		Filters: [][]cb.Filter{
			[]cb.Filter{
				cb.Filter{
					Field:    "item_id",
					Operator: "=",
					Value:    itemId,
				},
			},
		},
	}

	if err := userClient.DeleteData(colId, query); err != nil {
		context["returnValue"] = err.Error()
		return err
	}
	context["returnValue"] = nil
	return nil
}

func deleteAllItems(context map[string]interface{}, args []interface{}) error {
	userClient := context["userClient"].(*cb.UserClient)
	if err := argCheck(args, 1, ""); err != nil {
		return fmt.Errorf("deleteAll: Bad Arguments(s): %s\n", err.Error())
	}
	colName := args[0].(string)
	colId, err := collectionNameToId(colName)
	if err != nil {
		return err
	}
	if err := userClient.DeleteData(colId, &cb.Query{}); err != nil {
		context["returnValue"] = err.Error()
		return err
	}
	context["returnValue"] = nil
	return nil
}

//
//  Query format is:
//  ["query|select", "colName", [<columns>], [<optFilters>],
//			[<optOrderBy>], <optPageNumber>, <optPageSize>]
//
//  Where each of the <>'s are:
//		columns: array of strings, empty means all columns.
//		filter[[["field", "op", val],["a","==","b"]],[["j","!=","q"]]]
//      ^^^^^^^^This means ((field op val && a == b) || j != q)
//      orderby: opt -- ["fred", "desc", "flint", "asc"] -- [] == no ordering
//      page num and page size: ints
//
//
//  Examples:
//		["query", "theCollection"]
//			-- gets everything from theCollection
//		["query", "theCol", ["Age"]]
//			-- gets all the Age column values
//      ["query", "theCol", [], [[["id","==","xyz"]]]
//			-- gets all rows where id eq xyz

func query(context map[string]interface{}, args []interface{}) error {
	if err := argCheck(args, 1, "", []interface{}{}, []interface{}{}, []interface{}{}, 1, 1); err != nil {
		return fmt.Errorf("query: Bad argument(s): %s", err.Error())
	}
	myQuery := cb.Query{}
	var err error
	collectionName := args[0].(string)
	collection, err := collectionNameToId(collectionName)
	if err != nil {
		return err
	}
	if len(args) >= 2 {
		myQuery.Columns, err = buildColumns(valueOf(context, args[1]).([]interface{}))
	}
	if len(args) >= 3 {
		myQuery.Filters, err = buildFilter(valueOf(context, args[2]).([]interface{}))
		if err != nil {
			return fmt.Errorf("query: Bad filter: %s", err.Error())
		}
	}
	if len(args) >= 4 {
		if myQuery.Order, err = buildOrdering(args[3].([]interface{})); err != nil {
			return fmt.Errorf("Bad query ordering: %s", err.Error())
		}
	}
	if len(args) >= 5 {
		myQuery.PageNumber = int(args[4].(float64))
	}
	if len(args) >= 6 {
		myQuery.PageSize = int(args[5].(float64))
	}
	userClient := context["userClient"].(*cb.UserClient)
	stuff, err := userClient.GetData(collection, &myQuery)
	if err != nil {
		return err
	}
	context["returnValue"] = stuff["DATA"]
	context["returnCount"] = stuff["TOTAL"]

	return nil
}

func collectionNameToId(colName string) (string, error) {
	cols := scriptVars["collections"].(map[string]interface{})
	if colId, ok := cols[colName]; ok {
		return colId.(string), nil
	}
	return "", fmt.Errorf("Could not find collection %s", colName)
}

//  This is absurd
func buildFilter(f []interface{}) ([][]cb.Filter, error) {
	rval := [][]cb.Filter{}
	for orIdx, pitter := range f {
		rval = append(rval, []cb.Filter{})
		orList := pitter.([]interface{})
		for _, patter := range orList {
			filterArray := patter.([]interface{})
			filter, err := makeFilter(filterArray)
			if err != nil {
				return nil, err
			}
			rval[orIdx] = append(rval[orIdx], *filter)
		}
	}
	return rval, nil
}

func makeFilter(stuff []interface{}) (*cb.Filter, error) {
	if len(stuff) != 3 {
		return nil, fmt.Errorf("Each filter needs the values")
	}
	field := stuff[0].(string)
	operator := stuff[1].(string)
	value := stuff[2]
	rval := cb.Filter{
		Field:    field,
		Operator: operator,
		Value:    value,
	}

	return &rval, nil
}

// More absurdity
func buildColumns(cols []interface{}) ([]string, error) {
	rval := []string{}
	for _, col := range cols {
		switch col.(type) {
		case string:
			rval = append(rval, col.(string))
		default:
			return nil, fmt.Errorf("All columns must be strings")
		}
	}
	return rval, nil
}

func buildOrdering(ordering []interface{}) ([]cb.Ordering, error) {
	if len(ordering) == 0 {
		return nil, nil
	}
	if len(ordering)%2 != 0 {
		return nil, fmt.Errorf("Ordering arguments must be pairs")
	}
	rval := []cb.Ordering{}
	for i := 0; i < len(ordering); i += 2 {
		field := ordering[i].(string)
		direction := ordering[i+1].(string)
		dirBool := false
		if direction == "asc" || direction == "ascend" || direction == "ascending" {
			dirBool = true
		}
		rval = append(rval, cb.Ordering{SortOrder: dirBool, OrderKey: field})
	}
	return rval, nil
}

func deleteAllItemsHelp() string {
	return "query help not yet implemented"
}

func queryHelp() string {
	return "query help not yet implemented"
}

func createItemHelp() string {
	return "createItem help not yet implemented"
}

func deleteItemHelp() string {
	return "deleteItem help not yet implemented"
}
