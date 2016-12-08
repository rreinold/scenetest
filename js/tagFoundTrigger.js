var gTagId = "Tag id not set";

var getUserCallback = function(err, getResults) {
    if (err === true) {
        logStdErr("Query Failed");
        resp.failure(getResults);
    }
    logStdErr("getUser results: " + JSON.stringify(getResults));
    // Need to modify the user object
    var user = getResults["Data"][0];
    var email = user["email"];
    var state = ClearBlade.edgeId() + ":" + gTagId;
    var uq = ClearBlade.Query();
    uq.equalTo("email", email);
    ClearBlade.User().setUsers(uq, {"state": state});
    logStdErr("DONE SETTING USERS: " +  state);
};

function tagFoundTrigger(req, resp) {
    logStdErr("In Tag Found Trigger: " + req.params.body );
    ClearBlade.init({request:req});

    if (ClearBlade.isEdge() === false) {
        resp.success("Ignoring, I'm not an edge");
    }
    var edgeName = ClearBlade.edgeId();
    logStdErr("My edge name is: " + edgeName);
    var body = JSON.parse(req.params.body);

    var q = ClearBlade.Query();
    q.equalTo("tag_id", body.tagId);
    gTagId = body.tagId;

    ClearBlade.User().allUsers(q, getUserCallback);
    logStdErr("GONNA DO RESP.SUCCESS NOW");
    resp.success("Hmmm");
}
