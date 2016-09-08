function submitFeedback(id) {

    var obj = {
        comment: $("#comment-" + id).val()
    };

    $.post("addFeedback/" + id, obj, function(data) {

        if (data.ReturnValue) {
            $("#comment-body-" + id)[0].innerHTML += '<div style="float: right"><span onclick="deleteFeedback(' + data.SavedAsID + ')" style="cursor:pointer; color:#bbb;" class="glyphicon glyphicon-trash"></span></div><b>' + data.UserName + " schrieb:</b><p>" + $("#comment-" + id).val() + "</p><hr>";
            $("#comment-" + id).val("")
            updateBadges();
        }
    });
}

function updateBadges() {

    pnlgroup = $(".panel-group")

    for (var i = 0; i < pnlgroup.length; i++) {

        pnl = $(pnlgroup[i])
        count = pnl.find("hr").length

        if (count > 0) {
            pnl.find("h4 > a")[0].innerHTML = '<span class="badge">' + count + '</span> Zeige Feedback'
        } else {
            pnl.find("h4 > a")[0].innerHTML = 'Zeige Feedback'
        }
    }
}

function deleteFeedback(id) {

    if (confirm("Soll das abgegebene Feedback wirklich gelöscht werden?")) {

        $.get("/deleteFeedback/" + id, function(data) {
            if (data === true) {
                location.reload();
            }
        });
    }
}

$(function() {
    updateBadges()
});