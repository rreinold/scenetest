function publishTrigger(resp, sysKey, sysSec, userToken, msgClass, msgType) {
    var msgInfo = {
        msgClass:msgClass,
        msgType:msgType
    };

    var publishMessage = function() {
        log("DOH");
        var messaging = ClearBlade.Messaging({}, function(){});
        messaging.publish("/clearblade/internal/trigger", JSON.stringify(msgInfo));
    };

    var fiddlestix = {
        systemKey: sysKey,
        systemSecret: sysSec,
        userToken: userToken,
        callback: publishMessage
    };

    try {
        log("DAH");
        ClearBlade.init(fiddlestix);
        log("DEEG");
    }
    catch(err) {
        log("Caught exception: " + err.message)
    }
}


