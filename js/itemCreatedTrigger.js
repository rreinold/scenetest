
function itemCreatedTrigger(req, resp) {
    sysKey = req.systemKey;
    var sysSec = req.systemSecret;
    var userToken = req.userToken;
    publishTrigger(req, resp, "Data", "ItemCreated");
    resp.success("We win! Request object: " + JSON.stringify(req));
}
