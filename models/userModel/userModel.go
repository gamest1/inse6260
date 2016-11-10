// Copyright 2016 Esteban Garro. All rights reserved.

// Package userModel contains the model for users.
package userModel

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
  "github.com/goinggo/beego-mgo/utilities/location"
)

//** TYPES
type (
  //Availability and Location may have to move to their own package:
  Availability struct {
    Monday    int `bson:"monday" json:"monday"`
    Tuesday   int `bson:"tuesday" json:"tuesday"`
    Wednesday int `bson:"wednesday" json:"wednesday"`
    Thursday  int `bson:"thursday" json:"thursday"`
    Friday    int `bson:"friday" json:"friday"`
    Saturday  int `bson:"saturday" json:"saturday"`
    Sunday    int `bson:"sunday" json:"sunday"`
  }

	// Profile contains information for an individual user.
  Profile struct {
	  FirstName    string    `bson:"first_name" json:"first_name"`
	  LastName     string    `bson:"last_name" json:"last_name"`
	  Gender       string    `bson:"gender" json:"gender"`
	  Languages    []string  `bson:"languages" json:"languages"`
	  Location     location.Location  `bson:"location" json:"location"`
	  Type         string    `bson:"type" json:"type"`
	  Skills       []string  `bson:"skills" json:"skills"`
	  Availability Availability `bson:"availability" json:"availability"`
  }

	// User contains credentials for a user.
  User struct {
    ID       bson.ObjectId `bson:"_id,omitempty"`
  	Email	   string  `bson:"email" json:"email"`
  	Password string  `bson:"password" json:"password"`
  	Profile  Profile `bson:"profile" json:"profile"`
  }
)

// DisplayAvailability displays availability in human-readable format.
func (availability *Availability) DisplayAvailability() string {
	return fmt.Sprintf("Hours:\n\nMonday %d\nTuesday %d\nWednesday %d...", availability.Monday, availability.Tuesday, availability.Wednesday)
}
