// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE handle.

// Package main provides sample web application for beego and mgo.
package main

import (
	"net/http"
	"strings"
	"time"
	"os"

	"github.com/astaxie/beego"
	"github.com/goinggo/beego-mgo/localize"
	_ "github.com/goinggo/beego-mgo/routes"
	"github.com/goinggo/beego-mgo/utilities/helper"

	//Support for watchDB:
	"github.com/goinggo/beego-mgo/services"
	"github.com/goinggo/beego-mgo/services/requestService"
	"github.com/goinggo/beego-mgo/models/requestModel"
	"github.com/goinggo/beego-mgo/utilities/mongo"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/mgo.v2/bson"

	//Support for Socket.io
	sock "github.com/googollee/go-socket.io"
	"github.com/zhouhui8915/go-socket.io-client"

	log "github.com/goinggo/tracelog"
)

// watchConfiguration contains settings for running the watch service.
type watchConfiguration struct {
		Watchdb  	 			string
		Watchcollection string
}
//
type simpleObject struct {
	  ID         bson.ObjectId        `bson:"_id,omitempty"`
}
// watchRecord to be used to capture updates to the oplog.rs DB
type watchRecord struct {
		TS				 time.Time 						`bson:"ts" json:"ts"`
		Operation  string    						`bson:"op" json:"op"`
		Collection string	   						`bson:"ns" json:"ns"`
		Object     requestModel.Request `bson:"o"  json:"o"`
		SecObject  simpleObject 				`bson:"o2" json:"02"`
}

var AllReactiveUsers map[string]int = make(map[string]int)
var server *sock.Server

func DatabaseCheckup() {
  log.Started("main", "DatabaseCheckup started...")
	log.Started("main", "Connecting DatabaseCheckup to Socket Server...")

	opts := &socketio_client.Options {
        Transport: "websocket",
        Query:     make(map[string]string),
  }
	opts.Query["log"] = "false"
  opts.Query["origins"] = "*:*"

  uri := "http://localhost:5000/"

  client, err := socketio_client.NewClient(uri, opts)
  if err != nil {
			log.CompletedError(err, "DatabaseCheckup", "Problems creating socket client!")
  }
  client.On("error", func() {
		log.Trace("", "DatabaseCheckup", "Internal Socket Client received an error")
  })
  client.On("connection", func() {
		log.Trace("", "DatabaseCheckup", "Internal Socket Client successfully connected!")
  })
  client.On("disconnection", func() {
		log.Trace("", "DatabaseCheckup", "Internal Socket Client disconnected")
  })


	var Config watchConfiguration
	if err := envconfig.Process("mgo", &Config); err != nil {
		log.CompletedError(err, "DatabaseCheckup", "Watch configuration extraction failed...")
		return
	}

	//var PermanentMongoSession *mgo.Session
	PermanentMongoSession, err := mongo.CopyPermanentSession("")
	if err != nil {
		log.Error(err, "DatabaseCheckup", "Couldn't copy permanent session...")
		return
	}

	//var collection *mgo.Collection (this is our oplog.rs collection!)
	collection := PermanentMongoSession.DB(Config.Watchdb).C(Config.Watchcollection)
	var result watchRecord
  var lastId time.Time
	iter := collection.Find(nil).Sort("$natural").Tail(3 * time.Second)

	defer log.Trace("", "DatabaseCheckup", "Exiting go routine")
	defer iter.Close()
	defer mongo.CloseSession("", PermanentMongoSession)
	defer log.Trace("", "DatabaseCheckup", "DONE!")

	for {
     for iter.Next(&result) {
			 	 //We just care about updates and inserts:
				 if result.Operation == "u" {
					  log.Trace("", "DatabaseCheckup", "Processing update")
					  //var findService hold a service pointer to perform queries on the service_requests db
					 	findService := &services.Service{}
					 	findService.UserID = "DatabaseCheckup subsystem"
					 	MonotonicMongoSession, err := mongo.CopyMonotonicSession("")
					 	if err != nil {
					 		log.Error(err, "DatabaseCheckup subsystem", "Couldn't copy monotonic session")
					 	}
					 	findService.MongoSession = MonotonicMongoSession

						requestID := result.SecObject.ID
						log.Trace("", "DatabaseCheckup", "New update to DB, fetching request...")
						anUpdateRequest, err := requestService.FetchRequest(findService, requestID)
					 	if err != nil {
					 		log.CompletedErrorf(err, findService.UserID, "DatabaseCheckup", "FetchRequest")
					 	} else {
							message := anUpdateRequest.Originator + "::" + anUpdateRequest.Status + "::" + mongo.ToString(anUpdateRequest.ID)
							client.Emit("dbupdate", strings.Replace(message,"\"","",-1))
							log.Trace("", "DatabaseCheckup", "dbupdate: Message emitted %s", message)
						}
				 } else if result.Operation == "i" {
					 log.Trace("", "DatabaseCheckup", "Processing insert")
					 message := result.Object.Originator
					 client.Emit("dbrefresh", message)
					 log.Trace("", "DatabaseCheckup", "dbrefresh: Message emitted %s", message)
				 }
				 lastId = result.TS
     }
     if iter.Err() != nil {
         return
     }
     if iter.Timeout() {
			 	 //log.Trace("", "DatabaseCheckup", "Calling Timeout")
         continue
     }

		 log.Trace("", "DatabaseCheckup", "Polling database")
     query := collection.Find(bson.M{"ts": bson.M{"$gt": lastId}})
     iter = query.Sort("$natural").Tail(3 * time.Second)
	}

	log.Trace("", "DatabaseCheckup", "Infinite for loop completed")
}

func SetupSocketServer() {
  	log.Started("main", "SetupSocketServer started...")
    server, _ = sock.NewServer(nil)
    // if err != nil {
		// 		log.CompletedError(err, "SetupSocketServer", "Couldn't create new SocketIO Server")
    // }
    server.On("connection", func(so sock.Socket) {
				log.Trace("", "SetupSocketServer", "Incoming connection (onConnection)")
				so.On("create", func(userRoom string) {
						log.Trace("", "SetupSocketServer", "Joining user room[%s]",userRoom)
        		so.Join(userRoom)
						_, ok := AllReactiveUsers[userRoom]
						if ok {
								AllReactiveUsers[userRoom] = AllReactiveUsers[userRoom] + 1
						} else {
								AllReactiveUsers[userRoom] = 1
						}
						log.Trace("", "SetupSocketServer", "Current users in room %s: %d",userRoom, AllReactiveUsers[userRoom])
        })
        so.On("dbupdate", func(msg string) {
						log.Trace("", "SetupSocketServer", "Incoming dbupdate message[%+v]",msg)
						msgParts := strings.Split(msg, "::")
						log.Trace("", "SetupSocketServer", "Incoming dbupdate message[%+v]",msgParts)
						//Perform the lookup and broadcast to user room!
						thisRoom := msgParts[0]
						_, ok := AllReactiveUsers[thisRoom]
						if ok {
							  broadcast := msgParts[1] + "::" + msgParts[2]
								log.Trace("", "SetupSocketServer", "Broadcasting message[%s] to room: %s",broadcast, thisRoom)
							  server.BroadcastTo(thisRoom,"dbupdate",broadcast)
						} else {
								log.Trace("", "SetupSocketServer", "Room [%s] not found in %+v", thisRoom, AllReactiveUsers)
						}
        })
				so.On("dbrefresh", func(msg string) {
						log.Trace("", "SetupSocketServer", "Incoming dbrefresh message[%+v]",msg)
						thisRoom := strings.Split(msg, "::")[0]
						_, ok := AllReactiveUsers[thisRoom]
						if ok {
								log.Trace("", "SetupSocketServer", "Broadcasting refresh message to %s",thisRoom)
								server.BroadcastTo(thisRoom,"dbrefresh","")
						}
				})
        so.On("disconnection", func() {
						userRoom := so.Rooms()[0]
					  log.Trace("", "SetupSocketServer", "%s is dropping a connection (onDisconnection)", userRoom)
						thisUserConnections := AllReactiveUsers[userRoom]
						if thisUserConnections > 1 {
								AllReactiveUsers[userRoom] = AllReactiveUsers[userRoom] - 1
								log.Trace("", "SetupSocketServer", "Current users in room %s: %d",userRoom, AllReactiveUsers[userRoom])
						} else if thisUserConnections == 1 {
							log.Trace("", "SetupSocketServer", "%s is dropping its last connection! deleting room", userRoom)
							delete(AllReactiveUsers, userRoom)
						} else {
							log.Trace("", "SetupSocketServer", "**ERROR** this should never happen!")
						}
        })
    })
    server.On("error", func(so sock.Socket, err error) {
				log.CompletedError(err, "SetupSocketServer", "Some SocketIO Server Error")
    })

    http.Handle("/socket.io/", server)
		http.HandleFunc("/", saveHandler)

		//Start the permanent watch on database changes (which will be a client to the Socker Server):
		go DatabaseCheckup()

		log.Trace("", "SetupSocketServer", "Socket Server going up at localhost:5000...")
		http.ListenAndServe(":5000", nil)
		log.Trace("", "SetupSocketServer", "Socket Server went down :(")
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
    // allow cross domain AJAX requests
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:9003")
		requestedFile := r.URL.Path
		if strings.Contains(requestedFile,"socket.io/") {
			log.Trace("", "saveHandler", "Should not be here")
		} else {
			http.ServeFile(w, r, "/Users/gamest/Desktop/GO/src/github.com/goinggo/beego-mgo/static/js" + r.URL.Path)
	  }
}

func main() {
	log.Start(log.LevelTrace)

	// Init mongo
	log.Started("main", "Initializing Mongo")
	err := mongo.Startup(helper.MainGoRoutine)
	if err != nil {
		log.CompletedError(err, helper.MainGoRoutine, "initApp")
		os.Exit(1)
	}

	//Start the Socket.io server:
	go SetupSocketServer()

	// Load message strings
	localize.Init("en-US")

	beego.Run()
	//beego.SetStaticPath("/static","static")

	log.Completed(helper.MainGoRoutine, "Website Shutdown")
	log.Stop()
}
