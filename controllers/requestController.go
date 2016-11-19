// Copyright 2016 Esteban Garro. All rights reserved.

// Package controllers implements the controller layer for the request API.
package controllers

import (
  "time"
  "strings"
	bc "github.com/goinggo/beego-mgo/controllers/baseController"

  "github.com/goinggo/beego-mgo/services/userService"
	"github.com/goinggo/beego-mgo/services/requestService"
	"github.com/goinggo/beego-mgo/models/requestModel"

  "github.com/goinggo/beego-mgo/utilities/location"
	log "github.com/goinggo/tracelog"
)

//** TYPES

// RequestController manages the API for request related functionality.
type RequestController struct {
	bc.BaseController
}

//** WEB FUNCTIONS

// Index is the initial view for the agent's requests view.
func (controller *RequestController) Index() {
  email := "gamest1@gmail.com"
	log.Startedf(controller.UserID, "RequestController.Index", "email[%s]", email)

	patientRequests, err := requestService.FetchAllRequestsForUser(&controller.Service, email, "p")
	if err != nil {
		log.CompletedErrorf(err, controller.UserID, "RequestController.Index", "FetchAllRequestsForUser[%s]", email)
		controller.ServeError(err)
		return
	}

	controller.Data["Requests"] = patientRequests
	controller.Layout = "shared/basic-layout.html"
	controller.TplName = "requests/content.html"
	controller.LayoutSections = map[string]string{}
	controller.LayoutSections["PageHead"] = "requests/page-head.html"
	controller.LayoutSections["Header"] = "shared/header.html"
	controller.LayoutSections["Modal"] = "shared/modal.html"
}

//** AJAX FUNCTIONS

func (controller *RequestController) RequestsForUser() {
  log.Startedf(controller.UserID, "RequestsForUser", "")
	params := struct {
		Email string `form:":userId" valid:"Required; Email" error:"invalid_user_id"`
	}{controller.GetString(":userId")}

  log.Trace(controller.UserID, "RequestsForUser", "ParseAndValidate, validating params: %+v",params)
	if controller.ParseAndValidate(&params) == false {
		return
	}

  log.Trace(controller.UserID, "RequestsForUser", "Investigating type of user: %+v",params)
  userType, err := userService.TypeForUser(&controller.Service, params.Email)
	if err != nil {
		log.CompletedErrorf(err, controller.UserID, "RequestsForUser", "TypeForUser[%s]", params.Email)
		controller.ServeError(err)
		return
	}

  log.Trace(controller.UserID, "RequestsForUser", "ParseAndValidate, params validated: %s, %s",params.Email,userType)
  patientRequests, err := requestService.FetchAllRequestsForUser(&controller.Service, params.Email, userType)
	if err != nil {
		log.CompletedErrorf(err, controller.UserID, "RequestsForUser", "FetchAllRequestsForUser[%s]", params.Email)
		controller.ServeError(err)
		return
	}

	controller.Data["json"] = patientRequests
	controller.ServeJSON()
}

// RetrieveRequest handles the click on an email.
func (controller *RequestController) CreateRequest() {
  log.Startedf(controller.UserID, "CreateRequest", "")

	var params struct {
		Originator  string   `form:"originator" valid:"Required; Email" error:"invalid_email"`
		StartDate   string   `form:"servicedate" valid:"Required; MinSize(8)" error:"invalid_date"`
		StartTime   string   `form:"servicetime" valid:"Required; MinSize(5)" error:"invalid_time"`
		Apartment   string   `form:"apartment" valid:"MaxSize(3)" error:"invalid_apartment"`
		Number      int      `form:"streetnumber" valid:"Required; Min(0)" error:"invalid_number"`
		Street      string   `form:"streetname" valid:"Required; MaxSize(20)" error:"invalid_street_name"`
    City        string   `form:"cityname" valid:"Required; MaxSize(15)" error:"invalid_city_name"`
    Zip         string   `form:"postalcode" valid:"Required; AlphaNumeric; Match(/([ABCEGHJKLMNPRSTVXY]\d)([ABCEGHJKLMNPRSTVWXYZ]\d){2}/i)" error:"invalid_zip"`
    Duration    int      `form:"duration" valid:"Required; Range(1,3)" error:"invalid_duration"`
    Skill       string   `form:"skill" valid:"Required; Alpha; MaxSize(20)" error:"invalid_skill"`
    Gender      string   `form:"gender" valid:"Required; Alpha; MaxSize(1)" error:"invalid_gender"`
    Languages   []string `form:"languages" valid:"Required; MinSize(1)" error:"invalid_language_selection"`
	}

  log.Trace(controller.UserID, "CreateRequest", "ParseAndValidate, validating params: %+v",params)
	if controller.ParseAndValidate(&params) == false {
		return
	}
  log.Trace(controller.UserID, "CreateRequest", "ParseAndValidate, params validated: %+v",params)

  dateTimeString := params.StartDate + "T" + params.StartTime
  if( len(strings.Split(dateTimeString,":")) < 3 ) {
    dateTimeString = dateTimeString + ":00"
  }

  startTime, err := time.Parse("2006-01-02T15:04:05",dateTimeString)
  if err != nil {
		log.CompletedErrorf(err, controller.UserID, "RequestController.InsertNewRequest", "Error parsing time")
		controller.ServeError(err)
		return
	}

  //Create a new  map[string]interface{} (bson.M) Object using the given params:
  newLocation := &location.Location{0, 0, params.Apartment, params.Number, params.Street, params.City, "", params.Zip}
  newReqs := &requestModel.Requirements{params.Skill, params.Gender, params.Languages, *newLocation}
  newRequest := map[string]interface{}{
      "time": startTime,
      "duration":   params.Duration,
      "status": "pending",
      "request": newReqs,
      "originator": params.Originator,
  }
  //And insert it!
	err = requestService.InsertNewRequest(&controller.Service, newRequest)
	if err != nil {
		log.CompletedErrorf(err, controller.UserID, "RequestController.InsertNewRequest", "newRequest[%v]", newRequest)
		controller.ServeError(err)
		return
	}

  controller.Data["Request"] = newRequest
	controller.Layout = ""
	controller.TplName = "requests/pv_user-detail.html"

  view, _ := controller.RenderString()
	controller.AjaxResponse(0, "SUCCESS", view)
}
