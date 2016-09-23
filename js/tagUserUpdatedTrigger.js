
function tagUserUpdatedTrigger(req, resp) {
    logStdErr("tagUserUpdatedTrigger: " + JSON.stringify(req));
    var stuff = req.params.user.state;
    logStdErr("STUFF: " + stuff);
    var parts = stuff.split(":");
    logStdErr("Edge: " + parts[0] + ", Tag ID: " + parts[1]);
        

    ClearBlade.init({"request":req});
    var col = ClearBlade.Collection({"collectionName": "timer"});
    logStdErr("COLLECTION IS: " + JSON.stringify(col));


    resp.success("Done");
}
