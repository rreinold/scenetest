
function userDeletedTrigger(req, resp) {
    publishTrigger(req, resp, "User", "UserDeleted");
    resp.success("Should have published by now");
}
