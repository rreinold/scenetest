
function collectionDeletedTrigger(req, resp) {
    publishTrigger(req, resp, "Data", "CollectionDeleted");
    resp.success("Collection Deleted Trigger Fired");
}
