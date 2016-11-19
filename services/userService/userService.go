// Copyright 2016 Esteban Garro. All rights reserved.

// Package userService implements the service for the user functionality.
package userService

import (
	"github.com/goinggo/beego-mgo/models/userModel"
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
	// userConfiguration contains settings for running the user service.
	usersConfiguration struct {
		Database   string
		Collection string
	}
)

//** PACKAGE VARIABLES

// Config provides users configuration from the environment variables via the envconfig package
var Config usersConfiguration

//** INIT

func init() {
	// Pull in the configuration.
	if err := envconfig.Process("users", &Config); err != nil {
		log.CompletedError(err, helper.MainGoRoutine, "Init")
	}
}

//** PUBLIC FUNCTIONS
// InsertNewUser adds a Request object to mongoDB
func InsertNewUser(service *services.Service, newUser map[string]interface{}) error {
	log.Startedf(service.UserID, "InsertNewUser", "newUser:%+v", newUser)
  f := func(collection *mgo.Collection) error {
		queryMap := bson.M(newUser)

		log.Trace(service.UserID, "InsertNewUser", "MGO : db.%s.insert(%s)", Config.Collection, mongo.ToString(queryMap))
		return collection.Insert(queryMap)
	}

  if err := service.DBAction(Config.Database, Config.Collection, f); err != nil {
		if err != mgo.ErrNotFound {
			log.CompletedError(err, service.UserID, "InsertNewUser")
			return err
		}
	}

	log.Completedf(service.UserID, "InsertNewUser", "new user objected successfully injected!")
	return nil
}


func UpdateUserAvailability(service *services.Service, newUser map[string]interface{}) error {
	log.Startedf(service.UserID, "UpdateUserAvailability", "newUser[%+v]", newUser)

	f := func(collection *mgo.Collection) error {
		if newUser["currentpassword"] != nil {
			queryMap := bson.M{"$set":bson.M{"password":newUser["newpassword"],"profile.availability":newUser["profile"].(*userModel.Profile).Availability}}
			log.Trace(service.UserID, "UpdateUserAvailability", "MGO 1: db.%s.update(%s)", Config.Collection, mongo.ToString(queryMap))
			return collection.Update( bson.M{"email": newUser["email"], "password": newUser["currentpassword"]} , queryMap)
		} else {
			queryMap := bson.M{"$set":bson.M{"profile.availability":newUser["profile"].(*userModel.Profile).Availability}}
			log.Trace(service.UserID, "UpdateUserAvailability", "MGO 2: db.%s.update(%s)", Config.Collection, mongo.ToString(queryMap))
			return collection.Update( bson.M{"email": newUser["email"]} , queryMap)
		}
		return nil
	}

	if err := service.DBAction(Config.Database, Config.Collection, f); err != nil {
		if err != mgo.ErrNotFound {
			log.CompletedError(err, "UpdateUserAvailability", "Database update failed")
		} else {
			log.CompletedError(err, "UpdateUserAvailability", "User not found!")
		}
		return err
	}

	log.Completedf(service.UserID, "UpdateUserAvailability", "user profile: %+v", newUser)
	return nil
}

// Login returns an error if an email-password combination doesn't exist in the database or returns the users profile:
func Login(service *services.Service, email string, password string) (*userModel.Profile, error) {
	log.Startedf(service.UserID, "Login", "email[%s]", email)

	var user userModel.User
	f := func(collection *mgo.Collection) error {
		queryMap := bson.M{"email": email, "password": password}

		log.Trace(service.UserID, "Login", "MGO : db.%s.find(%s,{\"profile\": 1}).limit(1)", Config.Collection, mongo.ToString(queryMap))
		return collection.Find(queryMap).Select(bson.M{"profile": 1}).One(&user)
	}

	if err := service.DBAction(Config.Database, Config.Collection, f); err != nil {
		if err != mgo.ErrNotFound {
			log.CompletedError(err, "Login", "Database find failed")
		} else {
			log.CompletedError(err, "Login", "User not found!")
		}
		return nil, err
	}

	log.Completedf(service.UserID, "Login", "user profile: %+v", &user.Profile)
	return &user.Profile, nil
}

// FetchAllLanguagesForKind retrieves a list of all languages spoken by that kind of user.
// Normally, we are interested in all languages spoken by Care Givers
func FetchAllLanguagesForKind(service *services.Service, kind string) ([]string, error) {
	log.Startedf(service.UserID, "FetchAllLanguagesForKind", "kind[%s]", kind)

	var result []string
	f := func(collection *mgo.Collection) error {
		queryMap := bson.M{"profile.type": kind}

		log.Trace(service.UserID, "FetchAllLanguagesForKind", "MGO : db.%s.distinct(\"profile.languages\",%s)", Config.Collection, mongo.ToString(queryMap))
		return collection.Find(queryMap).Distinct("profile.languages", &result)
	}

	if err := service.DBAction(Config.Database, Config.Collection, f); err != nil {
		if err != mgo.ErrNotFound {
			log.CompletedError(err, service.UserID, "FetchAllLanguagesForKind")
			return nil, err
		}
	}

	log.Completedf(service.UserID, "FetchAllLanguagesForKind", "languages for %s: %+v", kind, result)
	return result, nil
}

// FetchAllSkills retrieves a list of all skills available.
func FetchAllSkills(service *services.Service) ([]string, error) {
	log.Startedf(service.UserID, "FetchAllSkills", "")

	var result []string
	f := func(collection *mgo.Collection) error {
		queryMap := bson.M{"profile.type": "cg"}

		log.Trace(service.UserID, "FetchAllSkills", "MGO : db.%s.distinct(\"profile.skills\",%s)", Config.Collection, mongo.ToString(queryMap))
		return collection.Find(queryMap).Distinct("profile.skills", &result)
	}

	if err := service.DBAction(Config.Database, Config.Collection, f); err != nil {
		if err != mgo.ErrNotFound {
			log.CompletedError(err, service.UserID, "FetchAllSkills")
			return nil, err
		}
	}

	log.Completedf(service.UserID, "FetchAllSkills", "care giver skills: %+v", result)
	return result, nil
}

// FetchProfile retrieves the profile of the user specified by their email
func FetchProfile(service *services.Service, email string) (*userModel.Profile, error) {
	log.Startedf(service.UserID, "FetchProfile", "email[%s]", email)

	var user userModel.User
	f := func(collection *mgo.Collection) error {
		queryMap := bson.M{"email": email}

		log.Trace(service.UserID, "FetchProfile", "MGO : db.%s.find(%s,%s).limit(1)", Config.Collection, mongo.ToString(queryMap),"{\"profile\": 1, \"_id\": 0}")
		return collection.Find(queryMap).Select(bson.M{"profile": 1, "_id": 0}).One(&user)
	}

	if err := service.DBAction(Config.Database, Config.Collection, f); err != nil {
		if err != mgo.ErrNotFound {
			log.CompletedError(err, service.UserID, "FetchProfile")
			return nil, err
		}
	}

	log.Completedf(service.UserID, "FetchProfile", "user profile: %+v", user.Profile)
	return &user.Profile, nil
}

// FindUsersOfType retrieves the users for the specified type
func FindUsersOfKind(service *services.Service, kind string) ([]userModel.User, error) {
	log.Startedf(service.UserID, "FindUsersOfKind", "kind[%s]", kind)

	var users []userModel.User
	f := func(collection *mgo.Collection) error {
		queryMap := bson.M{"profile.type": kind}

		log.Trace(service.UserID, "FindUsersOfKind", "Query : db.%s.find(%s)", Config.Collection, mongo.ToString(queryMap))
		return collection.Find(queryMap).Select(bson.M{"email": 1, "profile": 1, "_id": 0}).All(&users)
	}

	if err := service.DBAction(Config.Database, Config.Collection, f); err != nil {
		log.CompletedError(err, service.UserID, "FindUsersOfKind")
		return nil, err
	}

	log.Completedf(service.UserID, "FindUsersOfKind", "users of type %s %+v", kind, users)
	return users, nil
}

// FindCareGiversForLanguageAndSkill retrieves all care givers that speak a certain language and have certain skill
// A zero argument runs a query ignoring that argument:
func FindCareGiversForLanguageAndSkill(service *services.Service, language string, skill string) ([]userModel.User, error) {
	log.Startedf(service.UserID, "FindCareGiversForLanguageAndSkill", "language[%s] and skill[%s]", language,skill)

	var users []userModel.User
	f := func(collection *mgo.Collection) error {

		queryMap := bson.M{"profile.type": "cg"}
    if language != "" {
      queryMap["profile.languages"] = language
    }
    if skill != "" {
      queryMap["profile.skills"] = skill
    }

		log.Trace(service.UserID, "FindCareGiversForLanguageAndSkill", "Query : db.%s.find(%s)", Config.Collection, mongo.ToString(queryMap))
		return collection.Find(queryMap).Select(bson.M{"email": 1, "profile": 1, "_id": 0}).All(&users)
	}

	if err := service.DBAction(Config.Database, Config.Collection, f); err != nil {
		log.CompletedError(err, service.UserID, "FindCareGiversForLanguageAndSkill")
		return nil, err
	}

	log.Completedf(service.UserID, "FindCareGiversForLanguageAndSkill", "care givers that speak %s with skill %s %+v", language, skill, users)
	return users, nil
}
