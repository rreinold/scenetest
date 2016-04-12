function updateUser(req, resp) {
    ClearBlade.init({request:req});


    var q = ClearBlade.Query();


    q.equalTo("email", "foo@bar.com");


    ClearBlade.User().setUsers(q, {"dafizzle": "WOW"});


    resp.success("I think it worked");
}
