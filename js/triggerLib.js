function publishTrigger(resp, sysKey, sysSec, userToken, msgClass, msgType) {
    var msgInfo = {
        msgClass:msgClass,
        msgType:msgType
    };

    var publishMessage = function() {
        var messaging = ClearBlade.Messaging({}, function(){});
        messaging.publish("/clearblade/internal/trigger", JSON.stringify(msgInfo));
        log("PUBLISHED")
    };

    var fiddlestix = {
        systemKey: sysKey,
        systemSecret: sysSec,
        userToken: userToken,
        callback: publishMessage
    };

    ClearBlade.init(fiddlestix);
}
