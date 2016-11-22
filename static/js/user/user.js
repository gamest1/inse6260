$(document).ready(function() {
	$('.detail').click(function(e) {
		e.preventDefault();
		ShowDetail(this);
	});
	$('#usedate').change(function(e) {
    if ($(this).is(':checked')) {
      $("#daterange").show();
    } else {
      $("#daterange").hide();
			window.requestBox.restoreDataSet();
    }
	});
	$('input[name="daterange"]').daterangepicker({
    locale: {
      format: 'YYYY-MM-DD'
    },
    startDate: '2016-01-01',
    endDate: '2016-12-31'
	}, function(start, end, label) {
			var sDate = new Date(start);
			var eDate = new Date(end);
			console.log("A new date range was chosen: " + sDate + ' to ' + eDate);
			//console.log("A new date range was chosen: " + start.format('YYYY-MM-DD') + ' to ' + end.format('YYYY-MM-DD'));
			window.requestBox.filterDataSet(sDate,eDate);
	});
});

function Standard_Callback() {
    try {
        alert(this.ResultString);
    }

    catch (e) {
        alert(e);
    }
}

function Standard_ValidationCallback() {
    try {
        alert(this.ResultString);
    }

    catch (e) {
        alert(e);
    }
}

function Standard_ErrorCallback() {
    try {
        alert(this.ResultString);
    }

    catch (e) {
        alert(e);
    }
}

function ShowDetail(result) {
	try {
		var postData = {};
		postData["email"] = $(result).attr('data');

        var service = new ServiceResult();
        service.getJSONData("/users/retrieveuser",
                            postData,
                            ShowDetail_Callback,
                            Standard_ValidationCallback,
                            Standard_ErrorCallback
                            );
    }

    catch (e) {
        alert(e);
    }
}

function ShowDetail_Callback() {
	try {
		$('#system-modal-title').html("User Details");
		$('#system-modal-content').html(this.ResultObject);
		$("#systemModal").modal('show');
	}

	catch (e) {
        alert(e);
    }
}


// Class definition / constructor
var SystemUser = function SystemUser(apiObject) {
  // Initialization!
  for (var property in apiObject) {
      if (apiObject.hasOwnProperty(property)) {
        switch(property) {
          case "email":
            this["email"] = apiObject[property];
            break;
          case "profile":
            var subObject = apiObject[property];
            var newProfs = {};
            for (var prop in subObject) {
              if (subObject.hasOwnProperty(prop)) {
                switch(prop) {
                  case "gender":
                    if(subObject[prop] == "m") {
                      newProfs[prop] = "Male";
                    } else if(subObject[prop] == "f") {
                      newProfs[prop] = "Female";
                    } else {
                      newProfs[prop] = "N/A";
                    }
                    break;
									case "skills":
										console.log("Dealing with skills " + subObject[prop]);
										if (subObject[prop] && subObject[prop].length > 0) {
											newProfs[prop] = subObject[prop];
										} else {
											console.log("Replacing null or empty array");
											newProfs[prop] = ["N/A"];
										}
										break;
									case "location":
										console.log("Dealing with location");
										newProfs[prop] = subObject[prop];
										break;
                  default:
                    newProfs[prop] = subObject[prop];
                }
              }
            }
            this["Profile"] = newProfs;
            break;
          default:
            this[property] = apiObject[property];
        }
      }
  }
}

// Instance methods
SystemUser.prototype = {
  printLocation: function printLocation() {
      return this.Profile.location.apartment + "-" + this.Profile.location.number + " " + this.Profile.location.street + "\n" +
      this.Profile.location.city + ", " + this.Profile.location.state + "  " + this.Profile.location.zip;
    }
}
