function barcodeScanned(req, resp){
    ClearBlade.init({"request":req});
    if (ClearBlade.isEdge() === false) {
        logStdErr("barcodeScanned: Not an edge");
        resp.success("Not an edge");
    }
    var myEdge = ClearBlade.edgeId();
    logStdErr("In barcode scanned: " + JSON.stringify(req));
    var payload = JSON.parse(req.params.body);
    var theBarcode = payload.barcode;
    logStdErr("BARCODE IS " + theBarcode);
    var deviceName = req.userEmail;
    logStdErr("DEVICE NAME IS " + deviceName);
    ClearBlade.updateDevice(deviceName, {state:myEdge + "::" + theBarcode}, true);
}
