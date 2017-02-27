function HelloClearBladeTriggerHandler(req, resp){
    logStdErr("******************** HelloClearBladeTriggerHandler ********************")
    log(JSON.stringify(req));
    
    ClearBlade.init({request:req});
    var code = ClearBlade.Code();
    
    //Invoke another service
    code.execute("HelloClearBlade", null, true, function(err, data){
            if(err) {
                logStdErr("******************** Call Error " + data + "********************")
                log("Error encountered executing HelloClearBlade service");
            } else {
                logStdErr("******************** Call Success " + data + "********************")
                log("HelloClearBlade service executed successfully");
            }
            log(JSON.stringify(data));
    })
}
