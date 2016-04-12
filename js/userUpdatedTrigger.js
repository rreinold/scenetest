
function userUpdatedTrigger(req, resp) {
    sysKey = req.systemKey;
    var sysSec = req.systemSecret;
    var userToken = req.userToken;
    publishTrigger(req, resp, "User", "UserUpdated");
    resp.success("We win! Request object: " + JSON.stringify(req));
}
