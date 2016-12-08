function tagUserUpdatedTrigger(req, resp) {
    logStdErr("tagUserUpdatedTrigger: " + JSON.stringify(req));
    var stuff = req.params.user.state;
    logStdErr("STUFF: " + stuff);
    var parts = stuff.split(":");
    logStdErr("Edge: " + parts[0] + ", Tag ID: " + parts[1]);
      
    var edgeID = parts[0];
    var tagID = parts[1];

    ClearBlade.init({"request":req});
    var timerCol = ClearBlade.Collection({"collectionName": "MacFab Timer"});
    var taskCol = ClearBlade.Collection({"collectionName": "MacFab Task"});
    var taskQ = ClearBlade.Query();
    taskQ.equalTo("tagID", tagID);
    var taskId = "unknown";
    taskCol.fetch(taskQ, function(error, response) {
        if (error === true) {
            logStdErr("task fetch failed: " + JSON.stringify(response));
            resp.failure("Task query failure");
        }
        if (response.DATA.length === 0) {
            var errStr = "Could not find tag " + tagId;
            logStdErr(errStr);
            resp.failure(errStr);
        }
        taskId = response.DATA[0].task_id;
    });
    
    //  If the employee taps his tag more than once, throw away the extra taps
    if (employeeIsClockingIn(tagID) === true) {
        logStdErr("WOWOWOWO -- USER IS CLOCKING IN ALREAYD");
        resp.success("Employee is already trying to clock in");
    }

    var timerQ = ClearBlade.Query();
    logStdErr("TASK ID IS " + taskId);
    timerQ.equalTo("task", taskId);
	timerQ.equalTo("status", 0);
    timerCol.fetch(timerQ, function(error, response) {
        if (error === true) {
            logStdErr("timer fetch failed: " + JSON.stringify(response));
            resp.failure("Timer query failure");
        }
        logStdErr("BOTH QUERIES WORKED: " + JSON.stringify(response));
        var stuff = response.DATA;
        if (stuff.length === 0) {
            logStdErr("NEED TO CREATE NEW TIMER");
            beginStartTimerProcess(timerCol, edgeID, tagID, 102082, taskId);
        } else {
            logStdErr("NEED TO END CURRENT TIMER: " + stuff[0].ID);
            stopTimer(timerCol, stuff[0].ID, 0);
        }
    });

    resp.success("Done");
}
