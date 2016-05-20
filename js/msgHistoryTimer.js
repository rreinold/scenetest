function msgHistoryTimer(req, resp) {
    logStdErr("TIMER ON THE EDGE");
    ClearBlade.init({request:req});
    var messaging = ClearBlade.Messaging({}, function(){});
    if (ClearBlade.isEdge() === true) {
        messaging.getAndDeleteMessageHistory("device/updateForTimer", 100, function(failed, results) {
            if (failed) {
                logStdErr("getAndDeleteMessageHistory FAILED: " + JSON.stringify(results))
            } else {
                logStdErr("getAndDeleteMessageHistory SUCCEEDED: " + JSON.stringify(results))
                if (results.length > 0) {
                    resObj = JSON.parse(results[0].payload)
                    ClearBlade.updateDevice(resObj.deviceName, {state:resObj.state.toString()}, false);
                }
            }
        });
        resp.success(ClearBlade.edgeId() + " Everybody loves somebody sometime");
    }
}
