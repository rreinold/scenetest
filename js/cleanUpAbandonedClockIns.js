function cleanUpAbandonedClockIns(req, resp){
    ClearBlade.init({"request": req});
    if (ClearBlade.isEdge() === true) {
        resp.success("I am an edge. I'm not going to do anything");
        return;
    }

    var timerCol = ClearBlade.Collection({"collectionName": "MacFab Timer"});
    var tq = ClearBlade.Query();
    tq.equalTo("status", 3);
    var inProgress = null;
    timerCol.fetch(tq, function(error, r) {
        if (error === true) {
            logStdErr("clearUpAbandonedClockIns: timer fetch failed");
            resp.failure("Error in fetch of timer collection");
            return;
        }
        inProgress = r.DATA;
    });
    for (var i = 0; i < inProgress.length; i ++) {
        var cur = inProgress[i];
        logStdErr("");
        logStdErr("CLOCK IN CLEANUP: " + JSON.stringify(cur));
        logStdErr("");
        
        var createdStr = cur.created;
        var d = splitDateString(createdStr);
        var month = fixMonth(d[1]);
        var createdDate =  new Date(d[0], month, d[2], d[3], d[4], d[5], 0);
        var nowDate = new Date(Date.now());
        var diff = nowDate - createdDate;
        var minutes = diff / (1000*60);
        logStdErr("Date DIFF IS " + JSON.stringify(minutes));
        logStdErr("Created: " + createdDate.toString() + ", Now: " + nowDate.toString());
        if (diff >= 1) { // timeout can change to whatever you want
            var uq = ClearBlade.Query();
            uq.equalTo("ID", cur.ID);
            var updates = {"status": 1};
            timerCol.update(uq, updates, function(error, r) {
                if (error === true) {
                    logStdErr("Strange update error ignored");
                    return;
                }
            });
            logStdErr("We just clocked out this employee: " + JSON.stringify(cur));
        }
    }
    notifyUpdateClockInPortal();
    resp.success("CLEAN UP PROCESSED: " + JSON.stringify(inProgress)); 
}

function fixMonth(m) {
    var monthNum = Number(m) - 1;
    return monthNum.toString();
    
}

function splitDateString(ds) {
	var dt = ds.split(" ");
	var ymd = dt[0].split("-");
	var hms = dt[1].split(":");
	var rval = new Array(6);
	rval[0] = ymd[0];
	rval[1] = ymd[1];
	rval[2] = ymd[2];
	rval[3] = hms[0];
	rval[4] = hms[1];
	rval[5] = hms[2];
	return rval;
}
