function finishClockIn(req, resp){
    log(JSON.stringify(req));
    ClearBlade.init({"request": req});
    var taskId = req.params.taskId;
    var projectId = req.params.projectId;
    taskIdStr = JSON.stringify(taskId);
    projectIdStr = JSON.stringify(projectId);
    logStdErr("taskId and projectId are " + taskIdStr + ", " + projectIdStr);

    if (taskIdStr == undefined || taskIdStr == null || taskIdStr == "") {
        resp.error("Remember to set the taskId parameter");
    }
    if (projectIdStr == undefined || projectIdStr == null || projectIdStr == "") {
        resp.error("Remember to set the projectId parameter");
    }
    if (taskIdStr === projectIdStr) {
        logStdErr("SAME");
        resp.error("task and project are the same");
    }
    
    if (doFinishClockIn(taskId, projectId) === false) {
        resp.error("doFinishClockIn did not work");
    }
    resp.success("The Employee should be clocked in");
}