package controllers

import (
  "time"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"

	"github.com/goinggo/beego-mgo/services"
  "github.com/goinggo/beego-mgo/services/userService"
	"github.com/goinggo/beego-mgo/utilities/mongo"

	log "github.com/goinggo/tracelog"
)

type MainController struct {
	beego.Controller
  services.Service
}

func (this *MainController) activeContent(view string, scripts bool) {
  log.Trace("", "activeContent", "Loading active content for %s", view)
	this.Layout = "shared/basic-layout.html"
	this.LayoutSections = make(map[string]string)
  if scripts {
  	this.LayoutSections["Scripts"] = "scripts/" + view + ".html"
  }
	this.LayoutSections["Header"] = "shared/header.html"
	this.LayoutSections["Footer"] = "shared/footer.html"
	this.LayoutSections["Modal"] = "shared/modal.html"
	this.TplName = view + ".html"

	sess := this.GetSession("acme")
	if sess != nil {
		this.Data["InSession"] = 1
		m := sess.(map[string]interface{})
		this.Data["First"] = m["first"]
		this.Data["Email"] = m["username"]
		this.Data["Type"] = m["type"]
    if m["type"] == "p" {
      this.Data["TypeText"] = "Patient"
    } else if m["type"] == "a" {
      this.Data["TypeText"] = "Agent"
    } else if m["type"] == "cg" {
      this.Data["TypeText"] = "Care Giver"
    } else {
      log.Trace("", "Unknown type found in session", "Don't know what to do...")
  		return
    }
	} else {
    this.Data["InSession"] = 0
  }
}

func (this *MainController) Prepare() {
  log.Startedf("", "Preparing MainController for DB access", "")
	if err := this.Service.Prepare(); err != nil {
		log.Errorf(err, "", "MainController.Prepare", this.Ctx.Request.URL.Path)
		return
	}

	log.Trace("", "MainController.Prepare finished", "Path[%s]", this.Ctx.Request.URL.Path)
}

func (this *MainController) Finish() {
	defer func() {
		if this.MongoSession != nil {
			log.Trace("", "Finish", "Closing Session from sharedController")
			mongo.CloseSession("", this.MongoSession)
			this.MongoSession = nil
		}
	}()

	log.Completedf("", "Finish", this.Ctx.Request.URL.Path)
}


func (this *MainController) Get() {
  log.Startedf("", "Get MainController", "")
	this.activeContent("index", false)

	//******** Gmail style: If the previous user didn't logout, redirect back to the home page:
	sess := this.GetSession("acme")
	if sess != nil {
	  m := sess.(map[string]interface{})
    log.Trace(m["username"].(string), "Get MainController", "Session existed. Redirecting user to home at: %s", (m["timestamp"].(time.Time)).Local())
		this.Redirect("/user/home", 302)
		return
	}

  log.Trace("", "Get MainController", "Rendering login page")
}

func (this *MainController) Post() {
    log.Startedf("", "Post MainController", "")
    if this.Ctx.Input.Method() == "POST" {
        this.activeContent("index", false)

    		email := this.GetString("email")
    		password := this.GetString("password")
        log.Trace("POST", "Validating credentials: [%s],[%s]", email, password)
    		valid := validation.Validation{}
    		valid.Email(email, "email")
    		valid.Required(password, "password")
    		if valid.HasErrors() {
    			errormap := []string{}
    			for _, err := range valid.Errors {
    				errormap = append(errormap, "Validation failed on "+err.Key+": "+err.Message+"\n")
    			}
    			this.Data["Errors"] = errormap
          return
    		}
        log.Trace("POST", "Login credentials respect format [%s],[%s]", email, password)

        userProfile, err := userService.Login(&this.Service, email, password)
        if err != nil {
          errormap := []string{err.Error()}
          log.CompletedErrorf(err, email, "MainController.Post", "Login failed with error map: %+v", errormap)
          this.Data["Errors"] = errormap
          return
        }

    		//******** Create session and go back to previous page
    		m := make(map[string]interface{})
    		m["first"] = userProfile.FirstName
    		m["username"] = email
    		m["timestamp"] = time.Now()
        m["type"] = userProfile.Type
    		this.SetSession("acme", m)
    		this.Redirect("/user/home", 302)
    	} else {
        log.Trace("POST", ".Ctx.Input.Method [%s]", this.Ctx.Input.Method())
      }
}

func (this *MainController) Logout() {
	this.DelSession("acme")
  this.Redirect("/user/login", 302)
}

// ServeError prepares and serves an Error exception.
func (this *MainController) ServeError(err error) {
	this.Data["json"] = struct {
		Error string `json:"Error"`
	}{err.Error()}
	this.Ctx.Output.SetStatus(500)
	this.ServeJSON()
}

func (this *MainController) Notice() {
	this.activeContent("notice", false)

	flash := beego.ReadFromRequest(&this.Controller)
	if n, ok := flash.Data["notice"]; ok {
		this.Data["notice"] = n
	}
}
