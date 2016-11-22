// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE handle.

// Package routes initializes the routes for the web service.
package routes

import (
	"github.com/astaxie/beego"
	"github.com/goinggo/beego-mgo/controllers"
)

func init() {
	//Unauthenticated (no session needed):
	beego.Router("/user/login", new(controllers.MainController))
	beego.Router("/user/logout", new(controllers.MainController), "get:Logout")
	beego.Router("/user/register", new(controllers.MainController), "get,post:Register")

	//Authenticated (no session needed):
	beego.Router("/user/home", new(controllers.MainController), "get:Home")
	beego.Router("/user/profile", new(controllers.MainController), "get,post:Profile") //get: displays availability; post: updates availability
	beego.Router("/user/request", new(controllers.MainController), "get,post:Request") //get: displays a new form; post: creates new service request

	beego.Router("/user/day", new(controllers.MainController), "get:DisplayDay") //get: displays all daily activity

	//AJAX
	beego.Router("/requests/:userId", new(controllers.RequestController), "get:RequestsForUser")
	beego.Router("/user/display/:userId", new(controllers.MainController), "get:DisplayAll") //get: displays all users in the system if userId is an agent

	//Old sample routes:
	beego.Router("/users/retrieveuser", new(controllers.UserController), "post:RetrieveUser")
	beego.Router("/users", new(controllers.UserController), "get:Index")
	beego.Router("/users/retrieveuser", new(controllers.UserController), "post:RetrieveUser")

	beego.Router("/requests", new(controllers.RequestController), "get:Index")
	beego.Router("/requests/createnew", new(controllers.RequestController), "post:CreateRequest")

	beego.Router("/react", new(controllers.ReactController), "get:Index")
}
