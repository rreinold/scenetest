function msgHistoryTimer(req, resp) {
    logStdErr("TIMER ON THE EDGE");
    ClearBlade.init({request:req});
    var messaging = ClearBlade.Messaging({}, function(){});
    if (ClearBlade.isEdge() === true) {
        messaging.getAndDeleteMessageHistory("/device/updateForTimer", 100, function(failed, results) {
            if (failed) {
                logStdErr("getAndDeleteMessageHistory FAILED: " + JSON.stringify(results))
            } else {
                logStdErr("getAndDeleteMessageHistory SUCCEEDED: " + JSON.stringify(results))
                // Just for now -- will grab info when we see the format
                ClearBlade.updateDevice("Dorky Torque Wrench", {state:"Spiffy"})
            }
        });
        //publishTrigger(req, resp, "Device", "DeviceUpdated");
        resp.success(ClearBlade.edgeId() + " Everybody loves somebody sometime");
    }
}
