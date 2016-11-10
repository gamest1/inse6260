// Copyright 2016 Esteban Garro. All rights reserved.

// Package requestModel contains the model for all system Service Requests.
package requestModel

import (
  "time"
  "github.com/goinggo/beego-mgo/utilities/location"
	"gopkg.in/mgo.v2/bson"
)

//** TYPES
type (
  Requirements struct {
    Skill     string   `bson:"skill" json:"skill"`
    Gender    string   `bson:"gender" json:"gender"`
    Languages []string `bson:"languages" json:"languages"`
    Location  location.Location `bson:"location" json:"location"`
  }

  Request struct {
	  ID         bson.ObjectId        `bson:"_id,omitempty"`
	  StartTime  time.Time            `bson:"time" json:"time"`
	  Duration   int                  `bson:"duration" json:"duration"`
	  Status     string               `bson:"status" json:"status"`
	  Requirements Requirements 			`bson:"request" json:"request"`
	  Originator string  `bson:"originator" json:"originator"`
	  CareGiver  string  `bson:"care_giver" json:"care_giver"`
  }
)
