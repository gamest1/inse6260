// Copyright 2016 Esteban Garro. All rights reserved.

// Package requestService implements the service for the request functionality.
package requestService

import (
  "time"
  "errors"
	"github.com/goinggo/beego-mgo/models/requestModel"
	"github.com/goinggo/beego-mgo/services"
	"github.com/goinggo/beego-mgo/utilities/helper"
	"github.com/goinggo/beego-mgo/utilities/mongo"
	log "github.com/goinggo/tracelog"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//** TYPES

type (
	// requestConfiguration contains settings for running the request service.
	requestsConfiguration struct {
		Database   string
		Collection string
	}
)

//** PACKAGE VARIABLES

// Config provides requests configuration from the environment variables via the envconfig package
var Config requestsConfiguration

//** INIT

func init() {
	// Pull in the configuration.
	if err := envconfig.Process("requests", &Config); err != nil {
		log.CompletedError(err, helper.MainGoRoutine, "Init")
	}
}

//** PUBLIC FUNCTIONS
func AllocateRequest(service *services.Service, requestID string, careGiver string) error {
	log.Startedf(service.UserID, "AllocateRequest", "requestID[%s] to %s", requestID, careGiver)
  f := func(collection *mgo.Collection) error {
		selectorMap := bson.M{"_id": bson.ObjectIdHex(requestID)}

    if careGiver == "" {
      updateMap := bson.M{"$set": bson.M{"status": "pending", "care_giver" : ""}}
  		log.Trace(service.UserID, "AllocateRequest", "MGO : db.%s.update(%s,%s)", Config.Collection, mongo.ToString(selectorMap),mongo.ToString(updateMap))
  		return collection.Update(selectorMap,updateMap)
    } else {
      updateMap := bson.M{"$set": bson.M{"status": "allocated", "care_giver" : careGiver}}
      log.Trace(service.UserID, "AllocateRequest", "MGO : db.%s.update(%s,%s)", Config.Collection, mongo.ToString(selectorMap),mongo.ToString(updateMap))
      return collection.Update(selectorMap,updateMap)
    }
	}

  if err := service.DBAction(Config.Database, Config.Collection, f); err != nil {
		if err != mgo.ErrNotFound {
			log.CompletedError(err, service.UserID, "AllocateRequest")
			return err
		}
	}

	log.Completedf(service.UserID, "AllocateRequest", "request successfully allocated!")
	return nil
}

func UpdateRequest(service *services.Service, requestID string, status string) error {
	log.Startedf(service.UserID, "UpdateRequest", "requestID[%s]: %s", requestID, status)
  f := func(collection *mgo.Collection) error {

		selectorMap := bson.M{"_id": bson.ObjectIdHex(requestID)}
    updateMap := bson.M{"$set": bson.M{"status": status}}

		log.Trace(service.UserID, "UpdateRequest", "MGO : db.%s.update(%s,%s)", Config.Collection, mongo.ToString(selectorMap),mongo.ToString(updateMap))
		return collection.Update(selectorMap,updateMap)
	}

  if err := service.DBAction(Config.Database, Config.Collection, f); err != nil {
		if err != mgo.ErrNotFound {
			log.CompletedError(err, service.UserID, "UpdateRequest")
			return err
		}
	}

	log.Completedf(service.UserID, "UpdateRequest", "request successfully updated!")
	return nil
}

// InsertNewRequest adds a Request object to mongoDB
func InsertNewRequest(service *services.Service, request map[string]interface{}) error {
	log.Startedf(service.UserID, "InsertNewRequest", "request:%+v", request)
  f := func(collection *mgo.Collection) error {
		queryMap := bson.M(request)

		log.Trace(service.UserID, "InsertNewRequest", "MGO : db.%s.insert(%s)", Config.Collection, mongo.ToString(queryMap))
		return collection.Insert(queryMap)
	}

  if err := service.DBAction(Config.Database, Config.Collection, f); err != nil {
		if err != mgo.ErrNotFound {
			log.CompletedError(err, service.UserID, "InsertNewRequest")
			return err
		}
	}

	log.Completedf(service.UserID, "InsertNewRequest", "new request object successfully injected!")
	return nil
}

// FetchAllRequestsFromDateToDate retrieves the profile of all the requests between startTime and endTime
func FetchAllRequestsFromDateToDate(service *services.Service, startTime time.Time, endTime time.Time) ([]requestModel.Request, error) {
  if endTime.Before(startTime) {
    err := errors.New("FetchAllRequestsFromDateToDate: endTime before stratTime")
    return nil, err
  }
	log.Startedf(service.UserID, "FetchAllRequestsFromDateToDate", "startTime[%s]-endTime[%s]", startTime.String(), endTime.String())

	var requests []requestModel.Request
	f := func(collection *mgo.Collection) error {
		queryMap := bson.M{ "time" : bson.M{ "$gte" : startTime, "$lt" : endTime}}

		log.Trace(service.UserID, "FetchAllRequestsFromDateToDate", "MGO : db.%s.find(%s)", Config.Collection, mongo.ToString(queryMap))
		return collection.Find(queryMap).All(&requests)
	}

	if err := service.DBAction(Config.Database, Config.Collection, f); err != nil {
		if err != mgo.ErrNotFound {
			log.CompletedError(err, service.UserID, "FetchAllRequestsFromDateToDate")
			return nil, err
		}
	}

	log.Completedf(service.UserID, "FetchAllRequestsFromDateToDate", "requests on interval startTime[%s]-endTime[%s]: %+v", startTime.String(), endTime.String(), requests)
	return requests, nil
}

// FetchAllRequestsForUser retrieves the requests generated by or assigned to a user identified by its email
func FetchAllRequestsForUser(service *services.Service, email string, kind string) ([]requestModel.Request, error) {
	log.Startedf(service.UserID, "FetchAllRequestsForUser", "email[%s]", email)

	var requests []requestModel.Request
	f := func(collection *mgo.Collection) error {
		queryMap := bson.M{}
    if kind == "cg" {
      queryMap["care_giver"] = email
    } else if kind == "p" {
      queryMap["originator"] = email
    }

		log.Trace(service.UserID, "FetchAllRequestsForUser", "Query : db.%s.find(%s)", Config.Collection, mongo.ToString(queryMap))
		return collection.Find(queryMap).All(&requests)
	}

	if err := service.DBAction(Config.Database, Config.Collection, f); err != nil {
		log.CompletedError(err, service.UserID, "FetchAllRequestsForUser")
		return nil, err
	}

	log.Completedf(service.UserID, "FetchAllRequestsForUser", "requests related to %s %+v", email, requests)
	return requests, nil
}

//To be used by the DatabaseCheckup Socket Server:
// FetchRequest retrieves a service request object based on ID:
func FetchRequest(service *services.Service, ID bson.ObjectId) (*requestModel.Request, error) {
	log.Startedf(service.UserID, "FetchRequest", "ID[%+v]", ID)

	var request *requestModel.Request
	f := func(collection *mgo.Collection) error {
		queryMap := bson.M{ "_id" : ID}

		log.Trace(service.UserID, "FetchRequest", "MGO : db.%s.find(%s)", Config.Collection, mongo.ToString(queryMap))
		return collection.Find(queryMap).One(&request)
	}

	if err := service.DBAction(Config.Database, Config.Collection, f); err != nil {
			log.CompletedError(err, service.UserID, "FetchRequest")
			return nil, err
	}

	log.Completedf(service.UserID, "FetchRequest", "Request found: %+v", request)
	return request, nil
}

// To be used by the Scheduler:
// FetchCurrentSchedule retrieves all the requests with pending or allocated status:
func FetchCurrentSchedule(service *services.Service) ([]requestModel.Request, error) {
	log.Startedf(service.UserID, "FetchCurrentSchedule", "")

	var requests []requestModel.Request
	f := func(collection *mgo.Collection) error {
		queryMap := bson.M{"$or": []bson.M{bson.M{"status": "pending"}, bson.M{"status": "allocated"}}}

		log.Trace(service.UserID, "FetchCurrentSchedule", "MGO : db.%s.find(%s)", Config.Collection, mongo.ToString(queryMap))
		return collection.Find(queryMap).All(&requests)
	}

	if err := service.DBAction(Config.Database, Config.Collection, f); err != nil {
		if err != mgo.ErrNotFound {
			log.CompletedError(err, service.UserID, "FetchCurrentSchedule")
			return nil, err
		}
	}

	log.Completedf(service.UserID, "FetchCurrentSchedule", "Current schedule fetch completed %+v", requests)
	return requests, nil
}
