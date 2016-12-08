function processBarcodeScan(req, resp){
    ClearBlade.init({"request":req});
    logStdErr("processBarcodeScan: " + JSON.stringify(req));
}