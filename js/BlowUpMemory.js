function BlowUpMemory(req, resp){
    logStdErr("BLOW AND GO!")

    ClearBlade.init({request:req});

    var foo = ['a', 'b', 'c', 'd', 'e', 'f', 'g'];
    for (i = 0; i < 100; i++) {
        var bar = ['h', 'i', 'j', 'k', 'l', 'm', 'n'];
        foo.push.apply(foo, bar);
    }
    foo.push.apply(foo, aintNoThing);
    baz.push.apply(foo, aintNoThing);
}
