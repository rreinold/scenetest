function fetchSensorHistory(req, resp){
    var startTime = Date.now();
    ClearBlade.init({"request": req});
    var query = ClearBlade.Query({collectionName: "BACNET Data"});
    query.equalTo("device_name", "chain saw");
    query.greaterThan("time", "2017-07-29T13:41:38Z");
    query.columns(["time", "present_value"]);
    query.descending("time");
    query.setPage(1000000, 1);
    logStdErr("YOW!");
    query.fetch(function(err, data) {
        if (err) {
            log.error("fetch error : " + JSON.stringify(data));
        } else {
            logStdErr("here we go!");
            //log(data); // hahahahahaha
            ret = {sensorName: "chain saw", data: [] };
            for (var j = 0; j < data.DATA.length; j++) {
                ret.data.push({"x": data.DATA[j].time, "y": data.DATA[j].present_value});
            }
            ret.executionTime = JSON.stringify(((Date.now() - startTime) / 1000));
            logStdErr("recordcount: " + ret.data.length);
            resp.success(ret);
           // resp.success("It worked!");
        }
    });
}
