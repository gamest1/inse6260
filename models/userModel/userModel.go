// Copyright 2016 Esteban Garro. All rights reserved.

// Package userModel contains the model for users.
package userModel

import (
	"gopkg.in/mgo.v2/bson"
  "github.com/goinggo/beego-mgo/utilities/location"
  "github.com/goinggo/beego-mgo/utilities/availability"
)

//** TYPES
type (
	// Profile contains information for an individual user.
  Profile struct {
	  FirstName    string    `bson:"first_name" json:"first_name"`
	  LastName     string    `bson:"last_name" json:"last_name"`
	  Gender       string    `bson:"gender" json:"gender"`
	  Languages    []string  `bson:"languages" json:"languages"`
	  Location     location.Location  `bson:"location" json:"location"`
	  Type         string    `bson:"type" json:"type"`
	  Skills       []string  `bson:"skills" json:"skills"`
	  Availability availability.Availability `bson:"availability" json:"availability"`
  }

	// User contains credentials for a user.
  User struct {
    ID       bson.ObjectId `bson:"_id,omitempty"`
  	Email	   string  `bson:"email" json:"email"`
  	Password string  `bson:"password" json:"password"`
  	Profile  Profile `bson:"profile" json:"profile"`
  }
)
