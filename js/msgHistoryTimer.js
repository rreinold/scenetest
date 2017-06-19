function msgHistoryTimer(req, resp) {
    logStdErr("TIMER ON THE EDGE");
    ClearBlade.init({request:req});
    var messaging = ClearBlade.Messaging({}, function(){});
    if (ClearBlade.isEdge() === true) {
        messaging.getAndDeleteMessageHistory("device/updateForTimer", 0, null, null, null, function(failed, results) {
            if (failed) {
                logStdErr("getAndDeleteMessageHistory FAILED: " + JSON.stringify(results));
            } else {
                logStdErr("getAndDeleteMessageHistory SUCCEEDED: " + JSON.stringify(results));
                if (results.length > 0) {
                    logStdErr("GOT SOME\n");
                    highest = highestStateValue(results);
                    ClearBlade.updateDevice(highest.deviceName, {state:highest.state.toString()}, true);
                    messaging.publish("scenetest/deviceMessageCount", results.length.toString())
                }
            }
        });
        resp.success(ClearBlade.edgeId() + " Everybody loves somebody sometime");
    }
}

function highestStateValue(results) {
    log("a");
    highest = {"name": "Does not exist", "state": 0};
    for (i = 0; i < results.length; i ++) {
        log("b");
        payload = JSON.parse(results[i].payload);
        log(payload.state.toString());
        if (payload.state > highest.state) {
            logStdErr("YOWZA")
            highest = payload;
        }
    }
    log("c");
    logStdErr("HIGHEST: " + highest.state.toString());
    return highest;
}
