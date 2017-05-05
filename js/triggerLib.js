function publishTrigger(req, resp, msgClass, msgType) {
    var msgInfo = {
        msgClass:msgClass,
        msgType:msgType
    };
    ClearBlade.init({request:req});
    var messaging = ClearBlade.Messaging({}, function(){});
    logStdErr("PUBLISHING TO internal trigger topic: " + JSON.stringify(msgInfo));
    messaging.publish("/clearblade/internal/trigger", JSON.stringify(msgInfo));
    logStdErr("FINISHED WITH THE PUBLISH CALL")
}

function publishTriggerWithBody(req, resp, msgClass, msgType, msgBody) {
    var msgInfo = {
        msgClass:msgClass,
        msgType:msgType,
        msgBody:msgBody
    };
    ClearBlade.init({request:req});
    var messaging = ClearBlade.Messaging({}, function(){});
    logStdErr("PUBLISHING TO internal trigger topic: " + JSON.stringify(msgInfo));
    messaging.publish("/clearblade/internal/trigger", JSON.stringify(msgInfo));
    logStdErr("FINISHED WITH THE PUBLISH CALL")
}
