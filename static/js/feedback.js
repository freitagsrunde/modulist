function submitFeedback(moduleID, catID) {

    var obj = {
        category: catID,
        comment: $("#comment-form-" + catID).val()
    };

    $.post("/review/module/" + moduleID + "/add", obj, function(data) {

        if (data.Success) {

            for (var i = 0; i < data.Feedback.length; i++) {
                $("#comment-view-" + catID)[0].innerHTML += "<p>" + data.Feedback[i].Comment + "</p>";
            }
            $("#comment-form-" + catID).val("");
            $("#comment-header-" + catID)[0].innerHTML = "Feedback <span class = \"badge\">" + data.Count + "</span>";
        }
    });
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

function updateAllCounts(moduleID) {

    $.get("/review/module/" + moduleID + "/comments", function(data) {

        if (data.Success) {

            // TODO: Finish this.
            console.log(data.Feedback);
        }
    });
}

$(function() {

    var path = window.location.pathname;
    var moduleID = path.split("/")[3];

    updateAllCounts(moduleID);
})