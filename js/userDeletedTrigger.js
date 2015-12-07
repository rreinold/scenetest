
function userDeletedTrigger(req, resp) {
    sysKey = req.systemKey;
    var sysSec = req.systemSecret;
    var userToken = req.userToken;
    publishTrigger(req, resp, "User", "UserDeleted");
    resp.success("Should have published by now");
}
