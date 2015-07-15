package main

import (
	"fmt"
	cb "github.com/clearblade/Go-SDK"
)

func init() {
	funcMap["query"] = query
	funcMap["select"] = query // just a synonym
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
			rval[orIdx] = append(rval[orIdx], filter)
		}
	}
	return rval, nil
}

func makeFilter(stuff []interface{}) (cb.Filter, error) {
	rval := cb.Filter{}

	return rval, nil
}
