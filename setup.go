package main

import (
	"encoding/json"
	"fmt"
	cb "github.com/clearblade/Go-SDK"
	"os"
	"strings"
)

type AddRoleFunc func(string, string, []interface{})

var (
	devEmail        string
	devPassword     string
	sysKey          string
	sysSec          string
	setupState      map[string]interface{}
	globalSetupInfo map[string]interface{}
)

func init() {
	setupState = map[string]interface{}{}
	setupState["collections"] = []string{}
	setupState["connectCollections"] = []string{}
	setupState["services"] = []string{}
	setupState["libraries"] = []string{}
	setupState["triggers"] = []string{}
	setupState["timers"] = []string{}
	setupState["roles"] = []string{}
	setupState["users"] = []string{}
	setupState["developers"] = []string{}
	setupState["devices"] = []string{}
	setupState["edges"] = []string{}
	setupState["portals"] = []string{}
	setupState["systems"] = []string{}
}

func performSetup(setupInfo interface{}) {
	switch setupInfo.(type) {
	case map[string]interface{}:
		globalSetupInfo = setupInfo.(map[string]interface{})
		setupState["edgeDeploy"] = makeEdgeDeployStructure()
		setupState["edgeSync"] = makeEdgeSyncStructure()
		setupSystem(setupInfo.(map[string]interface{}))
	default:
		myPrintf("Incorrect type of outer json object")
		os.Exit(1)
	}
	saveSetupState(setupInfo.(map[string]interface{}))
	fmt.Printf("New System Key is: %s\n", sysKey)
}

func setupSystem(system map[string]interface{}) {
	if dev, ok := system["developer"]; ok {
		setupMainDeveloper(dev.(map[string]interface{}))
	} else {
		myPrintf("Must provide a developer for system setup\n")
		os.Exit(1)
	}

	createSystem(system)

	if roles, ok := system["roles"]; ok {
		setupRoles(roles.([]interface{}))
	} else {
		warn("No roles found")
	}

	if userColumns, ok := system["userColumns"]; ok {
		setupUserColumns(userColumns.([]interface{}))
	} else {
		warn("No user columns found")
	}

	if users, ok := system["users"]; ok {
		setupUsers(users.([]interface{}))
	} else {
		warn("No users found")
	}

	if developers, ok := system["developers"]; ok {
		setupDevelopers(developers.([]interface{}))
	} else {
		warn("No \"extra\" developers found")
	}

	if userTablePerms, ok := system["userTableRoles"]; ok {
		setupGenericPermissions("users", userTablePerms.(map[string]interface{}))
	} else {
		warn("No user table permissions ('userTableRoles') found")
	}

	if deviceTablePerms, ok := system["deviceTableRoles"]; ok {
		setupGenericPermissions("devices", deviceTablePerms.(map[string]interface{}))
	} else {
		warn("No device table permissions ('deviceTableRoles') found")
	}

	if messageHistoryPerms, ok := system["messageHistoryRoles"]; ok {
		setupGenericPermissions("msgHistory", messageHistoryPerms.(map[string]interface{}))
	} else {
		warn("No message history permissions ('messageHistoryRoles') found")
	}

	if collections, ok := system["collections"]; ok {
		setupCollections(collections.([]interface{}))
	} else {
		warn("No collections found")
	}

	if connectCollections, ok := system["connectCollections"]; ok {
		setupConnectCollections(connectCollections.([]interface{}))
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

	if deviceColumns, ok := system["deviceColumns"]; ok {
		setupDeviceColumns(deviceColumns.([]interface{}))
	} else {
		warn("No device columns found")
	}

	if devices, ok := system["devices"]; ok {
		setupDevices(devices.([]interface{}))
	} else {
		warn("No devices found")
	}

	if portals, ok := system["portals"]; ok {
		setupPortals(portals.([]interface{}))
	} else {
		warn("No portals found")
	}

	if edges, ok := system["edges"]; ok {
		setupEdges(edges.([]interface{}))
	} else {
		warn("No devices found")
	}
	setupEdgeDeployInfo()
}

func setupMainDeveloper(dev map[string]interface{}) {
	if _, ok := dev["email"]; !ok {
		fatal("Missing developer email field")
	}
	if _, ok := dev["password"]; !ok {
		fatal("Missing developer password field")
	}
	devEmail = dev["email"].(string)
	devPassword = dev["password"].(string)
	setupState["dev_email"] = devEmail
	setupState["dev_password"] = devPassword
	adminClient = cb.NewDevClient(devEmail, devPassword)

	fname := dev["firstname"].(string)
	lname := dev["lastname"].(string)
	org := dev["org"].(string)

	theDev, err := adminClient.RegisterDevUser(devEmail, devPassword, fname, lname, org)
	if err != nil {
		if !strings.Contains(err.Error(), "That user already exists") {
			fatal(err.Error())
		}
	}
	if authErr := adminClient.Authenticate(); authErr != nil {
		fatal(authErr.Error())
	}
	setupState["developer"] = theDev["user_id"]
	setupState["adminClient"] = adminClient
	scriptVars["developer"] = map[string]interface{}{
		"userId":   theDev["user_id"],
		"email":    devEmail,
		"password": devPassword,
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
		myPrintf("Created Role: %s\n", role.(string))
		rolesMap[role.(string)] = res.(map[string]interface{})["role_id"]
		appendState("roles", rolesMap[role.(string)].(string))
	}
}

func setupUserColumns(userColumns []interface{}) {
	for _, userColumn := range userColumns {
		setupUserColumn(userColumn.(map[string]interface{}))
	}
}

func setupUserColumn(userColumn map[string]interface{}) {
	sysKey := setupState["systemKey"].(string)
	adminClient = setupState["adminClient"].(*cb.DevClient)
	if err := adminClient.CreateUserColumn(sysKey, userColumn["column_name"].(string), userColumn["type"].(string)); err != nil {
		fatal(err.Error())
	}
	myPrintf("Added column to user table: %s\n", userColumn["column_name"].(string))
}

func setupDevelopers(developers []interface{}) {
	for _, developer := range developers {
		setupDeveloper(developer.(map[string]interface{}))
	}
}

func setupDeveloper(developer map[string]interface{}) {
	email := developer["email"].(string)
	pass := developer["password"].(string)
	fname := developer["fname"].(string)
	lname := developer["lname"].(string)
	org := developer["org"].(string)
	newDev, err := adminClient.RegisterDevUser(email, pass, fname, lname, org)
	if err != nil {
		if !strings.Contains(err.Error(), "That user already exists") {
			fatal(err.Error())
		}
		//  just fake the dev id and token for now. Our apis are lacking.
		newDev = map[string]interface{}{}
		newDev["user_id"] = fmt.Sprintf("<unknownId:%s>", email)
		newDev["dev_token"] = fmt.Sprintf("<unknownToken:%s>", email)
	}
	newId := newDev["user_id"].(string)

	devMap := scriptVars["developers"].(map[string]interface{})
	newDev["password"] = pass
	newDev["email"] = email
	newDev["fname"] = fname
	newDev["lname"] = lname
	newDev["org"] = org
	devMap[email] = newDev
	appendState("developers", newId)
	fmt.Printf("Set up developer: %s\n", email)
}

func setupUsers(users []interface{}) {
	for _, user := range users {
		setupUser(user.(map[string]interface{}))
	}
}

func setupUser(user map[string]interface{}) {
	email := user["email"].(string)
	password := user["password"].(string)
	sysKey := setupState["systemKey"].(string)
	sysSec := setupState["systemSecret"].(string)
	adminClient = setupState["adminClient"].(*cb.DevClient)
	newUser, err := adminClient.RegisterNewUser(email, password, sysKey, sysSec)
	if err != nil {
		fatal(err.Error())
	}

	newId := newUser["user_id"].(string)
	addUserToRoles(user, newId)

	usersMap := scriptVars["users"].(map[string]interface{})
	newUser["password"] = password
	usersMap[email] = newUser
	appendState("users", newId)

	// If there are custom fields, must call updateUser
	// for those...
	custFields := getCustomFields(user)
	if len(custFields) == 0 {
		return
	}

	if err = adminClient.UpdateUser(sysKey, newId, custFields); err != nil {
		fatal(err.Error())
	}
}

func getCustomFields(user map[string]interface{}) map[string]interface{} {
	rval := map[string]interface{}{}
	for key, val := range user {
		switch key {
		case "user_id", "email", "password", "creation_date", "roles":
		default:
			rval[key] = val
		}
	}
	return rval
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

func addDeviceToRoles(device map[string]interface{}, deviceName string) {
	if _, ok := device["roles"]; !ok {
		warn(fmt.Sprintf("No roles found for %s\n", device["name"].(string)))
		return
	}
	roleNames := device["roles"].([]interface{})
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
	err := adminClient.AddDeviceToRoles(sysKey, deviceName, roleIds)
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
	//  Create the collection
	colId, err := adminClient.NewCollection(sysKey, col["name"].(string))
	if err != nil {
		fatal(err.Error())
	}
	processEdgeDeployInfo("collection", colId, col)
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
	myPrintf("Set up collection %+v\n", col["name"])
}

func setupConnectCollections(cols []interface{}) {
	for _, col := range cols {
		setupConnectCollection(col.(map[string]interface{}))
	}
}

func setupConnectCollection(col map[string]interface{}) {
	//  Create the collection
	collectionConfig, err := cb.GenerateConnectCollection(col)
	if err != nil {
		fatal(err.Error())
	}
	colId, err := adminClient.NewConnectCollection(sysKey, collectionConfig)
	if err != nil {
		fatal(err.Error())
	}
	colsMap := scriptVars["collections"].(map[string]interface{})
	colsMap[col["name"].(string)] = colId
	appendState("collections", colId)

	setupCollectionRoles(colId, col)
	myPrintf("Set up connect collection %s\n", col["name"].(string))
}

func addThingToRoles(id string, roleNames []interface{}) {
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
	err := adminClient.AddUserToRoles(sysKey, id, roleIds)
	if err != nil {
		fatal(err.Error())
	}
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

func setupPortalRoles(portal, portalRoles map[string]interface{}) {
	roleDict := scriptVars["roles"].(map[string]interface{})
	for roleName, level := range portalRoles {
		roleId, ok := roleDict[roleName]
		if !ok {
			fatal(fmt.Sprintf("Unknown role: %s\n", roleName))
		}
		err := adminClient.AddPortalToRole(sysKey, portal["name"].(string), roleId.(string), int(level.(float64)))
		if err != nil {
			fatal(fmt.Sprintf("Could not add portal to role: %s\n", err.Error()))
		}
	}
}

func setupColumn(collectionId, columnName, columnType string) {
	myPrintf("Set up column %s(%s)\n", columnName, columnType)
	if err := adminClient.AddColumn(collectionId, strings.ToLower(columnName), columnType); err != nil {
		fatal(err.Error())
	}
}

func setupItem(items []map[string]interface{}, colId string) {
	_, err := adminClient.CreateData(colId, items)
	if err != nil {
		fatal(fmt.Sprintf("Error creating item: %s\n", err.Error()))
	}
	myPrintf("Created items for collection %s\n", colId)
}

func setupCodeServices(svcs []interface{}) {
	for _, svc := range svcs {
		setupCodeService(svc.(map[string]interface{}))
	}
}

func mkSvcParams(params []interface{}) []string {
	rval := []string{}
	for _, val := range params {
		rval = append(rval, val.(string))
	}
	return rval
}

func setupCodeService(svc map[string]interface{}) {
	processEdgeDeployInfo("service", svc["name"].(string), svc)
	svcName := getString(svc, "name")
	svcEuid := getString(svc, "euid")
	svcCode := getVarOrFile(svc, "code")
	svcParams := mkSvcParams(svc["parameters"].([]interface{}))
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
	if len(svcEuid) > 0 {
		if err := adminClient.SetServiceEffectiveUser(sysKey, svcName, svcEuid); err != nil {
			fatal(err.Error())
		}
	}
	setupCodeServiceRoles(svc)
	svcMap := scriptVars["codeServices"].(map[string]interface{})
	svcMap[svcName] = map[string]interface{}{
		"name":         svcName,
		"euid":         svcEuid,
		"code":         svcCode,
		"params":       svcParams,
		"dependencies": svcDeps,
	}
	appendState("services", svcName)
	myPrintf("Set up code service %s\n", svcName)
}

func setupCodeLibraries(libs []interface{}) {
	for _, lib := range libs {
		setupCodeLibrary(lib.(map[string]interface{}))
	}
}

func setupCodeLibrary(lib map[string]interface{}) {
	processEdgeDeployInfo("library", lib["name"].(string), lib)
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
	myPrintf("Set up code library %s\n", libName)
}

func setupTriggers(triggers []interface{}) {
	for _, trigger := range triggers {
		setupTrigger(trigger.(map[string]interface{}))
	}
}

func setupTrigger(trigger map[string]interface{}) {
	processEdgeDeployInfo("trigger", trigger["name"].(string), trigger)
	trigName := trigger["name"].(string)
	delete(trigger, "name")
	newTrig, err := adminClient.CreateEventHandler(sysKey, trigName, trigger)
	if err != nil {
		fatal(fmt.Sprintf("Could not create trigger: %s\n", err.Error()))
	}
	trigMap := scriptVars["triggers"].(map[string]interface{})
	trigMap[trigName] = newTrig
	appendState("triggers", trigName)
	myPrintf("Set up trigger %s\n", trigName)
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
	/* swm -- not needed anymore
	if startTime == "Now" {
		startTime = time.Now().Format(time.RFC3339)
	}
	*/
	timer["start_time"] = startTime
	newTimer, err := adminClient.CreateTimer(sysKey, timerName, timer)
	if err != nil {
		fatal(fmt.Sprintf("Could not create timer: %s\n", err.Error()))
	}
	timerMap := scriptVars["timers"].(map[string]interface{})
	timerMap[timerName] = newTimer
	appendState("timers", timerName)
	myPrintf("Set up timer %s\n", timerName)
}

func setupDeviceColumns(deviceColumns []interface{}) {
	for _, deviceColumn := range deviceColumns {
		setupDeviceColumn(deviceColumn.(map[string]interface{}))
	}
}

func setupDeviceColumn(deviceColumn map[string]interface{}) {
	sysKey := setupState["systemKey"].(string)
	adminClient = setupState["adminClient"].(*cb.DevClient)
	if err := adminClient.CreateDeviceColumn(sysKey, deviceColumn["column_name"].(string), deviceColumn["type"].(string)); err != nil {
		fatal(err.Error())
	}
	myPrintf("Added column to %s_edge_devices table: %s\n", sysKey, deviceColumn["column_name"].(string))
}

func setupDevices(devices []interface{}) {
	for _, device := range devices {
		setupDevice(device.(map[string]interface{}))
	}
}

func setupDevice(device map[string]interface{}) {
	deviceName := device["name"].(string)
	roles := device["roles"].([]interface{})
	delete(device, "roles")
	newDevice, err := adminClient.CreateDevice(sysKey, deviceName, device)
	if err != nil {
		fatal(fmt.Sprintf("Could not create device: %s\n", err.Error()))
	}
	newDevice["roles"] = roles
	addDeviceToRoles(newDevice, deviceName)

	deviceMap := scriptVars["devices"].(map[string]interface{})
	deviceMap[deviceName] = newDevice
	appendState("devices", deviceName)
	myPrintf("Set up device %s\n", deviceName)
}

func setupPortals(portals []interface{}) {
	for _, portal := range portals {
		setupPortal(portal.(map[string]interface{}))
	}
}

func setupPortal(portal map[string]interface{}) {
	portalName := portal["name"].(string)
	config := portal["config"].(map[string]interface{})
	configBytes, err := json.Marshal(config)
	if err != nil {
		fatal(fmt.Sprintf("Could not marshal portal config: %s", err))
	}
	var portalRoles map[string]interface{}
	portalRoles, _ = portal["roles"].(map[string]interface{})
	delete(portal, "roles")

	portal["config"] = string(configBytes)
	newPortal, err := adminClient.CreatePortal(sysKey, portalName, portal)
	if err != nil {
		fatal(fmt.Sprintf("Could not create portal: %s\n", err.Error()))
	}
	if portalRoles != nil {
		setupPortalRoles(portal, portalRoles)
	}
	portalMap := scriptVars["portals"].(map[string]interface{})
	portalMap[portalName] = newPortal
	appendState("portals", portalName)
	myPrintf("Set up portal %s\n", portalName)
}

func setupEdges(edges []interface{}) {
	for _, edge := range edges {
		setupEdge(edge.(map[string]interface{}))
	}
}

func setupEdge(edge map[string]interface{}) {
	edgeName := edge["name"].(string)
	edge["system_key"] = sysKey
	edge["system_secret"] = sysSec
	delete(edge, "name")
	newEdge, err := adminClient.CreateEdge(sysKey, edgeName, edge)
	if err != nil {
		fatal(fmt.Sprintf("Could not create edge: %s\n", err.Error()))
	}
	edgeMap := scriptVars["edges"].(map[string]interface{})
	edgeMap[edgeName] = newEdge
	appendState("edges", edgeName)
	myPrintf("Set up edge %s\n", edgeName)
}

func setupMessageHistoryPerms(name string, theGoods map[string]interface{}) {
	roleMap := scriptVars["roles"].(map[string]interface{})
	for roleName, permsIF := range theGoods {
		perms := int(permsIF.(float64))
		roleId, ok := roleMap[roleName].(string)
		if !ok {
			fatalf("Unknown role '%s'", roleName)
		}
		if err := adminClient.AddGenericPermissionToRole(sysKey, roleId, name, perms); err != nil {
			fatalf("Could not add user/device table permissions for %s, %s, %d\n", roleName, name, perms)
		}
		fmt.Printf("Set up message history permissions\n")
	}
}

func setupGenericPermissions(name string, theGoods map[string]interface{}) {
	roleMap := scriptVars["roles"].(map[string]interface{})
	for roleName, permsIF := range theGoods {
		perms := int(permsIF.(float64))
		roleId, ok := roleMap[roleName].(string)
		if !ok {
			fatalf("Unknown role '%s'", roleName)
		}
		if err := adminClient.AddGenericPermissionToRole(sysKey, roleId, name, perms); err != nil {
			fatalf("Could not add user/device table permissions for %s, %s, %d\n", roleName, name, perms)
		}
	}
}

func warn(msg string) {
	myPrintf("Warning: %s\n", msg)
}

func fatal(msg string) {
	myPrintf("Fatal Error: %s\n", msg)
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
		//warn(fmt.Sprintf("The value for %s does not exist\n", theThing))
		return ""
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

/*
func processEdgeSyncInfo(resourceType, resourceName string, resource map[string]interface{}) {
	edgeInfo, ok := resource["syncToEdges"]
	if !ok {
		return
	}
	delete(resource, "syncToEdges")
}
*/

func processEdgeDeployInfo(resourceType, resourceName string, resource map[string]interface{}) {
	edgeInfo, ok := resource["deployToEdges"]
	if !ok {
		return
	}
	delete(resource, "deployToEdges")
	edgesToProcess := gatherAppropriateEdges(edgeInfo)
	edgeSyncStuff := setupState["edgeDeploy"].(map[string][]string)
	key := makeKeyForEdgeStuff(resourceType, resourceName)
	edgeSyncStuff[key] = edgesToProcess

	/*
		for _, edgeName := range edgesToProcess {
			oneEdgeSync := edgeSyncStuff[edgeName]
			resourceSlice := oneEdgeSync[string(resourceType)]
			resourceSlice = append(resourceSlice, resourceName)
			oneEdgeSync[string(resourceType)] = resourceSlice
		}
	*/
}

func gatherAppropriateEdges(edgeInfo interface{}) []string {
	switch edgeInfo.(type) {
	case string:
		edgeStr := edgeInfo.(string)
		if edgeStr == "all" {
			return getAllEdgesNames()
		}
		if edgeStr == "none" {
			return []string{}
		}
		return []string{edgeStr}
	case []interface{}:
		edgeNameSlice := edgeInfo.([]interface{})
		edgesToProcess := make([]string, len(edgeNameSlice))
		for i, edgeIF := range edgeNameSlice {
			edgesToProcess[i] = edgeIF.(string)
		}
		return edgesToProcess
	case []string:
		return edgeInfo.([]string)
	default:
		fatalf("Bad type for key 'edges': %T\n", edgeInfo)
	}
	return []string{} // Not reached -- just makes compiler happy
}

func getAllEdgesNames() []string {
	edgeList, ok := globalSetupInfo["edges"].([]interface{})
	if !ok {
		return []string{}
	}
	rval := make([]string, len(edgeList))
	for i, edgeIF := range edgeList {
		edge := edgeIF.(map[string]interface{})
		rval[i] = edge["name"].(string)
	}
	return rval
}

func makeEdgeDeployStructure() map[string][]string {
	return map[string][]string{}
}

func makeEdgeSyncStructure() map[string]map[string][]string {
	// The key is now defined by makeKeyForEdgeStuff() below. The
	// slice-of-strings value is a list of edge names.
	// For now, scenetest always deploys stuff to the platform.
	theThing := map[string]map[string][]string{}
	allEdges := getAllEdgesNames()
	for _, edge := range allEdges {
		theThing[edge] = map[string][]string{
			cb.ServiceSync: []string{},
			cb.LibrarySync: []string{},
			cb.TriggerSync: []string{},
		}
	}
	return theThing
}

func makeKeyForEdgeStuff(resourceType, resourceName string) string {
	return resourceType + ":::" + resourceName
}

func unpackKeyForEdgeStuff(key string) (string, string) {
	pieces := strings.Split(key, ":::")
	return pieces[0], pieces[1]
}

func makeEdgeListQuery(edges []string) *cb.Query {
	var outerQ *cb.Query

	for _, edgeName := range edges {
		if outerQ == nil {
			outerQ = cb.NewQuery()
			outerQ.EqualTo("name", edgeName)
			continue
		}
		q := cb.NewQuery()
		q.EqualTo("name", edgeName)
		outerQ.Or(q)
	}
	return outerQ
}

func setupEdgeDeployInfo() {
	theInfo := setupState["edgeDeploy"].(map[string][]string)
	fmt.Printf("THE INFO: %+v\n", theInfo)
	for key, edges := range theInfo {
		resourceType, resourceName := unpackKeyForEdgeStuff(key)
		q := makeEdgeListQuery(edges)
		if _, err := adminClient.CreateDeployResourcesForSystem(sysKey, resourceName, resourceType, true, q); err != nil {
			myPrintf("Call to CreateDeployResourcesForSystem failed: %%s\n", err)
		}
		fmt.Printf("ADDED: (%s): %s, %s, %+v\n", sysKey, resourceType, resourceName, edges)
	}
	/*
		for edgeName, edgeStuff := range theInfo {
			fmt.Printf("EDGE NAME: %s\n", edgeName)
			q := makeEdgeQuery(edgeName)
			for resourceType, resourceNames := range edgeStuff {
				for _, resourceName := range resourceNames {
					if _, err := adminClient.CreateDeployResourcesForSystem(sysKey, resourceName, resourceType, true, q); err != nil {
						myPrintf("Call to CreateDeployResourcesForSystem failed: %%s\n", err)
					}
					fmt.Printf("ADDED: (%s): %s, %s, %s\n", sysKey, edgeName, resourceType, resourceName)
				}
			}
		}
	*/
}
