function Stage1(req, resp) {
	logStdErr("stage 1");

	ClearBlade.init({request:req});
    var messaging = ClearBlade.Messaging();
	messaging.publish("/clearblade/test/jobQueue/stage1/results", "Adam");

    resp.success("OK");
}
