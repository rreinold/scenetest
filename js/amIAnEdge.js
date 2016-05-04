function amIAnEdge(req, resp) {
    ClearBlade.init({request:req});
    if (ClearBlade.isEdge() === true) {
        resp.success("We are an edge: " + ClearBlade.edgeId());
    } else {
        resp.success("We are NOT an edge");
    }
}
