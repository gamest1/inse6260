$(document).ready(function() {
  $('#location').change(function(e) {
    if ($(this).is(':checked')) {
      $("#newlocation").show();
    } else {
      $("#newlocation").hide();
    }
	});
});
