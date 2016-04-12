function publishTrigger(req, resp, msgClass, msgType) {
    var msgInfo = {
        msgClass:msgClass,
        msgType:msgType
    };
    ClearBlade.init({request:req});
    var messaging = ClearBlade.Messaging({}, function(){});
    messaging.publish("/clearblade/internal/trigger", JSON.stringify(msgInfo));
}
