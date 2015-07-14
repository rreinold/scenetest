
function itemCreatedTrigger(req, resp) {
    sysKey = req.systemKey;
    var sysSec = req.systemSecret;
    var userToken = req.userToken;
    log("BEGIN");
    publishTrigger(resp, sysKey, sysSec, userToken, "Data", "ItemCreated");
    log("ENDDDDD");
    resp.success("Should have published by now");
}
