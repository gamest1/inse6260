package controllers

import (
  "regexp"
  "reflect"
  "time"
  "strings"

	log "github.com/goinggo/tracelog"
	"github.com/astaxie/beego/validation"

  "github.com/goinggo/beego-mgo/models/requestModel"
  "github.com/goinggo/beego-mgo/services/requestService"
	"github.com/goinggo/beego-mgo/models/userModel"
  "github.com/goinggo/beego-mgo/services/userService"
  "github.com/goinggo/beego-mgo/utilities/location"
  "github.com/goinggo/beego-mgo/utilities/availability"
)

type RequestsResponse struct {
    UserType string
    Requests []requestModel.Request
}

func (this *MainController) addSystemLanguagesAndSkills() {
  systemSkills, err := userService.FetchAllSkills(&this.Service)
  if err != nil {
    log.CompletedErrorf(err, "", "MainController.Register", "FetchAllSkills")
  }

  systemLanguages, err := userService.FetchAllLanguagesForKind(&this.Service, "cg")
  if err != nil {
    log.CompletedErrorf(err, "", "MainController.Register", "FetchAllLanguagesForKind[%s]", "cg")
  }

  this.Data["Skills"] = systemSkills
  this.Data["Languages"] = systemLanguages
}

func (this *MainController) Register() {
	this.activeContent("user/register", true)
  this.addSystemLanguagesAndSkills();
  this.Data["Days"] = []string{"monday","tuesday","wednesday","thursday","friday","saturday","sunday"}
	if this.Ctx.Input.Method() == "POST" {
    //This could be put into a struct, perhaps?
    firstName := this.GetString("first")
    lastName := this.GetString("last")
    gender := this.GetString("gender")
    email := this.GetString("email")
    password := this.GetString("password")
    password2 := this.GetString("password2")
    apartment := this.GetString("apartment")
    streetNumber, err := this.GetInt("streetnumber")
    if err != nil {
      log.Trace("", "POST Registration", "Cannot parse number: %v", err)
    }

    streetName := this.GetString("streetname")
    cityName := this.GetString("cityname")
    postalCode := this.GetString("postalcode")
    kind := this.GetString("type")

    languages := this.GetStrings("languages")
    skills := this.GetStrings("skills")

    if len(languages) == 1 {
      if languages[0] == "other" {
        languages = []string{this.GetString("otherlanguage")}
      }
    }
    if len(skills) == 1 {
      if skills[0] == "other" {
        skills = []string{this.GetString("otherskill")}
      }
    }

    log.Trace("", "POST Registration", "Validating input: [%s],[%s],[%s], %d, %+v, %+v", email, firstName, gender, streetNumber, languages, skills)

    valid := validation.Validation{}

    valid.MinSize(firstName, 1, "First Name")
    valid.MinSize(lastName, 1, "Last Name")
    valid.MinSize(gender, 1, "Gender")
    valid.Email(email, "Email")
    valid.MinSize(password, 6, "Password")
    valid.MinSize(apartment, 1, "Apartment")

    nType := reflect.TypeOf(streetNumber)
    if k := nType.Kind(); k != reflect.Int {
        valid.Numeric(56, "Street Number") //This will fillup an error.
	  }

    valid.MinSize(streetName, 1, "Street Name")
    valid.MinSize(cityName, 1, "City Name")
    valid.Match(postalCode, regexp.MustCompile(`(?i)([ABCEGHJKLMNPRSTVXY]\d)([ABCEGHJKLMNPRSTVWXYZ]\d){2}`), "Postal Code")
    valid.MinSize(languages, 1, "Language selection")
    valid.MinSize(kind, 1, "User Type")

    if kind == "a" {
      token := this.GetString("token")
      if (token != "secrettoken") {
        errormap := []string{"Secret token doesn't match"}
        this.Data["Errors"] = errormap
        log.Trace("", "POST Registration", "3) Token: %+v", errormap)
        return
      }
    }

    newAvailability := &availability.Availability{0, 0, 0, 0, 0, 0, 0}
    if kind == "cg" {
      valid.MinSize(skills, 1, "Skill selection")
      for i := 0; i < 7 ; i++ {
          switch i {
            case 0:
              monValue, _ := this.GetInt("monday")
              newAvailability.Monday = monValue
              break;
            case 1:
              tueValue, _ := this.GetInt("tuesday")
              newAvailability.Tuesday = tueValue
              break;
            case 2:
              wedValue, _ := this.GetInt("wednesday")
              newAvailability.Wednesday = wedValue
              break;
            case 3:
              thuValue, _ := this.GetInt("thursday")
              newAvailability.Thursday = thuValue
              break;
            case 4:
              friValue, _ := this.GetInt("friday")
              newAvailability.Friday = friValue
              break;
            case 5:
              satValue, _ := this.GetInt("saturday")
              newAvailability.Saturday = satValue
              break;
            case 6:
              sunValue, _ := this.GetInt("sunday")
              newAvailability.Sunday = sunValue
              break;
            default:
              break;
          }
	    }
    }

    if valid.HasErrors() {
      errormap := []string{}
      for _, err := range valid.Errors {
        errormap = append(errormap, "Validation failed on " + err.Key + ": " + err.Message + "\n")
      }
      this.Data["Errors"] = errormap
      log.Trace("", "POST Registration", "1) valid.HasErrors[%+v]", errormap)
      return
    } else if password != password2 {
      errormap := []string{"Provided passwords don't match"}
      this.Data["Errors"] = errormap
      log.Trace("", "POST Registration", "2) Passwords: %+v", errormap)
      return
    }

    //Create a new  map[string]interface{} (bson.M) Object using the given params:

    newLocation := &location.Location{0, 0, apartment, streetNumber, streetName, cityName, "", postalCode}
    newProfile := &userModel.Profile{firstName, lastName, gender, languages, *newLocation, kind, skills, *newAvailability}
    newUser := map[string]interface{}{
        "email":    email,
        "password": password,
        "profile":  newProfile,
    }

    log.Trace("", "POST Registration", "Creating newUser object %+v", newUser)
    err = userService.InsertNewUser(&this.Service, newUser)
    if err != nil {
      errormap := []string{err.Error()}
      log.CompletedErrorf(err, email, "MainController.Post", "Cannot create new user: %+v", errormap)
      this.Data["Errors"] = errormap
      return
    }

    this.Data["Success"] = 1
    log.Trace("POST Registration", "Done! Redirecting to home:", "")
	} else {
    log.Trace("", "Register GET", "I'm here!!!")
  }
}

func (this *MainController) Profile() {
  //******** This page requires login
  sess := this.GetSession("acme")
  if sess == nil {
    this.Redirect("/user/login", 302)
    return
  }
	m := sess.(map[string]interface{})
	this.activeContent("user/profile", true)

  type Day struct {
    Key string
    Value int
  }

  var days = []Day {
    Day{"monday",m["availability"].(availability.Availability).Monday},
    Day{"tuesday",m["availability"].(availability.Availability).Tuesday},
    Day{"wednesday",m["availability"].(availability.Availability).Wednesday},
    Day{"thursday",m["availability"].(availability.Availability).Thursday},
    Day{"friday",m["availability"].(availability.Availability).Friday},
    Day{"saturday",m["availability"].(availability.Availability).Saturday},
    Day{"sunday",m["availability"].(availability.Availability).Sunday},
  }

  this.Data["Success"] = 0
  this.Data["Days"] = days
  this.Data["First"] =  m["first"].(string)
  this.Data["Last"] = m["last"].(string)
  this.Data["Email"] = m["username"].(string)
  log.Trace("", "Profile", "Current availability: %+v", days)

  if this.Ctx.Input.Method() == "POST" {
    //This could be put into a struct, perhaps?
    log.Trace("", "POST Profile", "")
    errormap := []string{}

    email := m["username"]
    currentPassword := this.GetString("currentpassword")
    newPassword := this.GetString("password")
    newPassword2 := this.GetString("password2")

      newAvailability := &availability.Availability{0, 0, 0, 0, 0, 0, 0}
      for i := 0; i < 7 ; i++ {
          switch i {
            case 0:
              monValue, _ := this.GetInt("monday")
              newAvailability.Monday = monValue
              break;
            case 1:
              tueValue, _ := this.GetInt("tuesday")
              newAvailability.Tuesday = tueValue
              break;
            case 2:
              wedValue, _ := this.GetInt("wednesday")
              newAvailability.Wednesday = wedValue
              break;
            case 3:
              thuValue, _ := this.GetInt("thursday")
              newAvailability.Thursday = thuValue
              break;
            case 4:
              friValue, _ := this.GetInt("friday")
              newAvailability.Friday = friValue
              break;
            case 5:
              satValue, _ := this.GetInt("saturday")
              newAvailability.Saturday = satValue
              break;
            case 6:
              sunValue, _ := this.GetInt("sunday")
              newAvailability.Sunday = sunValue
              break;
            default:
              break;
          }
      }

    log.Trace("", "POST Request", "Validating input: [%s],[%s],[%s], %+v", currentPassword, newPassword, newPassword2, newAvailability)

    updatesPassword := false
    if currentPassword != "" && newPassword != "" && newPassword2  != "" {
      //Person wanted to update password
      valid := validation.Validation{}
      valid.MinSize(newPassword, 6, "Password")
      if valid.HasErrors() {
        for _, err := range valid.Errors {
          errormap = append(errormap, "Validation failed on " + err.Key + ": " + err.Message + "\n")
        }
        this.Data["Errors"] = errormap
        log.Trace("", "POST Profile", "valid.HasErrors[%+v]", errormap)
        return
      } else if newPassword != newPassword2 {
        errormap := []string{"Provided passwords don't match"}
        this.Data["Errors"] = errormap
        log.Trace("", "POST Profile", "2) Passwords: %+v", errormap)
        return
      }
      //Update password:
      updatesPassword = true
    }

    //Update availability
    newProfile := &userModel.Profile{ Availability : *newAvailability }
    newUser := map[string]interface{}{
      "email" : email,
      "profile" : newProfile,
    }

    if updatesPassword {
      newUser["currentpassword"] = currentPassword
      newUser["newpassword"] = newPassword
    }

    log.Trace("", "POST Profile", "Updating %s with newUser object %+v", email, newUser)
    err := userService.UpdateUserAvailability(&this.Service, newUser)
    if err != nil {
      errormap := []string{err.Error()}
      log.CompletedErrorf(err, email.(string), "MainController.Post", "Cannot update user: %+v", errormap)
      this.Data["Errors"] = errormap
      return
    }

    this.Data["Success"] = 1
    m["availability"] = *newAvailability
    log.Trace("POST Profile", "Done!", "")
	} else {
    log.Trace("", "GET Profile", "")
  }
}

func (this *MainController) Request() {
  //******** This page requires login
  sess := this.GetSession("acme")
  if sess == nil {
    this.Redirect("/user/login", 302)
    return
  }
	m := sess.(map[string]interface{})
	this.activeContent("user/request", true)
  this.addSystemLanguagesAndSkills();
  this.Data["Success"] = 0

  if this.Ctx.Input.Method() == "POST" {
    //This could be put into a struct, perhaps?
    log.Trace("", "POST Request", "")
    errormap := []string{}

    email := m["username"]
    apartment := m["location"].(location.Location).Apartment
    streetNumber := m["location"].(location.Location).Number
    streetName := m["location"].(location.Location).Street
    cityName := m["location"].(location.Location).City
    stateName := m["location"].(location.Location).State
    postalCode := m["location"].(location.Location).Zip

    loc := this.GetString("location")
    if loc != "" {
      apartment = this.GetString("apartment")
      streetNumber, _ = this.GetInt("streetnumber")
      streetName = this.GetString("streetname")
      cityName = this.GetString("cityname")
      stateName = ""
      postalCode = this.GetString("postalcode")
    }

    serviceDate := this.GetString("servicedate")
    serviceTime := this.GetString("servicetime")
    dateTimeString := serviceDate + "T" + serviceTime
    if( len(strings.Split(dateTimeString,":")) < 3 ) {
      dateTimeString = dateTimeString + ":00"
    }

    startTime, err := time.Parse("2006-01-02T15:04:05",dateTimeString)
    if err != nil {
      log.Trace("", "POST Request", "Error parsing time: %v", err)
      errormap = append(errormap, "Validation failed on Service Date or Service Time: " + err.Error() + "\n")
    }

    duration, err := this.GetInt("duration")
    if err != nil {
      log.Trace("", "POST Request", "2) Cannot parse number: %v", err)
      errormap = append(errormap, "Validation failed on duration: " + err.Error() + "\n")
    }

    skill := this.GetString("skill")
    gender := this.GetString("gender")
    languages := this.GetStrings("languages")
    if len(languages) == 1 {
      if languages[0] == "other" {
        languages = []string{this.GetString("otherlanguage")}
      }
    }

    log.Trace("", "POST Request", "Validating input: [%s],[%s],[%s], %d, %+v", skill, startTime, gender, duration, languages)

    valid := validation.Validation{}

    valid.MinSize(gender, 1, "Gender")
    valid.Email(email, "Email")
    valid.MinSize(apartment, 1, "Apartment")

    nType := reflect.TypeOf(streetNumber)
    if k := nType.Kind(); k != reflect.Int {
        valid.Numeric(56, "Street Number") //This will fillup an error.
	  }

    valid.MinSize(streetName, 1, "Street Name")
    valid.MinSize(cityName, 1, "City Name")
    valid.Match(postalCode, regexp.MustCompile(`(?i)([ABCEGHJKLMNPRSTVXY]\d)([ABCEGHJKLMNPRSTVWXYZ]\d){2}`), "Postal Code")
    valid.MinSize(languages, 1, "Language selection")

    if valid.HasErrors() {
      for _, err := range valid.Errors {
        errormap = append(errormap, "Validation failed on " + err.Key + ": " + err.Message + "\n")
      }
      this.Data["Errors"] = errormap
      log.Trace("", "POST Request", "valid.HasErrors[%+v]", errormap)
      return
    }

    //Create a new  map[string]interface{} (bson.M) Object using the given params:
    newLocation := &location.Location{0, 0, apartment, streetNumber, streetName, cityName, stateName, postalCode}
    newReqs := &requestModel.Requirements{skill, gender, languages, *newLocation}
    newRequest := map[string]interface{}{
        "time": startTime,
        "duration":   duration,
        "status": "pending",
        "request": newReqs,
        "originator": email,
    }

    log.Trace("", "POST Request", "Creating newRequest object %+v", newRequest)
    //And insert it!
  	err = requestService.InsertNewRequest(&this.Service, newRequest)
  	if err != nil {
      errormap := []string{err.Error()}
      log.Trace("", "POST Request", "Cannot create new request: %+v", errormap)
      this.Data["Errors"] = errormap
  		return
  	}

    this.Data["Success"] = 1
    log.Trace("POST Request", "Done! Redirecting to home:", "")
	} else {
    log.Trace("", "Request GET", "")
  }
}

func (this *MainController) Home() {
  log.Trace("", "Home GET", "")
  //******** This page requires login
  sess := this.GetSession("acme")
  if sess == nil {
    this.Redirect("/user/login", 302)
    return
  }

  this.activeContent("user/home", true)
}

func (this *MainController) DisplayDay() {
  log.Trace("", "Home DisplayDay???", "I'm here!!!")
  this.activeContent("user/displayd", false)
}

//** AJAX FUNCTIONS
func (this *MainController) DisplayAll() {
  log.Startedf("MainController", "DisplayAll", "")

  errormap := []string{}
  email := this.GetString(":userId")

  valid := validation.Validation{}
  valid.Email(email, "Email")

  if valid.HasErrors() {
    for _, err := range valid.Errors {
      errormap = append(errormap, "Validation failed on " + err.Key + ": " + err.Message + "\n")
    }
    log.Trace("", "DisplayAll", "valid.HasErrors[%+v]", errormap)
    return
  }

  log.Trace("MainController", "DisplayAll", "Investigating type of user: %s",email)
  userType, err := userService.TypeForUser(&this.Service, email)
	if err != nil {
		log.CompletedErrorf(err, "MainController", "DisplayAll", "TypeForUser[%s]", email)
		return
	}

  if userType == "a" {
    log.Trace("MainController", "DisplayAll", "Agent performing system fetch: %s - %s",email,userType)
    systemUsers, err := userService.FindUsersOfKind(&this.Service, "all")
    if err != nil {
      log.CompletedErrorf(err, "MainController", "DisplayAll", "FindUsersOfKind[%s]", "all")
      return
    }
    this.Data["json"] = systemUsers
    this.ServeJSON()
  } else {
    log.Trace("MainController", "DisplayAll", "Unauthorized user attempting system fetch: %s - %s",email,userType)
    return
  }
}

func (this *MainController) RequestsForUser() {
  log.Startedf("MainController", "RequestsForUser", "")

  errormap := []string{}
  email := this.GetString(":userId")

  valid := validation.Validation{}
  valid.Email(email, "Email")

  if valid.HasErrors() {
    for _, err := range valid.Errors {
      errormap = append(errormap, "Validation failed on " + err.Key + ": " + err.Message + "\n")
    }
    log.Trace("", "RequestsForUser", "valid.HasErrors[%+v]", errormap)
    return
  }

  log.Trace("MainController", "RequestsForUser", "Investigating type of user: %s",email)
  userType, err := userService.TypeForUser(&this.Service, email)
	if err != nil {
		log.CompletedErrorf(err, "MainController", "RequestsForUser", "TypeForUser[%s]", email)
		return
	}

  log.Trace("MainController", "RequestsForUser", "Fetching requests for user: %s - %s",email,userType)
  patientRequests, err := requestService.FetchAllRequestsForUser(&this.Service, email, userType)
	if err != nil {
		log.CompletedErrorf(err, "MainController", "RequestsForUser", "FetchAllRequestsForUser[%s]", email)
		return
	}

	this.Data["json"] = &RequestsResponse{ userType, patientRequests }
	this.ServeJSON()
}

func (this *MainController) CancelRequest() {
  this.updateRequest("canceled")
}
func (this *MainController) CompleteRequest() {
  this.updateRequest("completed")
}

func (this *MainController) updateRequest(status string) {
  log.Startedf("MainController", "UpdateRequest[%s]", status)

  errormap := []string{}
  reqId := this.GetString(":reqId")

  valid := validation.Validation{}
  valid.MinSize(reqId, 20, "Request ID")

  if valid.HasErrors() {
    for _, err := range valid.Errors {
      errormap = append(errormap, "Validation failed on " + err.Key + ": " + err.Message + "\n")
    }
    log.Trace("", "UpdateRequest", "valid.HasErrors[%+v]", errormap)
    return
  }

  log.Trace("MainController", "UpdateRequest", "Updating request [%s]", reqId)
  err := requestService.UpdateRequest(&this.Service, reqId, status)
	if err != nil {
		log.CompletedErrorf(err, "MainController", "UpdateRequest", "UpdateRequest[%s]: %s", reqId, status)
		return
	}

//	this.Data["json"] = &interface{}{}
	this.ServeJSON()
}
