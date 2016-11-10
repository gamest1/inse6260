// Copyright 2016 Esteban Garro. All rights reserved.

// Package controllers implements the controller layer for the user API.
package controllers

import (
	bc "github.com/goinggo/beego-mgo/controllers/baseController"
	"github.com/goinggo/beego-mgo/services/userService"
	log "github.com/goinggo/tracelog"
)

//** TYPES

// UserController manages the API for user related functionality.
type UserController struct {
	bc.BaseController
}

//** WEB FUNCTIONS

// Index is the initial view for the agent's users view.
func (controller *UserController) Index() {
	kind := "cg"
	log.Startedf(controller.UserID, "UserController.Index", "Kind[%s]", kind)

	careGivers, err := userService.FindUsersOfKind(&controller.Service, kind)
	if err != nil {
		log.CompletedErrorf(err, controller.UserID, "UserController.Index", "FindUsersOfKind[%s]", kind)
		controller.ServeError(err)
		return
	}

	skills, err := userService.FetchAllSkills(&controller.Service)
	if err != nil {
		log.CompletedErrorf(err, controller.UserID, "UserController.Index", "FetchAllSkills")
		controller.ServeError(err)
		return
	}

	languages, err := userService.FetchAllLanguagesForKind(&controller.Service, kind)
	if err != nil {
		log.CompletedErrorf(err, controller.UserID, "UserController.Index", "FetchAllLanguagesForKind[%s]", kind)
		controller.ServeError(err)
		return
	}


	controller.Data["CareGivers"] = careGivers
	controller.Data["Skills"] = skills
	controller.Data["Languages"] = languages
	controller.Layout = "shared/basic-layout.html"
	controller.TplName = "users/content.html"
	controller.LayoutSections = map[string]string{}
	controller.LayoutSections["PageHead"] = "users/page-head.html"
	controller.LayoutSections["Header"] = "shared/header.html"
	controller.LayoutSections["Modal"] = "shared/modal.html"
}

//** AJAX FUNCTIONS

// RetrieveUser handles the click on an email.
func (controller *UserController) RetrieveUser() {
	var params struct {
		UserEmail string `form:"email" valid:"Required; MinSize(4)" error:"invalid_email"`
	}

	if controller.ParseAndValidate(&params) == false {
		return
	}

	profile, err := userService.FetchProfile(&controller.Service, params.UserEmail)
	if err != nil {
		log.CompletedErrorf(err, controller.UserID, "UserController.RetrieveUser", "UserEmail[%s]", params.UserEmail)
		controller.ServeError(err)
		return
	}

	controller.Data["UserProfile"] = profile
	controller.Layout = ""
	controller.TplName = "users/pv_user-detail.html"
	view, _ := controller.RenderString()

	controller.AjaxResponse(0, "SUCCESS", view)
}
