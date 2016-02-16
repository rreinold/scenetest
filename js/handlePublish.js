
function handlePublish(req, resp) {
    sysKey = req.systemKey;
    var sysSec = req.systemSecret;
    var userToken = req.userToken;
    justPublishOneMessage(req, resp, "/timer/popped", "Not much to say");
    publishTrigger(req, resp, "Messaging", "Publish");
    resp.success("Should have published by now");
}
