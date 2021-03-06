$(document).ready(function() {
  $('#new-request-button').click(function() {
		DisplayNewRequestForm();
	});
  $('#update-profile-button').click(function() {
		DisplayUpdateProfileForm();
	});
  $('#create-request-button-json').click(function() {
		CreateNewRequest();
	});
});

function DisplayNewRequestForm() {
	try {
		$('#system-modal-title').html("New Request Creation Form");
		$('#system-modal-content').html("<div>Put your form here...</div>");
		$("#systemModal").modal('show');
	}	catch (e) {
        alert(e);
    }
}

function DisplayUpdateProfileForm() {
	try {
		$('#system-modal-title').html("Update Availability Form");
		$('#system-modal-content').html("<div>Put your availability here...</div>");
		$("#systemModal").modal('show');
	}	catch (e) {
        alert(e);
    }
}

$.fn.serializeObject = function() {
    var o = {};
    var a = this.serializeArray();
    $.each(a, function() {
        if (o[this.name] !== undefined) {
            if (!o[this.name].push) {
                o[this.name] = [o[this.name]];
            }
            o[this.name].push(this.value || '');
        } else {
           var e = document.getElementsByName(this.name)[0];
           if (e.type == "number") {
            o[this.name] = parseInt(this.value || '');
           } else {
            o[this.name] = this.value || '';
          }
        }
    });
    return o;
};

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

function CreateNewRequest() {
	try {
        var postData = $('#createRequestForm').serializeObject();
        var service = new ServiceResult();
        service.getJSONData("/requests/createnew",
                            postData,
                            CreateNewRequest_Callback,
                            Standard_ValidationCallback,
                            Standard_ErrorCallback
                            );
    }

    catch (e) {
        alert(e);
    }
}

function CreateNewRequest_Callback() {
	try {
		$('#system-modal-title').html("Request Creation Status");
		$('#system-modal-content').html(this.ResultObject);
		$("#systemModal").modal('show');
	}

	catch (e) {
        alert(e);
    }
}

// Class definition / constructor
var ServiceRequest = function ServiceRequest(apiObject,userType) {
  // Initialization!
  this["Actions"] = userType;
  for (var property in apiObject) {
      if (apiObject.hasOwnProperty(property)) {
        switch(property) {
          case "time":
            var d = new Date(apiObject[property]);
            this["StartTime"] = d.toString();
            break;
          case "request":
            var subObject = apiObject[property];
            var newReqs = {};
            for (var prop in subObject) {
              if (subObject.hasOwnProperty(prop)) {
                switch(prop) {
                  case "gender":
                    if(subObject[prop] == "m") {
                      newReqs[prop] = "Male";
                    } else if(subObject[prop] == "f") {
                      newReqs[prop] = "Female";
                    } else {
                      newReqs[prop] = "N/A";
                    }
                    break;
                  default:
                    newReqs[prop] = subObject[prop];
                }
              }
            }
            this["Requirements"] = newReqs;
            break;
          case "care_giver":
              //Here you could change empty strings for N/A or things like that.
              this["careGiver"] = apiObject[property];
              break;
          default:
            this[property] = apiObject[property];
        }
      }
  }
}

// Instance methods
ServiceRequest.prototype = {
  printLocation: function printLocation() {
      return this.Requirements.location.apartment + "-" + this.Requirements.location.number + " " + this.Requirements.location.street + "\n" +
      this.Requirements.location.city + ", " + this.Requirements.location.state + "  " + this.Requirements.location.zip;
    }
}
