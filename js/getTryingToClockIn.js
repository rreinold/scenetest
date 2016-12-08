function getTryingToClockIn(req, resp){
    ClearBlade.init({"request":req});
    posers = doGetTryingToClockIn();
    if (posers === null) {
        log("doGetTryingToClockIn return nill");
        resp.failure("Error trying to fetch posers");
    }
    resp.success(posers);
}