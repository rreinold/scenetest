
function itemCreatedTrigger(req, resp) {
    sysKey = req.systemKey;
    var sysSec = req.systemSecret;
    var userToken = req.userToken;
    publishTrigger(resp, sysKey, sysSec, userToken, "Data", "ItemCreated");
    resp.success("Should have published by now");
}
