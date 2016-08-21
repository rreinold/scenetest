var getUserCallback = function(err, getResults) {
    if (err === true) {
        logStdErr("Query Failed");
            resp.failure(getResults);
    }
    logStdErr("getUser results: " + JSON.stringify(getResults));
};

function tagFoundTrigger(req, resp) {
    logStdErr("In Tag Found Trigger: " + req.params.body )
    ClearBlade.init({request:req});

    if (ClearBlade.isEdge() == false) {
        resp.success("Ignoring, I'm not an edge");
    }
    var edgeName = ClearBlade.edgeId();
    logStdErr("My edge name is: " + edgeName);
    var body = JSON.parse(req.params.body);

    var q = ClearBlade.Query();
    q.equalTo("tag_id", body["tagId"]);

    ClearBlade.User().allUsers(q, getUserCallback);
    resp.success("Hmmm");
}
