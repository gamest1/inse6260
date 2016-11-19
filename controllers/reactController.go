// Copyright 2016 Esteban Garro. All rights reserved.

// Package controllers implements the controller layer for the request API.
package controllers

import (
	bc "github.com/goinggo/beego-mgo/controllers/baseController"
	"github.com/goinggo/beego-mgo/services/requestService"
	log "github.com/goinggo/tracelog"
)

//** TYPES

// ReactController manages the API for request related functionality.
type ReactController struct {
	bc.BaseController
}

//** WEB FUNCTIONS

// Index is the initial view for the agent's requests view.
func (controller *ReactController) Index() {
  email := "gamest1@gmail.com"
	log.Startedf(controller.UserID, "ReactController.Index", "email[%s]", email)

	patientRequests, err := requestService.FetchAllRequestsForUser(&controller.Service, email, "p")
	if err != nil {
		log.CompletedErrorf(err, controller.UserID, "ReactController.Index", "FetchAllRequestsForUser[%s]", email)
		controller.ServeError(err)
		return
	}

	controller.Data["Requests"] = patientRequests
	controller.Layout = "shared/react-layout.html"
	controller.TplName = "react/content.html"
	controller.LayoutSections = map[string]string{}
	controller.LayoutSections["PageHead"] = "react/page-head.html"
	controller.LayoutSections["Header"] = "shared/header.html"
}

//** AJAX FUNCTIONS

// RetrieveReact handles the click on an email.
func (controller *ReactController) CreateReact() {

}
