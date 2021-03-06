// Copyright 2013 Ardan Studios. All rights reserved.
// Use of baseController source code is governed by a BSD-style
// license that can be found in the LICENSE handle.

// Package baseController implements boilerplate code for all baseControllers.
package baseController

import (
	"reflect"
	"runtime"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/goinggo/beego-mgo/localize"
	"github.com/goinggo/beego-mgo/services"
	"github.com/goinggo/beego-mgo/utilities/mongo"
	log "github.com/goinggo/tracelog"
)

//** TYPES

type (
	// BaseController composes all required types and behavior.
	BaseController struct {
		beego.Controller
		services.Service
	}
)

//** INTERCEPT FUNCTIONS

// Prepare is called prior to the baseController method.
func (baseController *BaseController) Prepare() {
	baseController.UserID = baseController.GetString("userID")
	if baseController.UserID == "" {
		baseController.UserID = baseController.GetString(":userID")
	}
	if baseController.UserID == "" {
		baseController.UserID = "Unknown"
	}

	if err := baseController.Service.Prepare(); err != nil {
		log.Errorf(err, baseController.UserID, "BaseController.Prepare", baseController.Ctx.Request.URL.Path)
		baseController.ServeError(err)
		return
	}

	log.Trace(baseController.UserID, "BaseController.Prepare", "UserID[%s] Path[%s]", baseController.UserID, baseController.Ctx.Request.URL.Path)
}

// Finish is called once the baseController method completes.
func (baseController *BaseController) Finish() {
	defer func() {
		if baseController.MongoSession != nil {
			log.Trace(baseController.UserID, "Finish", "Closing Session from baseController")
			mongo.CloseSession(baseController.UserID, baseController.MongoSession)
			baseController.MongoSession = nil
		}
	}()

	log.Completedf(baseController.UserID, "Finish", baseController.Ctx.Request.URL.Path)
}

//** VALIDATION

// ParseAndValidate will run the params through the validation framework and then
// response with the specified localized or provided message.
func (baseController *BaseController) ParseAndValidate(params interface{}) bool {
	// This is not working anymore :(
	log.Trace(baseController.UserID, "ParseAndValidate", "About to call ParseForm")
	if err := baseController.ParseForm(params); err != nil {
		log.Trace(baseController.UserID, "ParseAndValidate", "ParseForm returned error: %s", err)
		baseController.ServeError(err)
		return false
	}
	//ParseForm fails to fetch arrays into the params variable. We need to do that manually:
	allPostData := baseController.Ctx.Request.PostForm
	log.Trace(baseController.UserID, "ParseAndValidate", "Lookup table: %+v", allPostData)

  st := reflect.ValueOf(params).Elem()
  typeOfST := st.Type()
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		switch field.Kind() {
		case reflect.Slice:
			lookup := typeOfST.Field(i).Tag.Get("form")
			language := allPostData[lookup]
			var array []string
			if language == nil {
					lookup = typeOfST.Field(i).Tag.Get("form") + "[]" //Beego seems to attach [] to slices!
					array = allPostData[lookup]
			} else {
				 array = append(array,language[0])
			}
			log.Trace(baseController.UserID, "ParseAndValidate", "Modifying slice internals with lookup: %s using %+v of type: %v", lookup, array, reflect.TypeOf(array))
			//Make the new slice:
			slice := reflect.MakeSlice(reflect.TypeOf(array), len(array), len(array))
			for idx, element := range array {
					v := slice.Index(idx)
					v.Set(reflect.ValueOf(element))
			}
			field.Set(slice)
			log.Trace(baseController.UserID, "ParseAndValidate", "Slice set!!")
		}
	}

	var valid validation.Validation
	log.Trace(baseController.UserID, "ParseAndValidate", "About to call Valid with completed params: %+v",params)
	ok, err := valid.Valid(params)
	if err != nil {
		log.Trace(baseController.UserID, "ParseAndValidate", "Valid returned error")
		baseController.ServeError(err)
		return false
	}

	if ok == false {
		log.Trace(baseController.UserID, "ParseAndValidate", "but ok is false :(")
		// Build a map of the Error messages for each field
		messages2 := make(map[string]string)

		val := reflect.ValueOf(params).Elem()
		for i := 0; i < val.NumField(); i++ {
			// Look for an Error tag in the field
			typeField := val.Type().Field(i)
			tag := typeField.Tag
			tagValue := tag.Get("Error")

			// Was there an Error tag
			if tagValue != "" {
				messages2[typeField.Name] = tagValue
			}
		}

		// Build the Error response
		log.Trace(baseController.UserID, "ParseAndValidate", "Building error response")
		var errors []string
		for _, err := range valid.Errors {
			// Match an Error from the validation framework Errors
			// to a field name we have a mapping for
			message, ok := messages2[err.Field]
			if ok == true {
				// Use a localized message if one exists
				errors = append(errors, localize.T(message))
				continue
			}

			// No match, so use the message as is
			errors = append(errors, err.Message)
		}

		log.Trace(baseController.UserID, "ParseAndValidate", "About to return errors: %+v", errors)
		baseController.ServeValidationErrors(errors)
		return false
	}

	return true
}

//** EXCEPTIONS

// ServeError prepares and serves an Error exception.
func (baseController *BaseController) ServeError(err error) {
	baseController.Data["json"] = struct {
		Error string `json:"Error"`
	}{err.Error()}
	baseController.Ctx.Output.SetStatus(500)
	baseController.ServeJSON()
}

// ServeValidationErrors prepares and serves a validation exception.
func (baseController *BaseController) ServeValidationErrors(Errors []string) {
	baseController.Data["json"] = struct {
		Errors []string `json:"Errors"`
	}{Errors}
	baseController.Ctx.Output.SetStatus(409)
	baseController.ServeJSON()
}

//** CATCHING PANICS

// CatchPanic is used to catch any Panic and log exceptions. Returns a 500 as the response.
func (baseController *BaseController) CatchPanic(functionName string) {
	if r := recover(); r != nil {
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		log.Warning(baseController.Service.UserID, functionName, "PANIC Defered [%v] : Stack Trace : %v", r, string(buf))

		baseController.ServeError(fmt.Errorf("%v", r))
	}
}

//** AJAX SUPPORT

// AjaxResponse returns a standard ajax response.
func (baseController *BaseController) AjaxResponse(resultCode int, resultString string, data interface{}) {
	response := struct {
		Result       int
		ResultString string
		ResultObject interface{}
	}{
		Result:       resultCode,
		ResultString: resultString,
		ResultObject: data,
	}

	baseController.Data["json"] = response
	baseController.ServeJSON()
}
