$(document).ready(function() {
	$('.detail').click(function(e) {
		e.preventDefault();
		ShowDetail(this);
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
