function updateClockedInPortal(req, resp){
    ClearBlade.init({"request":req});
    notifyUpdateClockInPortal();
    resp.success("");
}