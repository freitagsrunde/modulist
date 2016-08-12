function checkModule(id) {

    $.get("/done/" + id, function(data) {

        if (data.Done !== undefined) {

            if (data.Done) {
                $("#module-"+data.ID+" > td > a > span").removeClass("glyphicon-unchecked")
                $("#module-"+data.ID+" > td > a > span").addClass("glyphicon-check")
            } else {
                $("#module-"+data.ID+" > td > a > span").removeClass("glyphicon-check")
                $("#module-"+data.ID+" > td > a > span").addClass("glyphicon-unchecked")
            }

            updateVisibility();
        }
    })
}

function updateVisibility(fast) {

    rows = $("table tr")

    if ($("#modules-list-approved-checkbox")[0].checked) {

        for (var i = 1; i < rows.length; i++) {
            row = $($("table tr")[i])

            if (row.find("a span").hasClass("glyphicon-check")) {

                if (fast) {
                    row.hide()
                } else {
                    row.fadeOut()
                }
            }
        }

        createCookie("hideDone", 1)
    } else {

        for (var i = 1; i < rows.length; i++) {
            row = $($("table tr")[i])

            if (fast) {
                row.show()
            } else {
                row.fadeIn()
            }
        }

        createCookie("hideDone", 0)
    }
}

function feedbackCheck(id) {

    $.get("/done/" + id, function(data) {

        if (data.Done !== undefined) {

            if (data.Done) {
                $("h2")[0].innerHTML += '<sup><span class="badge">DONE</span></sup>'
                $($("h2")[1]).next().find("a").text("Modul als unfertig markieren")
            } else {
                $("h2 > sup").remove()
                $($("h2")[1]).next().find("a").text("Modul als fertig markieren")
            }
        }
    })
}

function createCookie(name, value, days) {

    var expires;

    if (days) {
        var date = new Date();
        date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
        expires = "; expires=" + date.toGMTString();
    } else {
        expires = "";
    }

    document.cookie = encodeURIComponent(name) + "=" + encodeURIComponent(value) + expires + "; path=/";
}

function readCookie(name) {

    var nameEQ = encodeURIComponent(name) + "=";
    var ca = document.cookie.split(';');

    for (var i = 0; i < ca.length; i++) {

        var c = ca[i];
        while (c.charAt(0) === ' ') {
            c = c.substring(1, c.length);
        }

        if (c.indexOf(nameEQ) === 0) {
            return decodeURIComponent(c.substring(nameEQ.length, c.length));
        }
    }

    return null;
}

$(function() {

    if (readCookie("hideDone") === "1") {
        $("#modules-list-approved-checkbox")[0].checked=true
        updateVisibility(true);
    }
})
