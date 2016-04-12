function updateDeviceTrigger(req, resp) {
    ClearBlade.init({request:req});
    var body = JSON.parse(req.params.body);
    body.state = body.state.toString();


    ClearBlade.updateDevice("Dorky Torque Wrench", body)

    resp.success("I think it worked");
}
