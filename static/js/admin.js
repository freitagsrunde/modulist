function saveMailHeader() {
  var newtext = {content: $("#mail-header").val()};
  $.post("/admin/send-feedback/mail-header", newtext, function(data) {
    if (data) {
      $("#mailheader-success").fadeIn();
      setTimeout(function() {
        $("#mailheader-success").fadeOut();
      }, 3000)
    }
  })
}

function saveMailFooter() {
  var newtext = {content: $("#mail-footer").val()};
  $.post("/admin/send-feedback/mail-footer", newtext, function(data) {
    if (data) {
      $("#mailfooter-success").fadeIn();
      setTimeout(function() {
        $("#mailfooter-success").fadeOut();
      }, 3000)
    }
  })
}

function sendOutFeedback(depMailAddress, depMailContent) {

    var feedback = {
        mail: depMailAddress,
        content: depMailContent
    }

    $.post("/admin/send-feedback", feedback, function(retData) {

        if(retData.ReturnValue) {
            $("#button-" + depMailAddress).innerHTML = '<p class = "text-info">Versendet!</p>'
        }
    })
}
