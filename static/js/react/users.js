var UserList = React.createClass({
  render: function() {
    console.log("Rendering UserList..." + this.props.data);
    var userRows = this.props.data.map(function (user,idx) {
      console.log("User: " + user);
      return (
        <SysUser key={idx} user={user} / >
      );
    });
    return (
      <table className="table table-bordred table-striped">
        <thead>
          <tr>
            <th>TYPE</th>
            <th>Username</th>
            <th>LastName, FirstName</th>
            <th>Gender</th>
            <th>Location</th>
            <th>Languages</th>
            <th>Skills</th>
          </tr>
        </thead>
        <tbody>
          {userRows}
        </tbody>
      </table>
    );
  }
});

var SysUser = React.createClass({
  render: function() {
    console.log("Rendering User...");
    var userLanguages = this.props.user.Profile.languages.reduce(function (acc, language, idx, arr) {
      if(idx==arr.length - 1) {
          return acc + ', and ' + language;
      } else {
          return acc  + ', ' + language;
      }
    });
    var userSkills = this.props.user.Profile.skills.reduce(function (acc, skill, idx, arr) {
      if(idx==arr.length - 1) {
          return acc + ', and ' + skill;
      } else {
          return acc  + ', ' + skill;
      }
    });
    const colStyle = {
      color: 'blue',
      textAlign: 'center',
    };
    const loc = this.props.user.printLocation();
    return (
      <tr>
        <td>{this.props.user.Profile.type}</td>
        <td>{this.props.user.email}</td>
        <td>{this.props.user.Profile.last_name}, {this.props.user.Profile.first_name}</td>
        <td>{this.props.user.Profile.gender}</td>
        <td>{loc}</td>
        <td>{userLanguages}</td>
        <td>{userSkills}</td>
      </tr>
    );
  }
});

var usocket;

var UserBox = React.createClass({
  initSocket: function() {
    //The requests should come from the socket upon connection to this array:
    var self = this;
    console.log("initSocket...");
    //var allowedOrigins = "http://localhost:9003";
    usocket = io('http://localhost:5000');
    // , {
    //     allowedOrigins : allowedOrigins,
    //     withCredentials : false
    // });
    //socket = io();
    usocket.on('connect', function(){
      console.log("Connection successful!");
    });
    usocket.on('disconnect', function(){
      console.log("Reactive connection with server dropped...");
    });
    usocket.on('urefresh', this.fetchAllUsers);
  },
  createRoom: function() {
    usocket.emit('create', cuserid);
    console.log("Room created/join for this user/application!");
    this.fetchAllUsers();
  },
  fetchAllUsers: function() {
    var self = this;
    $.ajax({
	     method: "GET",
       url: this.props.url + cuserid,
       dataType: 'json',
       cache: false,
       success: function(data) {
          var all = [];
          if (data) {
            for (let o of data) {
              all.push(new SystemUser(o));
            }
            console.log("Got " + all.length + " results");
          } else {
            console.log("Got no user table results");
          }
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
    console.log("Rendering UserBox...");
    return (
      <div className="userBox">
        <h2>All System Users</h2>
        <UserList data={this.state.data} />
      </div>
    );
  }
});

var cuserid = document.getElementById('currentUser').innerHTML;
ReactDOM.render(<UserBox url="http://localhost:9003/user/display/" />, document.getElementById('userslist'));
