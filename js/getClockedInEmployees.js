function getClockedInEmployees(req, resp) {
    ClearBlade.init({"request":req});
    var stuff = doGetClockedInEmployees();
    if (stuff === null) {
        logStdErr("Call to getEmployees failed");
        resp.failure("Call to getEmployees failed");
    } else {
        resp.success(stuff);
    }
}
