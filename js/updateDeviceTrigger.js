function updateDeviceTrigger(req, resp) {
    logStdErr("In update device TRIGGER")
    ClearBlade.init({request:req});
    var body = JSON.parse(req.params.body);
    body.state = body.state.toString();


    ClearBlade.updateDevice("Dorky Torque Wrench", body, true);

    resp.success("I think it worked");
}
