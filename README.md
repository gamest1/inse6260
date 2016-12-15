# INSE 6260 Spectrix Health Care Scheduling System

This project was developed using the *Beego Mgo Example* project as a base using the following technologies:

* GoLang
* MongoDB
* GoConvey (Testing Framework)
* React + Socket.IO

-- Run the web service
```bash
cd $GOPATH/src/github.com/goinggo/beego-mgo/zscripts
./runlocal.sh
```

-- Test Web App
Run the home page!
	http://localhost:9003

-- Run the test cases
```bash
cd $GOPATH/src/github.com/goinggo/beego-mgo
./testconvey.sh
```


## Beego Mgo Example

Copyright 2013 Ardan Studios. All rights reserved.  
Use of this source code is governed by a BSD-style license that can be found in the LICENSE handle.

This application provides a sample to use the beego web framework and the Go MongoDB driver mgo. This program connects to a public MongoDB at MongoLab. A single collection is available for testing.

The project includes several shell scripts in the zscripts folder to make building, running and testing the web application easier.

GoingGo.net Post:  
http://www.goinggo.net/2013/12/sample-web-application-using-beego-and.html

Ardan Studios  
12973 SW 112 ST, Suite 153  
Miami, FL 33186  
bill@ardanstudios.com

### Installation

	-- YOU MUST HAVE BAZAAR INSTALLED
	http://wiki.bazaar.canonical.com/Download

	-- Get, build and install the code
	go get github.com/goinggo/beego-mgo

	-- Run the web service
	cd $GOPATH/src/github.com/goinggo/beego-mgo/zscripts
	./runbuild.sh

	-- Run the tests
	cd $GOPATH/src/github.com/goinggo/beego-mgo/zscripts
	./runtests.sh

	-- Test Web Service API's
	Run the home page and go through the tabs
	http://localhost:9003

### Notes About Architecture

I have been asked why I have organized the code in this way?

The models folder contains the data structures for the individual services. Each service places their models in a separate folder.

The services folder contain the raw service calls that the business layer would use to implement higher level functionality.

The controller methods handle and process the requests.

The more that can be abstracted into the base controller and base service the better. This way, adding a new functionality is simple and you don't need to worry about forgetting to do something important. Authentication always comes to mind.

The utilities folder is just that, support for the web application, mostly used by the services. You have exception handling support, extended logging support and the mongo support.

The abstraction layer for executing MongoDB queries and commands help hide the boilerplate code away into the base service and mongo utility code.

Using environmental variables for the configuration parameters provides a best practice for minimizing security risks. The scripts in the zscripts folder contains the environment variables required to run the web application. In a real project these settings would never be saved in source control.
