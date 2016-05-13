function realUpdateDeviceTrigger(req, resp) {
    ClearBlade.init({request:req});
    if (ClearBlade.isEdge() === true) {
        logStdErr("update device trigger on EDGE");
        publishTrigger(req, resp, "Device", "DeviceUpdated");
        resp.success(ClearBlade.edgeId() + " Got the update device trigger pull");
    } else {
        logStdErr("update device trigger on NOVI");
        publishTrigger(req, resp, "Device", "DeviceUpdated");
        resp.success("Novi Ignoring update device trigger pull");
    }
}
