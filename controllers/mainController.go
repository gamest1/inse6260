package controllers

import (
  "regexp"
  "reflect"

	log "github.com/goinggo/tracelog"
	"github.com/astaxie/beego/validation"

	"github.com/goinggo/beego-mgo/models/userModel"
  "github.com/goinggo/beego-mgo/services/userService"
  "github.com/goinggo/beego-mgo/utilities/location"
  "github.com/goinggo/beego-mgo/utilities/availability"
)

func (this *MainController) Register() {
	this.activeContent("user/register", true)
  systemSkills, err := userService.FetchAllSkills(&this.Service)
  if err != nil {
    log.CompletedErrorf(err, "", "MainController.Register", "FetchAllSkills")
    this.ServeError(err)
    return
  }

  systemLanguages, err := userService.FetchAllLanguagesForKind(&this.Service, "cg")
  if err != nil {
    log.CompletedErrorf(err, "", "MainController.Register", "FetchAllLanguagesForKind[%s]", "cg")
    this.ServeError(err)
    return
  }

  this.Data["Skills"] = systemSkills
  this.Data["Languages"] = systemLanguages
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
    valid.MinSize(skills, 1, "Skill selection")
    valid.MinSize(kind, 1, "User Type")

    if kind == "a" {
      token := this.GetString("token")
      if (token != "secret token") {
        errormap := []string{"Secret token doesn't match"}
        this.Data["Errors"] = errormap
        log.Trace("", "POST Registration", "3) Token: %+v", errormap)
        return
      }
    }

    newAvailability := &availability.Availability{0, 0, 0, 0, 0, 0, 0}
    if kind == "cg" {
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

    log.Trace("POST Registration", "Should redirect or create session or something", "")

	} else {
    log.Trace("", "Register GET", "I'm here!!!")
  }
}

func (this *MainController) Profile() {
	this.activeContent("user/profile", false)

  if this.Ctx.Input.Method() == "POST" {
    log.Trace("", "Profile POST", "I'm here!!!")
	} else {
    log.Trace("", "Profile GET", "I'm here!!!")
  }
}

func (this *MainController) Request() {
	this.activeContent("user/request", false)

  if this.Ctx.Input.Method() == "POST" {
    log.Trace("", "Request POST", "I'm here!!!")
	} else {
    log.Trace("", "Request GET", "I'm here!!!")
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

  this.activeContent("user/home", false)
}

func (this *MainController) DisplayAll() {
  log.Trace("", "Home DisplayAll???", "I'm here!!!")
  this.activeContent("user/displayall", false)
}

func (this *MainController) DisplayDay() {
  log.Trace("", "Home DisplayDay???", "I'm here!!!")
  this.activeContent("user/displayd", false)
}
