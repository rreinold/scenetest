function stopTimer(timerCollection, id, complete) {
    var nowDate = new Date(Date.now());
    var tmpStr = nowDate.toString();
    var theDateStr = tmpStr.substring(0, tmpStr.length - 10); //  Hack to get rid of milliseconds and the offset from GMT
    var changes = { "status": 1, "end": theDateStr};
    var updateQ = ClearBlade.Query();
    updateQ.equalTo("ID", id);
    timerCollection.update(updateQ, changes, function (error, response) {
        if (error === true) {
            logStdErr("UPDATE TIMER ENTRY FAILED: " + JSON.stringify(response));
        } else {
            logStdErr("UPDATE TIMER ENTRY SUCCEEDED: " + JSON.stringify(response));
        }
    });
    notifyUpdateClockInPortal();
}

function beginStartTimerProcess(timerCollection, edgeName, tagId, projectId, taskId) {
    var newTimerEntry = {
        "task": taskId,
        "project": projectId,
        "status": 3,
        "edge_name": edgeName + ":" + tagId,
        "notes": "RFID"
    };

    timerCollection.create(newTimerEntry, function(error, response) {
        if (error === true) {
            logStdErr("CREATE NEW ROW IN TIMER TABLE FAILED: " + JSON.stringify(response));
        } else {
            logStdErr("SUCCESSFULLY CREATED NEW ROW TIMER TABLE: " + JSON.stringify(response));
        }
    });
    notifyUpdateClockInPortal();
}

function notifyUpdateClockInPortal() {
    var stuff =  doGetClockedInEmployees();
    var tryClockIn = doGetTryingToClockIn();
    var messaging = ClearBlade.Messaging({}, function(){});
    messaging.publish("/w2macfab/clock/clockinout", JSON.stringify(stuff));
    messaging.publish("/w2macfab/clock/tryclockin", JSON.stringify(tryClockIn));
}

function doGetClockedInEmployees() {
    logStdErr("B");
    var timerC = ClearBlade.Collection({"collectionName": "MacFab Timer"});
    var timerQ = ClearBlade.Query();
    timerQ.equalTo("status", 0);
    timerQ.notEqualTo("project", 0); // ignore bogus entries
    timerQ.notEqualTo("task", 0); // ignore bogus entries
    var clockedIn;
    timerC.fetch(timerQ, function(error, timerResp) {
        if (error === true) {
            logStdErr("Timer query failed: " + JSON.stringify(timerResp));
            return null;
        }
        clockedIn = timerResp.DATA;
    });
    var cLen = clockedIn.length;
    var result = new Array(cLen);
    logStdErr("C");

    var projCol = ClearBlade.Collection({"collectionName": "MacFab Project"});
    var taskCol = ClearBlade.Collection({"collectionName": "MacFab Task"});
    var projQ = null;
    var taskQ = null;
    if (cLen === 0) {
        return [];
    }
    for (var i = 0; i < cLen; i ++) {
        var curRec = clockedIn[i];

        //  Chain up or queries for project
        var subPQ = ClearBlade.Query();
        subPQ.equalTo("project_id", curRec.project);
        if (projQ === null) {
            projQ = subPQ;
        } else {
            projQ.or(subPQ);
        }
        logStdErr("");
        logStdErr("PROJQ: " + JSON.stringify(projQ));
        logStdErr("");

        //  Chain up or queries for task
        var subTQ = ClearBlade.Query();
        subTQ.equalTo("task_id", curRec.task);
        if (taskQ === null) {
            taskQ = subTQ;
        } else {
            taskQ.or(subTQ);
        }

    }

    //  Now, fetch the tasks and projects related to the timer entries. What this all
    //  boils down to is that we're faking out joins.
    var projects, tasks;

    projCol.fetch(projQ, function(error, r) {
        if (error === true) {
            logStdErr("project query failed: " + JSON.stringify(r));
            return null;
        }
        logStdErr("RESPONSE: " + JSON.stringify(r));
       
        projects = r.DATA;
    });

    taskCol.fetch(taskQ, function(error, r) {
        if (error === true) {
            logStdErr("task query failed: " + JSON.stringify(r));
            return null;
        }
        tasks = r.DATA;
    });

    // Now, loop again, creating the "joined" rows/records
    for (i = 0; i < cLen; i ++) {
        var timer = clockedIn[i];
        theTask = findObjInList(tasks, "task_id", timer.task);
        theProj = findObjInList(projects, "project_id", timer.project);
        //logStdErr("PROJECTS: " + JSON.stringify(projects));
        if (theTask === null) {
            logStdErr("AAAAAAHHHHH -- couldn't find task for " + timer.task);
            return null;
        }
        if (theProj === null) {
            logStdErr("AAAAAAHHHHH -- couldn't find project for " + timer.project);
            return null;
        }

        result[i] = {
            "Name": theTask.name,
            "Project": theProj.name,
            "ProjectID": timer.project,
            "ClockedIn": timer.start,
            "Notes": timer.notes
        };
    }

    return result;
}

function findObjInList(list, key, value) {
    for (var i = 0; i < list.length; i ++) {
        if (list[i][key] == value) {
            return list[i];
        }
    }
    return null;
}

function doGetTryingToClockIn() {
    logStdErr("DO GET TRYING TO CLOCK IN");
    var timerCol = ClearBlade.Collection({"collectionName": "MacFab Timer"});
    var taskCol = ClearBlade.Collection({"collectionName": "MacFab Task"});
    var timerQ = ClearBlade.Query();
    timerQ.equalTo("status", 3);
    var posers = null;
    var results = [];
    logStdErr("BOFFO!");
    timerCol.fetch(timerQ, function(error, r) {
        if (error === true) {
            log("fetch of timer table failed");
            resp.failure("Fetch of timer table failed");
        }
        posers = r.DATA;
    });
    for (var i = 0; i < posers.length; i++) {
        logStdErr("LOOP");
        var curPoser = posers[i];
        var edgeAndTag = curPoser.edge_name;
        if (edgeAndTag === "") {
            log("Skipping log in attempt with no edge/tag: " + JSON.stringify(curPoser));
            continue;
        }
        var splitEdgeAndTag = edgeAndTag.split(":");
        if (splitEdgeAndTag.length !== 2) {
            log("edge_name field has unknown format: " + edgeAndTag);
            continue;
        }
        var edgeName = splitEdgeAndTag[0];
        var tagID = splitEdgeAndTag[1];
        var taskQ = ClearBlade.Query();
        taskQ.equalTo("tagID", tagID);
        taskCol.fetch(taskQ, function(error, r) {
            if (error === true) {
                log("task collection fetch failed");
                resp.failure("task collection fetch failed.");
            }
            if (r.DATA.length === 0) {
                log("could not find tag: " + tagID + " in task table");
                return;
            }
            var t = r.DATA[0];
            var newEntry = {
                "Name": t.name,
                "ClockInStarted": curPoser.created,
                "TagID": tagID,
                "Edge": edgeName,
                "Task": curPoser.task
            };
            results.push(newEntry);
        });
    }
    return results;
}

function employeeIsClockingIn(tagID) {
    var timerCol = ClearBlade.Collection({"collectionName": "MacFab Timer"});
    var tq = ClearBlade.Query();

    posers = doGetTryingToClockIn();
    if (posers === null) {
        resp.failure("Couldn't get folks trying to clock in");
    }
    
    for (var i = 0; i < posers.length; i ++) {
        var poser = posers[i];
        if (tagID == poser.TagID) {
            return true;
        }
    }
    return false;
}

function doFinishClockIn(taskId, projectId) {
    var q = ClearBlade.Query();
    var timerCol = ClearBlade.Collection({"collectionName": "MacFab Timer"});
    q.equalTo("task", taskId);
    q.equalTo("status", 3);
    var nowDate = new Date(Date.now());
    var tmpStr = nowDate.toString();
    var theDateStr = tmpStr.substring(0, tmpStr.length - 10);
    logStdErr("DATE IS " + theDateStr);
    var updates = {
        "status": 0,
        "project": projectId,
        "start": theDateStr
    };
    var rval = true;
    timerCol.update(q, updates, function(error, r) {
    //timerCol.fetch(q, function(error, r) {
        if (error === true) {
            log("timer fetch failed");
            logStdErr("timer fetch failed");
            rval = false;
            return
        }
        /*
        if (r.DATA.length === 0) {
            log("couldn't find anybody to log in");
            logStderr("Could not find anybody to log in");
            rval = false;
            return;
        }
        */
        logStdErr("SHOULD BE LOGGING IN");
    });
    
    notifyUpdateClockInPortal();
    return rval;
}
