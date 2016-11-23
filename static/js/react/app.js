var RequestList = React.createClass({
  render: function() {
    console.log("Rendering RequestList..." + this.props.data);
    var requestRows = this.props.data.map(function (request,idx) {
      console.log("Request: " + request);
      return (
        <Request key={idx} request={request} / >
      );
    });
    return (
      <table className="table table-bordred table-striped">
        <thead>
          <tr>
            <th>ID</th>
            <th>Date and Time</th>
            <th>Duration</th>
            <th>Status</th>
            <th>Location</th>
            <th>Details</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {requestRows}
        </tbody>
      </table>
    );
  }
});

var CancelButton = React.createClass({
    cancelClick: function () {
      $.ajax({
  	     method: "POST",
         url: this.props.url + this.props.request,
         dataType: 'json',
         cache: false,
         success: function(data) {
              console.log("Sucessfully transmitted cancel for " + this.props.request);
         },
         error: function(xhr, status, err) {
           console.error(this.props.url, status, err.toString());
         }
      });
    },
    render: function() {
          return <button className="btn btn-danger" onClick={this.cancelClick}>Cancel</button>;
    }
});

var CompleteButton = React.createClass({
    completeClick: function () {
      $.ajax({
         method: "POST",
         url: this.props.url + this.props.request,
         dataType: 'json',
         cache: false,
         success: function(data) {
              console.log("Sucessfully transmitted completed message for " + this.props.request);
         },
         error: function(xhr, status, err) {
           console.error(this.props.url, status, err.toString());
         }
      });
    },
    render: function() {
          return <button className="btn btn-success" onClick={this.completeClick}>Complete</button>;
    }
});

var Request = React.createClass({
  render: function() {
    var self = this;
    console.log("Rendering Request...");
    var requestLanguages = this.props.request.Requirements.languages.reduce(function (acc, language, idx, arr) {
      if(idx==arr.length - 1) {
          return acc + ', or ' + language;
      } else {
          return acc  + ', ' + language;
      }
    });
    const colStyle = {
      color: 'blue',
      textAlign: 'center',
    };
    const shortID = this.props.request.ID.substr(this.props.request.ID.length-6,6);
    const loc = this.props.request.printLocation();
    const actions = function(s) {
      if (s == "pending") {
          return <CancelButton request={self.props.request.ID} url="http://localhost:9003/requests/cancel/"/>;
      } else if (s == "allocated") {
          if (self.props.request.Actions == "cg") {
            //Care Givers can Cancel AND mark as Completed!
            return  (
              <div>
              <CancelButton request={self.props.request.ID} url="http://localhost:9003/requests/cancel/"/>
              <CompleteButton request={self.props.request.ID} url="http://localhost:9003/requests/complete/"/>
              </div>
            );
          } else {
            return <CancelButton request={self.props.request.ID} url="http://localhost:9003/requests/cancel/"/>;
          }
      }
    }(this.props.request.status);
    return (
      <tr>
        <td>{shortID}</td>
        <td>{this.props.request.StartTime}</td>
        <td style={colStyle}>{this.props.request.duration}</td>
        <td>{this.props.request.status}</td>
        <td>{loc}</td>
        <td>Request for a {this.props.request.Requirements.Gender} {this.props.request.Requirements.skill} able to speak:<br/>
        {requestLanguages}
        </td>
        <td>{actions}</td>
      </tr>
    );
  }
});

var socket;
var oData;

var RequestBox = React.createClass({
  initSocket: function() {
    //The requests should come from the socket upon connection to this array:
    var self = this;
    console.log("initSocket...");
    //var allowedOrigins = "http://localhost:9003";
    socket = io('http://localhost:5000');
    // , {
    //     allowedOrigins : allowedOrigins,
    //     withCredentials : false
    // });
    //socket = io();
    socket.on('connect', function(){
      console.log("Connection successful!");
    });
    socket.on('disconnect', function(){
      console.log("Reactive connection with server dropped...");
    });
    socket.on('dbupdate', function(message) {
      console.log("Should update " + message);
      var res = message.split("::");
      for (var i = 0; i < self.state.data.length; i++) {
        if (self.state.data[i].ID == res[1]) {
            console.log("Found entry! changing status to " + res[0]);
            var newState = self.state.data.slice();
            newState[i].status = res[0];
            oData = newState;
            self.replaceState({data : newState});
            break;
        }
      }
    });
    socket.on('dbrefresh', this.fetchAllRequests);
  },
  createRoom: function() {
    socket.emit('create', userid);
    console.log("Room created/join for this user/application!");
    this.fetchAllRequests();
  },
  fetchAllRequests: function() {
    var self = this;
    $.ajax({
	     method: "GET",
       url: this.props.url + userid,
       dataType: 'json',
       cache: false,
       success: function(data) {
          var all = [];
          if (data.Requests) {
            for (let o of data.Requests) {
              all.push(new ServiceRequest(o,data.UserType));
            }
            console.log("Got " + all.length + " results");
          } else {
            console.log("Got no LERT table results");
          }
          oData = all;
          self.replaceState({data : all});
       },
       error: function(xhr, status, err) {
         console.error(this.props.url, status, err.toString());
       }
    });
  },
  getInitialState: function() {
     console.log("getInitialState...");
     this.initSocket();
     oData = [];
     return {data: []};
  },
  componentDidMount: function() {
   console.log("componentDidMount...");
   this.createRoom();
   window.requestBox = this;
  },
  restoreDataSet: function() {
    this.replaceState({data : oData});
  },
  filterDataSet: function(start,end) {
    var currentState = this.state.data.slice();
    oData = currentState;
    var newState = currentState.filter(function (el) {
          var d = new Date(el.StartTime);
          return d <= end && d >= start;
    });
    this.replaceState({data : newState});
  },
  render: function() {
    console.log("Rendering RequestBox...");
    return (
      <div className="requestBox">
        <h2>All requests</h2>
        <RequestList data={this.state.data} />
      </div>
    );
  }
});

var userid = document.getElementById('currentUser').innerHTML;
ReactDOM.render(<RequestBox url="http://localhost:9003/requests/" />, document.getElementById('content'));
