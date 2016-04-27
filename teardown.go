package main

import (
	"fmt"
	cb "github.com/clearblade/Go-SDK"
)

var adm *cb.DevClient

func performTeardown() {
	teardownSystem(scriptVars["teardown"].(map[string]interface{}))
}

func authTheDevGod(system map[string]interface{}) {
	adm = cb.NewDevClient(system["dev_email"].(string), system["dev_password"].(string))
	if err := adm.Authenticate(); err != nil {
		fatal(fmt.Sprintf("Could not auth the dev god: %s\n", err.Error()))
	}
}

func setUrlValues(system map[string]interface{}) {
}

func teardownSystem(system map[string]interface{}) {
	setUrlValues(system)
	sysKey = system["systemKey"].(string)
	authTheDevGod(system)

	deleteEdges(system)
	deleteDevices(system)
	deleteTimers(system)
	deleteTriggers(system)
	deleteLibraries(system)
	deleteServices(system)
	deleteCollections(system)
	deleteUsers(system)
	deleteRoles(system)
	deleteSystem(system)
	deleteDeveloper(system)

}

func deleteEdges(system map[string]interface{}) {
	for _, edgeName := range system["edges"].([]interface{}) {
		if err := adm.DeleteEdge(sysKey, edgeName.(string)); err != nil {
			warn(fmt.Sprintf("Could not delete edge %v: %s", edgeName, err.Error()))
		} else {
			myPrintf("Deleted edge %v\n", edgeName)
		}
	}
}

func deleteDevices(system map[string]interface{}) {
	for _, deviceName := range system["devices"].([]interface{}) {
		if err := adm.DeleteDevice(sysKey, deviceName.(string)); err != nil {
			warn(fmt.Sprintf("Could not delete device %v: %s", deviceName, err.Error()))
		} else {
			myPrintf("Deleted device %v\n", deviceName)
		}
	}
}

func deleteTimers(system map[string]interface{}) {
	for _, timerName := range system["timers"].([]interface{}) {
		if err := adm.DeleteTimer(sysKey, timerName.(string)); err != nil {
			warn(fmt.Sprintf("Could not delete timer %v: %s -- it may have already expired", timerName, err.Error()))
		} else {
			myPrintf("Deleted timer %v\n", timerName)
		}
	}
}

func deleteTriggers(system map[string]interface{}) {
	for _, triggerName := range system["triggers"].([]interface{}) {
		if err := adm.DeleteEventHandler(sysKey, triggerName.(string)); err != nil {
			warn(fmt.Sprintf("Could not delete trigger %v: %s", triggerName, err.Error()))
		} else {
			myPrintf("Deleted trigger %v\n", triggerName)
		}
	}
}

func deleteLibraries(system map[string]interface{}) {
	for _, libraryName := range system["libraries"].([]interface{}) {
		if err := adm.DeleteLibrary(sysKey, libraryName.(string)); err != nil {
			warn(fmt.Sprintf("Could not delete library %v: %s", libraryName, err.Error()))
		} else {
			myPrintf("Deleted library %v\n", libraryName)
		}
	}
}

func deleteServices(system map[string]interface{}) {
	for _, serviceName := range system["services"].([]interface{}) {
		if err := adm.DeleteService(sysKey, serviceName.(string)); err != nil {
			warn(fmt.Sprintf("Could not delete service %v: %s", serviceName, err.Error()))
		} else {
			myPrintf("Deleted service %v\n", serviceName)
		}
	}
}

func deleteCollections(system map[string]interface{}) {
	for _, colId := range system["collections"].([]interface{}) {
		if err := adm.DeleteCollection(colId.(string)); err != nil {
			warn(fmt.Sprintf("Could not delete collection %v: %s", colId, err.Error()))
		} else {
			myPrintf("Deleted collection %v\n", colId)
		}
	}
}

func deleteUsers(system map[string]interface{}) {
	for _, userId := range system["users"].([]interface{}) {
		if err := adm.DeleteUser(sysKey, userId.(string)); err != nil {
			warn(fmt.Sprintf("Could not delete user %v: %s", userId, err.Error()))
		} else {
			myPrintf("Deleted user %v\n", userId)
		}
	}
}

func deleteRoles(system map[string]interface{}) {
	for _, roleId := range system["roles"].([]interface{}) {
		if err := adm.DeleteRole(sysKey, roleId.(string)); err != nil {
			warn(fmt.Sprintf("Could not delete role %v: %s", roleId, err.Error()))
		} else {
			myPrintf("Deleted role %v\n", roleId)
		}
	}
}

func deleteSystem(system map[string]interface{}) {
	if err := adm.DeleteSystem(sysKey); err != nil {
		warn(fmt.Sprintf("Could not delete system: %s\n", err.Error()))
	} else {
		myPrintf("Deleted system %v\n", sysKey)
	}
}

func deleteDeveloper(system map[string]interface{}) {
	myPrintf("Developer %s is trying to commit suicide, but the platform won't let him\n", system["dev_email"].(string))
}
