package controllers

import (
  "github.com/astaxie/beego"
  log "github.com/goinggo/tracelog"
)

type FileController struct {
	beego.Controller
}

// @router / [get]
func (o *FileController) Index() {
  log.Trace("", "Index", "File controller should just serve file back and set headers??")
	//o.Data["json"] = {}
	//o.ServeJson()
}
