function JustUpdateDevice(req, resp) {
    ClearBlade.init({request:req});
    logStdErr("Just Update Device: " + JSON.stringify(req.params))


    ClearBlade.updateDevice("Dorky Torque Wrench", req.params.updates, true);

    resp.success("I think it worked");
}
