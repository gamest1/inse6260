$(document).ready(function() {

	$('#type').change(function(e) {
    if ($(this).val()=="cg") {
      $("#extendedSkills").show();
    } else {
      $("#extendedSkills").hide();
    }
		if ($(this).val()=="a") {
      $("#invitationToken").show();
    } else {
      $("#invitationToken").hide();
    }
	});
	$('#skills').change(function(e) {
    if ($(this).val()=="other") {
      $("#otherskilldiv").show();
    } else {
      $("#otherskilldiv").hide();
    }
	});
  $('#languages').change(function(e) {
    if ($(this).val()=="other") {
      $("#otherlanguagediv").show();
    } else {
      $("#otherlanguagediv").hide();
    }
  });

	var weekdays = ["monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"];
	for (var i = 0; i < 7 ; i++ ) {
		var day = weekdays[i];
		var sel = "#" + day + "-slider-range";
		var f = function (e, ui) {
					  var h1 = Math.floor(ui.values[0] / 60);
		        var hours1 = h1;
		        var minutes1 = ui.values[0] - (hours1 * 60);

		        if (hours1.length == 1) hours1 = '0' + hours1;
		        if (minutes1.length == 1) minutes1 = '0' + minutes1;
		        if (minutes1 == 0) minutes1 = '00';
		        if (hours1 >= 12) {
		            if (hours1 == 12) {
		                hours1 = hours1;
		                minutes1 = minutes1 + " PM";
		            } else {
		                hours1 = hours1 - 12;
		                minutes1 = minutes1 + " PM";
		            }
		        } else {
		            hours1 = hours1;
		            minutes1 = minutes1 + " AM";
		        }
		        if (hours1 == 0) {
		            hours1 = 12;
		            minutes1 = minutes1;
		        }

						var dia = this.id.split("-")[0];
						var sel2 = "." + dia + "-slider-time";
		        $(sel2).html(hours1 + ':' + minutes1);

		        var h2 = Math.floor(ui.values[1] / 60);
		        var hours2 = h2;
		        var minutes2 = ui.values[1] - (hours2 * 60);

		        if (hours2.length == 1) hours2 = '0' + hours2;
		        if (minutes2.length == 1) minutes2 = '0' + minutes2;
		        if (minutes2 == 0) minutes2 = '00';
		        if (hours2 >= 12) {
		            if (hours2 == 12) {
		                hours2 = hours2;
		                minutes2 = minutes2 + " PM";
		            } else if (hours2 == 24) {
		                hours2 = 11;
		                minutes2 = "59 PM";
		            } else {
		                hours2 = hours2 - 12;
		                minutes2 = minutes2 + " PM";
		            }
		        } else {
		            hours2 = hours2;
		            minutes2 = minutes2 + " AM";
		        }

		        $(sel2 + "2").html(hours2 + ':' + minutes2);

						var acc = 0;
						for (var i = h1; i < h2; i++) {
							acc += Math.pow(2, i);
						}

						$("#" + dia).val(acc);
		};

		console.log($("#" + day).val())
		$(sel).slider({
		    range: true,
		    min: 0,
		    max: 1440,
		    step: 60,
		    values: [600, 720],
		    slide: f
		});
	}
});
