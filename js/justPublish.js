function justPublishOneMessage(req, resp, topic, payload) {
    var theTruth = {
        email: req.userEmail,
        systemKey: req.systemKey,
        systemSecret: req.systemSecret,
        userToken: req.userToken
    }
    ClearBlade.init({request:req});
    //ClearBlade.init(theTruth);
    var messaging = ClearBlade.Messaging({}, function(){});
    result = messaging.publish(topic, payload);
}
