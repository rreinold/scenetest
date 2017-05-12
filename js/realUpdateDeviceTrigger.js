function realUpdateDeviceTrigger(req, resp) {
    ClearBlade.init({request:req});
    if (ClearBlade.isEdge() === true) {
        logStdErr("update device trigger on EDGE" + JSON.stringify(req));
        publishTriggerWithBody(req, resp, "Device", "DeviceUpdated", JSON.stringify(req));
        //resp.success(ClearBlade.edgeId() + " Got the update device trigger pull");
    } else {
        logStdErr("update device trigger on NOVI" + JSON.stringify(req));
        publishTriggerWithBody(req, resp, "Device", "DeviceUpdated", JSON.stringify(req));
        //resp.success("Novi Ignoring update device trigger pull");
    }
}
