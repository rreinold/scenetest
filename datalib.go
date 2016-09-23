package main

import (
	"fmt"
	cb "github.com/clearblade/Go-SDK"
)

type queryStmt struct{}
type createItemStmt struct{}
type deleteItemStmt struct{}
type deleteAllItemsStmt struct{}
type createCollectionStmt struct{}

//type createConnectCollectionStmt struct{}
type deleteCollectionStmt struct{}
type allCollectionsStmt struct{}

func init() {
	funcMap["query"] = &queryStmt{}
	funcMap["createItem"] = &createItemStmt{}
	funcMap["deleteItem"] = &deleteItemStmt{}
	funcMap["deleteAllItems"] = &deleteAllItemsStmt{}
	funcMap["createCollection"] = &createCollectionStmt{}
	//funcMap["createConnectCollection"] = &createConnectCollectionStmt{}
	funcMap["deleteCollection"] = &deleteCollectionStmt{}
	funcMap["allCollections"] = &allCollectionsStmt{}
}

func (c *createItemStmt) run(context map[string]interface{}, args []interface{}) (interface{}, error) {
	userClient := context["userClient"].(cb.Client)
	if err := argCheck(args, 2, "", map[string]interface{}{}); err != nil {
		return nil, fmt.Errorf("Usage: [createItem, <colName>, <argMap>]")
	}
	colName := args[0].(string)
	rowInfo := args[1].(map[string]interface{})
	colId, err := collectionNameToId(colName)
	if err != nil {
		return nil, err
	}
	resp, err := userClient.CreateData(colId, rowInfo)
	if err != nil {
		return nil, err
	}
	if len(resp) != 1 {
		return nil, fmt.Errorf("Wrong number if items returned by createItem: %d", len(resp))
	}
	val := resp[0].(map[string]interface{})
	return val["item_id"], nil
}

func (d *deleteItemStmt) run(context map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 2, "", ""); err != nil {
		return nil, fmt.Errorf("query: Bad argument(s): %s", err.Error())
	}
	colName := args[0].(string)
	colId, err := collectionNameToId(colName)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return nil, nil
}

func (d *deleteAllItemsStmt) run(context map[string]interface{}, args []interface{}) (interface{}, error) {
	userClient := context["userClient"].(*cb.UserClient)
	if err := argCheck(args, 1, ""); err != nil {
		return nil, fmt.Errorf("deleteAll: Bad Arguments(s): %s\n", err.Error())
	}
	colName := args[0].(string)
	colId, err := collectionNameToId(colName)
	if err != nil {
		return nil, err
	}
	if err := userClient.DeleteData(colId, &cb.Query{}); err != nil {
		context["returnValue"] = err.Error()
		return nil, err
	}
	return nil, nil
}

func (c *createCollectionStmt) run(context map[string]interface{}, args []interface{}) (interface{}, error) {
	scriptVarsLock.Lock()
	defer scriptVarsLock.Unlock()
	if err := argCheck(args, 1, ""); err != nil {
		return nil, err
	}
	colName := args[0].(string)
	adminClient := context["adminClient"].(*cb.DevClient)
	systemKey := scriptVars["systemKey"].(string)
	rval, err := adminClient.NewCollection(systemKey, colName)
	if err != nil {
		return nil, err
	}
	allCollections := scriptVars["collections"].(map[string]interface{})
	allCollections[colName] = rval
	return rval, nil
}

func (c *createCollectionStmt) help() string {
	return "[\"createCollection\", \"<collectionName>\"]"
}

/*
func (c *createConnectCollectionStmt) run(context map[string]interface{}, args []interface{}) (interface{}, error) {
	scriptVarsLock.Lock()
	defer scriptVarsLock.Unlock()
	if err := argCheck(args, 1, map[string]interface{}{}); err != nil {
		return nil, err
	}
	collectionConfig := args[0].(map[string]interface{})
	adminClient := context["adminClient"].(*cb.DevClient)
	systemKey := scriptVars["systemKey"].(string)
	rval, err := adminClient.NewConnectCollection(systemKey, collectionConfig)
	if err != nil {
		return nil, err
	}
	allCollections := scriptVars["collections"].(map[string]interface{})

	//allCollections[colName] = rval
	return rval, nil
}

func (c *createConnectCollectionStmt) help() string {
	return "[\"createConnectCollection\", \"{<collectionConfig>}\"]"
}
*/

func (d *deleteCollectionStmt) run(context map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 1, ""); err != nil {
		return nil, err
	}
	colName := args[0].(string)
	collectionId, err := collectionNameToId(colName)
	if err != nil {
		return nil, err
	}
	adminClient := context["adminClient"].(*cb.DevClient)
	if err := adminClient.DeleteCollection(collectionId); err != nil {
		return nil, err
	}
	return nil, nil
}

func (d *deleteCollectionStmt) help() string {
	return "[\"deleteCollection\", \"<collectionName>\"]"
}

func (a *allCollectionsStmt) run(context map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 0); err != nil {
		return nil, err
	}
	adminClient := context["adminClient"].(*cb.DevClient)
	sysKey := scriptVars["systemKey"].(string)
	cols, err := adminClient.GetAllCollections(sysKey)
	if err != nil {
		return nil, err
	}
	return cols, nil
}

func (a *allCollectionsStmt) help() string {
	return "[\"allCollections\"]"
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

func (q *queryStmt) run(context map[string]interface{}, args []interface{}) (interface{}, error) {
	if err := argCheck(args, 1, "", []interface{}{}, []interface{}{}, []interface{}{}, 1, 1); err != nil {
		return nil, fmt.Errorf("query: Bad argument(s): %s", err.Error())
	}
	myQuery := cb.Query{}
	var err error
	collectionName := args[0].(string)
	collection, err := collectionNameToId(collectionName)
	if err != nil {
		return nil, err
	}
	if len(args) >= 2 {
		myQuery.Columns, err = buildColumns(args[1].([]interface{}))
	}
	if len(args) >= 3 {
		myQuery.Filters, err = buildFilter(args[2].([]interface{}))
		if err != nil {
			return nil, fmt.Errorf("query: Bad filter: %s", err.Error())
		}
	}
	if len(args) >= 4 {
		if myQuery.Order, err = buildOrdering(args[3].([]interface{})); err != nil {
			return nil, fmt.Errorf("Bad query ordering: %s", err.Error())
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
		return nil, err
	}
	context["returnCount"] = stuff["TOTAL"]

	return stuff["DATA"], nil
}

func collectionNameToId(colName string) (string, error) {
	scriptVarsLock.RLock()
	defer scriptVarsLock.RUnlock()
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

func (d *deleteAllItemsStmt) help() string {
	return "query help not yet implemented"
}

func (q *queryStmt) help() string {
	return "query help not yet implemented"
}

func (c *createItemStmt) help() string {
	return "createItem help not yet implemented"
}

func (d *deleteItemStmt) help() string {
	return "deleteItem help not yet implemented"
}
