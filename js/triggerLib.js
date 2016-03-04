function publishTrigger(req, resp, msgClass, msgType) {
    var msgInfo = {
        msgClass:msgClass,
        msgType:msgType
    };
    var yowza = {
        steve: "This is steven",
        email: req.userEmail,
        systemKey: req.systemKey,
        systemSecret: req.systemSecret,
        userToken: req.userToken
    }
    ClearBlade.init({request:req});
    //ClearBlade.init(yowza);
    var messaging = ClearBlade.Messaging({}, function(){});
    messaging.publish("/clearblade/internal/trigger", JSON.stringify(msgInfo));
}
