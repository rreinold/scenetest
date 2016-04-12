function justPublishOneMessage(req, resp, topic, payload) {
    ClearBlade.init({request:req});
    var messaging = ClearBlade.Messaging({}, function(){});
    result = messaging.publish(topic, payload);
}
