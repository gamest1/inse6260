$(document).ready(function() {

	var weekdays = ["monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"];
  var upperBounds = [];
  var lowerBounds = [];

	for (var i = 0; i < 7 ; i++ ) {
		var day = weekdays[i];
		var sel = "#" + day + "-slider-range";
    var upperValue = 660;
    var lowerValue = 660;
    var lookingOnes = true;
		var done = false;
		var sliderValue = $("#" + day).val()

    for (var j = 23; j >= 1 ; j-- ) {
      cursor = sliderValue / Math.pow(2, j);
      if (cursor > 1 && lookingOnes) {
        upperValue = 60 * (j+1);
        lookingOnes = false;
        sliderValue -=  Math.pow(2, j);
      } else if (cursor > 1 && !lookingOnes) {
        sliderValue -=  Math.pow(2, j);
      } else if (cursor == 1 && !lookingOnes) {
        lowerValue = 60 * j;
				done = true
        break;
      }
    }

		if (!lookingOnes && !done) {
			lowerValue = 0;
		}

		$(sel).slider({
		    range: true,
		    min: 0,
		    max: 1440,
		    step: 60,
		    values: [lowerValue, upperValue],
		    slide: slideUpdate
		});

    upperBounds.push(upperValue);
    lowerBounds.push(lowerValue);
	}

  for (var i = 0; i < 7 ; i++ ) {
		var dia = weekdays[i];
    TimeLabelUpdate(dia, lowerBounds[i], upperBounds[i])
  }
});
