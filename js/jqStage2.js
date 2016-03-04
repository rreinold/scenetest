function Stage2(req, resp) {
	logStdErr("stage 2!");

	ClearBlade.init({request:req});
    var messaging = ClearBlade.Messaging();
	messaging.publish("/clearblade/test/jobQueue/stage2/results", "Adam");

    resp.success("OK");
}
