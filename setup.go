package main

import (
	"encoding/json"
	"fmt"
	cb "github.com/clearblade/Go-SDK"
	"os"
	"strings"
	"time"
)

var (
	devEmail    string
	devPassword string
	adminClient *cb.DevClient
	sysKey      string
	sysSec      string
	setupState  map[string]interface{}
)

func init() {
	setupState = map[string]interface{}{}
	setupState["collections"] = []string{}
	setupState["services"] = []string{}
	setupState["libraries"] = []string{}
	setupState["triggers"] = []string{}
	setupState["timers"] = []string{}
	setupState["roles"] = []string{}
	setupState["users"] = []string{}
}

func performSetup(setupInfo interface{}) {
	//  We're passed either an array of systems, or just one system.
	//  Thus, the type checking stuff.
	//fmt.Printf("THE WHOLE SHABANG IS: %+v\n", setupInfo)
	switch setupInfo.(type) {
	case map[string]interface{}:
		setupSystem(setupInfo.(map[string]interface{}))
	case []interface{}:
		setupSystems(setupInfo.([]interface{}))
	default:
		fmt.Printf("Incorrect type of outer json object")
		os.Exit(1)
	}
	saveSetupState(setupInfo.(map[string]interface{}))

	marshalled, err := json.MarshalIndent(scriptVars, "", "    ")
	if err != nil {
		fatal(fmt.Sprintf("Could not marshal: %s\n", err.Error()))
	}
	fmt.Printf("HERE'S THE STUFF: %s\n", string(marshalled))
}

func setupSystems(systems []interface{}) {
	for _, system := range systems {
		setupSystem(system.(map[string]interface{}))
	}
}

func setupSystem(system map[string]interface{}) {
	if dev, ok := system["developer"]; ok {
		setupDeveloper(dev.(map[string]interface{}))
	} else {
		fmt.Printf("Must provide a developer for system setup\n")
		os.Exit(1)
	}

	createSystem(system)

	if roles, ok := system["roles"]; ok {
		setupRoles(roles.([]interface{}))
	} else {
		warn("No roles found")
	}

	if users, ok := system["users"]; ok {
		setupUsers(users.([]interface{}))
	} else {
		warn("No users found")
	}

	if collections, ok := system["collections"]; ok {
		setupCollections(collections.([]interface{}))
	} else {
		warn("No collections found")
	}

	if codeServices, ok := system["codeServices"]; ok {
		setupCodeServices(codeServices.([]interface{}))
	} else {
		warn("No code services found")
	}

	if codeLibraries, ok := system["codeLibraries"]; ok {
		setupCodeLibraries(codeLibraries.([]interface{}))
	} else {
		warn("No code libraries found")
	}

	if subscriptions, ok := system["subscriptions"]; ok {
		setupSubscriptions(subscriptions.([]interface{}))
	} else {
		warn("No subscriptions found")
	}

	if triggers, ok := system["triggers"]; ok {
		setupTriggers(triggers.([]interface{}))
	} else {
		warn("No triggers found")
	}

	if timers, ok := system["timers"]; ok {
		setupTimers(timers.([]interface{}))
	} else {
		warn("No timers found")
	}
}

func setupDeveloper(dev map[string]interface{}) {
	if _, ok := dev["email"]; !ok {
		fatal("Missing developer email field")
	}
	if _, ok := dev["password"]; !ok {
		fatal("Missing developer password field")
	}
	devEmail = dev["email"].(string)
	devPassword = dev["password"].(string)
	adminClient = cb.NewDevClient(devEmail, devPassword)

	fname := dev["firstname"].(string)
	lname := dev["lastname"].(string)
	org := dev["org"].(string)

	if theDev, err := adminClient.RegisterUser(devEmail, devPassword, fname, lname, org); err == nil {
		setupState["developer"] = theDev["user_id"]
		setupState["dev_email"] = devEmail
		setupState["dev_password"] = devPassword
		return
	} else if strings.Contains(err.Error(), "That user already exists") {
		if authErr := adminClient.Authenticate(); authErr != nil {
			fatal(authErr.Error())
		}
	} else {
		fatal(err.Error())
	}
}

func createSystem(system map[string]interface{}) {
	name := system["name"].(string)
	userAuth := system["userAuth"].(bool)
	descr := system["description"].(string)

	sysStr, sysErr := adminClient.NewSystem(name, descr, userAuth)
	if sysErr != nil {
		fatal(sysErr.Error())
	}
	realSystem, getErr := adminClient.GetSystem(sysStr)
	if getErr != nil {
		fatal(getErr.Error())
	}
	scriptVars["systemKey"] = realSystem.Key
	scriptVars["systemSecret"] = realSystem.Secret
	sysKey = realSystem.Key
	sysSec = realSystem.Secret
	setupState["systemKey"] = sysKey
	setupState["systemSecret"] = sysSec
	// might set more vars...
}

func setupRoles(roles []interface{}) {
	rolesMap := scriptVars["roles"].(map[string]interface{})
	rolesMap["Authenticated"] = "Authenticated"
	rolesMap["Anonymous"] = "Anonymous"
	rolesMap["Admin"] = "Admin"
	for _, role := range roles {
		res, err := adminClient.CreateRole(sysKey, role.(string))
		if err != nil {
			fatal(err.Error())
		}
		fmt.Printf("Created Role: %+v\n", res)
		rolesMap[role.(string)] = res.(map[string]interface{})["role_id"]
		appendState("roles", rolesMap[role.(string)].(string))
	}
}

func setupUsers(users []interface{}) {
	for _, user := range users {
		setupUser(user.(map[string]interface{}))
	}
}

func setupUser(user map[string]interface{}) {
	email := user["email"].(string)
	password := user["password"].(string)
	userClient := cb.NewUserClient(sysKey, sysSec, email, password)
	newUser, err := userClient.RegisterUser(email, password)
	if err != nil {
		fatal(err.Error())
	}

	addUserToRoles(user, newUser["user_id"].(string))

	fmt.Printf("Set up user %+v\n", newUser)
	usersMap := scriptVars["users"].(map[string]interface{})
	newUser["password"] = password
	usersMap[email] = newUser
	appendState("users", newUser["user_id"].(string))
}

func addUserToRoles(user map[string]interface{}, userId string) {
	if _, ok := user["roles"]; !ok {
		warn(fmt.Sprintf("No roles found for %s\n", user["email"].(string)))
		return
	}
	roleNames := user["roles"].([]interface{})
	roleIds := []string{}
	roleMap := scriptVars["roles"].(map[string]interface{})
	for i, _ := range roleNames {
		roleName := roleNames[i].(string)
		if roleId, ok := roleMap[roleName]; ok {
			roleIds = append(roleIds, roleId.(string))
		} else {
			fatal(fmt.Sprintf("Undefined role: %s\n", roleName))
		}
	}
	if len(roleIds) == 0 {
		return
	}
	err := adminClient.AddUserToRoles(sysKey, userId, roleIds)
	if err != nil {
		fatal(err.Error())
	}
}

func setupCollections(cols []interface{}) {
	for _, col := range cols {
		setupCollection(col.(map[string]interface{}))
	}
}

func setupCollection(col map[string]interface{}) {
	fmt.Printf("Setting up collection %+v\n", col["name"])
	//  Create the collection
	colId, err := adminClient.NewCollection(sysKey, col["name"].(string))
	if err != nil {
		fatal(err.Error())
	}
	colsMap := scriptVars["collections"].(map[string]interface{})
	colsMap[col["name"].(string)] = colId
	appendState("collections", colId)

	// Add the columns (this is silly) one at a time
	if _, ok := col["columns"]; !ok {
		fatal(fmt.Sprintf("No columns found for collection %s\n", col["name"].(string)))
	}
	columns := col["columns"].(map[string]interface{})
	for colName, colType := range columns {
		setupColumn(colId, colName, colType.(string))
	}

	setupCollectionRoles(colId, col)

	// Now, add the items
	if _, ok := col["items"]; !ok {
		return
	}
	items := col["items"].([]interface{})
	allData := []map[string]interface{}{}
	for _, item := range items {
		theItem := item.(map[string]interface{})
		count := getCount(theItem)
		delete(theItem, "count")
		for ; count > 0; count-- {
			allData = append(allData, theItem)
		}
	}
	setupItem(allData, colId)
}

func setupCollectionRoles(colId string, col map[string]interface{}) {
	if _, ok := col["roles"]; !ok {
		warn(fmt.Sprintf("No roles found for collection %s\n", col["name"].(string)))
		return
	}
	roles := col["roles"].(map[string]interface{})

	roleDict := scriptVars["roles"].(map[string]interface{})
	for roleName, level := range roles {
		roleId, ok := roleDict[roleName]
		if !ok {
			fatal(fmt.Sprintf("Unknown role: %s\n", roleName))
		}
		err := adminClient.AddCollectionToRole(sysKey, colId, roleId.(string), int(level.(float64)))
		if err != nil {
			fatal(fmt.Sprintf("Could not add collection to role: %s\n", err.Error()))
		}
	}
}

func setupCodeServiceRoles(svc map[string]interface{}) {
	if _, ok := svc["roles"]; !ok {
		warn(fmt.Sprintf("No roles found for collection %s\n", svc["name"].(string)))
		return
	}
	roles := svc["roles"].(map[string]interface{})

	roleDict := scriptVars["roles"].(map[string]interface{})
	for roleName, level := range roles {
		roleId, ok := roleDict[roleName]
		if !ok {
			fatal(fmt.Sprintf("Unknown role: %s\n", roleName))
		}
		err := adminClient.AddServiceToRole(sysKey, svc["name"].(string), roleId.(string), int(level.(float64)))
		if err != nil {
			fatal(fmt.Sprintf("Could not add collection to role: %s\n", err.Error()))
		}
	}
}

func setupColumn(collectionId, columnName, columnType string) {
	fmt.Printf("Adding column %s(%s)\n", columnName, columnType)
	if err := adminClient.AddColumn(collectionId, strings.ToLower(columnName), columnType); err != nil {
		fatal(err.Error())
	}
}

func setupItem(items []map[string]interface{}, colId string) {
	newItems, err := adminClient.CreateData(colId, items)
	if err != nil {
		fatal(fmt.Sprintf("Error creating item: %s\n", err.Error()))
	}
	fmt.Printf("Created item(s): %+v\n", newItems)
}

func setupCodeServices(svcs []interface{}) {
	for _, svc := range svcs {
		setupCodeService(svc.(map[string]interface{}))
	}
}

func setupCodeService(svc map[string]interface{}) {
	svcName := getString(svc, "name")
	svcCode := getVarOrFile(svc, "code")
	svcCode = strings.Replace(svcCode, "\n", "", -1)
	svcParams := svc["parameters"].([]interface{})
	svcDeps := ""
	if _, ok := svc["dependencies"]; ok {
		svcDeps = svc["dependencies"].(string)
	}
	if err := adminClient.NewServiceWithLibraries(sysKey, svcName, svcCode, svcDeps, svcParams); err != nil {
		fatal(err.Error())
	}
	if err := adminClient.EnableLogsForService(sysKey, svcName); err != nil {
		fatal(err.Error())
	}
	setupCodeServiceRoles(svc)
	svcMap := scriptVars["codeServices"].(map[string]interface{})
	svcMap[svcName] = map[string]interface{}{
		"name":         svcName,
		"code":         svcCode,
		"params":       svcParams,
		"dependencies": svcDeps,
	}
	appendState("services", svcName)
	fmt.Printf("Set up code service %+v\n", svcMap[svcName])
}

func setupCodeLibraries(libs []interface{}) {
	for _, lib := range libs {
		setupCodeLibrary(lib.(map[string]interface{}))
	}
}

func setupCodeLibrary(lib map[string]interface{}) {
	fmt.Printf("Setting up code library %+v\n", lib["name"])
	libName := lib["name"].(string)
	delete(lib, "name")
	libCode := getVarOrFile(lib, "code")
	lib["code"] = libCode
	newLib, err := adminClient.CreateLibrary(sysKey, libName, lib)
	if err != nil {
		fatal(fmt.Sprintf("Could not create code library: %s\n", err.Error()))
	}
	libMap := scriptVars["codeLibraries"].(map[string]interface{})
	libMap[libName] = newLib
	appendState("libraries", libName)
	fmt.Printf("Set up code library %+v\n", newLib)
}

func setupSubscriptions(subs []interface{}) {
	for _, sub := range subs {
		setupSubscription(sub.(map[string]interface{}))
	}
}

func setupSubscription(sub map[string]interface{}) {
	fmt.Printf("Setting up subscription for %+v\n", sub["topic"])
}

func setupTriggers(triggers []interface{}) {
	for _, trigger := range triggers {
		setupTrigger(trigger.(map[string]interface{}))
	}
}

func setupTrigger(trigger map[string]interface{}) {
	trigName := trigger["name"].(string)
	delete(trigger, "name")
	newTrig, err := adminClient.CreateEventHandler(sysKey, trigName, trigger)
	if err != nil {
		fatal(fmt.Sprintf("Could not create trigger: %s\n", err.Error()))
	}
	trigMap := scriptVars["triggers"].(map[string]interface{})
	trigMap[trigName] = newTrig
	appendState("triggers", trigName)
	fmt.Printf("Set up trigger %+v\n", newTrig)
}

func setupTimers(timers []interface{}) {
	for _, timer := range timers {
		setupTimer(timer.(map[string]interface{}))
	}
}

func setupTimer(timer map[string]interface{}) {
	timerName := timer["name"].(string)
	delete(timer, "name")
	startTime := timer["start_time"].(string)
	if startTime == "Now" {
		startTime = time.Now().Format(time.RFC3339)
	}
	timer["start_time"] = startTime
	newTimer, err := adminClient.CreateTimer(sysKey, timerName, timer)
	if err != nil {
		fatal(fmt.Sprintf("Could not create timer: %s\n", err.Error()))
	}
	timerMap := scriptVars["timers"].(map[string]interface{})
	timerMap[timerName] = newTimer
	appendState("timers", timerName)
	fmt.Printf("Set up timer %+v\n", newTimer)
}

func warn(msg string) {
	fmt.Printf("Warning: %s\n", msg)
}

func fatal(msg string) {
	fmt.Printf("Fatal Error: %s\n", msg)
	os.Exit(1)
}

func getCount(stuff map[string]interface{}) int {
	if _, ok := stuff["count"]; !ok {
		return 1
	}
	return int(stuff["count"].(float64))
}

func getString(stuff map[string]interface{}, theThing string) string {
	if _, ok := stuff[theThing]; !ok {
		fatal(fmt.Sprintf("The value for %s does not exist\n", theThing))
	}
	switch stuff[theThing].(type) {
	case string:
	default:
		fatal(fmt.Sprintf("The value of %s is not a string\n", theThing))
	}
	return stuff[theThing].(string)
}

func appendState(stateKey, value string) {
	theList := setupState[stateKey].([]string)
	theList = append(theList, value)
	setupState[stateKey] = theList
}
