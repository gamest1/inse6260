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
          </tr>
        </thead>
        <tbody>
          {requestRows}
        </tbody>
      </table>
    );
  }
});

var Request = React.createClass({
  render: function() {
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
      </tr>
    );
  }
});

var socket;

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
            self.replaceState({data : newState});
            break;
        }
      }
    });
    socket.on('dbrefresh', this.fetchAllRequests);
  },
  createRoom: function() {
    socket.emit('create', this.props.userid);
    console.log("Room created for this user!");
    this.fetchAllRequests();
  },
  fetchAllRequests: function() {
    console.log("Do your Ajax magic to fetchAllRequests...")
    var self = this;
    $.ajax({
	     method: "GET",
       url: this.props.url + this.props.userid,
       dataType: 'json',
       cache: false,
       success: function(data) {
          var all = [];
          for (let o of data) {
            all.push(new ServiceRequest(o));
          }
          console.log("Got " + all.length + " results");
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
     return {data: []};
  },
  componentDidMount: function() {
   console.log("componentDidMount...");
   this.createRoom();
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

ReactDOM.render(
  <RequestBox  url="http://localhost:9003/requests/" userid="gamest1@gmail.com"/>,
  document.getElementById('content')
);
