function Stage3(req, resp) {
	logStdErr("stage 3!!");

	ClearBlade.init({request:req});
	var collection = ClearBlade.Collection({collectionName:"CityInfo"});
	collection.create({name:req["params"]["body"]}, function(err,val){
		if (err) {
			logStdErr(val);
		}
		resp.success();
	});
}
