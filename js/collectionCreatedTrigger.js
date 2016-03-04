
function collectionCreatedTrigger(req, resp) {
    publishTrigger(req, resp, "Data", "CollectionCreated");
    resp.success("Collection Created Trigger Fired");
}
