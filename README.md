# INSE 6260 Spectrix Health Care Scheduling System

This project was developed using the *Beego Mgo Example* project as a base using the following technologies:

* GoLang
* MongoDB
* GoConvey (Testing Framework)
* React + Socket.IO

To run Spectrix on your local machine, you must install Bazaar, mongo, and golang first. Make sure you have a working GoLang system on your computer: $GOPATH must be set and make sure you have reading access to the project's code repository located [https://github.com/gamest1/inse6260.git](https://github.com/gamest1/inse6260.git), if you don't have access to the repository, unzip the project's source file inside the $GOPATH/src/ folder and fix all dependancies problem doing *go get* everytime a package dependancy is not found.

Then:

```bash
go get github.com/gamest1/inse6260/...
```

(The three dots at the end guarantee that you will get all project dependancies) For more instructions about setting up your GoLang working environment, visit: [https://golang.org/doc/install](https://golang.org/doc/install)

Create a directory:

```bash
mkdir -p /data/
```

where all the database data will be stored.

Now, you must run two mongo deamons in a primary secondary replica set configuration, and load some data onto a database called inse6260. To setup the primary/secondary configuration, you may use the provided .conf files in the /dbup folder using:

```bash
mongod -f [your path]/dbup/mongod1.conf
mongod -f [your path]/dbup/mongod2.conf
```

Additionally, use:

```bash
 mongorestore [your path]/dbup/
```

To load pre-existing data to your database  (you may need to delete the .conf files from [your path]/dbup/ for the previous command to work).

Finally, a user with readWrite permissions must be created to access the inse6260 database using the following credentials: username: guest, password: welcome

Open a mongo shell using:

```bash
 mongo
```

and use the following command for this purpose:

```bash
db.createUser( { user: "guest", pwd: "welcome", roles: [ { role: "readWrite", db: "inse6260" } ] } )
```

Once your database system is up and running:


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
