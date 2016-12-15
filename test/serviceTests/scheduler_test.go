// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE handle.

// Package serviceTests implements tests for all system services.
package serviceTests

import (
	"time"
	"testing"

  . "github.com/goinggo/beego-mgo/utilities/scheduler"
	. "github.com/smartystreets/goconvey/convey"

	"gopkg.in/mgo.v2/bson"
)

// Test_CreateSchedule checks that a new schedule can be created via the scheduler
/*func Test_CreateSchedule(t *testing.T) {

	scheduler.CreateSchedule()

	Convey("Subject: Test CreateSchedule", t, func() {
		Convey("Should Be Able To Start the scheduler", func() {
			So(nil, ShouldEqual, nil)
		})
	})
}*/

func Test_IMRequests(t *testing.T) {
	myFormat := "Jan _2 2006 15:04:05"
	t1, _ := time.Parse(myFormat, "Dec 12 2016 06:00:00")
	app1 := Appointment{t1,1} // 6-7 appointment on Dec 12
	t2, _ := time.Parse(myFormat, "Dec 12 2016 08:00:00")
	app2 := Appointment{t2,2} // 8-10 appointment on Dec 12
	t3, _ := time.Parse(myFormat, "Dec 12 2016 09:00:00")
	app3 := Appointment{t3,2} // 9-11 appointment on Dec 12

	req1 := IMRequest{bson.NewObjectId(), nil, "pending", "", app1}
	req2 := IMRequest{bson.NewObjectId(), nil, "pending", "", app2}
	req3 := IMRequest{bson.NewObjectId(), nil, "pending", "", app3}
	req4 := IMRequest{bson.NewObjectId(), nil, "allocated", "cg1@test.com", app1}

	aSolution := []IMRequest{req1,req2,req4}

	SS := &SchedulingSolution{aSolution}
	BB := &SchedulingSolution{}

	m := make(map[string]IMBookings)
	m["cg1@test.com"] = IMBookings{make([]Appointment, 0)}
	m["cg2@test.com"] = IMBookings{make([]Appointment, 0)}
	cSolution := []IMRequest{req1,req2,req3}
	TT := &SchedulingSolution{cSolution}

	board := &SolutionBoard{m, TT, SS, false}

	Convey("Subject: Test In-Memory Requests", t, func() {
		Convey("Score function should return 33 when only one of three requests are allocated", func() {
			So(SS.Score(), ShouldEqual, 33)
		})
		Convey("Should Have data after copying and should report the exact same score", func() {
			BB.Copy(SS)
			So(BB.Size(), ShouldBeGreaterThan, 0)
			So(BB.Score(), ShouldEqual, 33)
		})
	})
	Convey("Subject: Test Solution Board", t, func() {
		Convey("Assign should improve current solution score", func() {
			board.Assign("cg1@test.com",req1,0)
			So(board.CurrentSolution.Score(), ShouldBeGreaterThan, 0)
		})
		Convey("Unassign should return the score back to what it was", func() {
			board.Unassign("cg1@test.com",req1,0)
			So(board.CurrentSolution.Score(), ShouldEqual, 0)
		})
		Convey("CanPlay should be false if a conflicting request comes in", func() {
			board.Assign("cg1@test.com",req2,0)
			board.CanPlay("cg1@test.com",req3)
			So(board.CanPlay("cg1@test.com",req3), ShouldBeFalse)
		})
		Convey("CanPlay should be ok if a non-conflicting request comes in", func() {
			board.CanPlay("cg1@test.com",req3)
			So(board.CanPlay("cg2@test.com",req3), ShouldBeTrue)
		})
	})
}

func Test_Appointments(t *testing.T) {
	myFormat := "Jan _2 2006 15:04:05"
	t2, _ := time.Parse(myFormat, "Dec 12 2016 06:00:00")
	app1 := Appointment{t2,1}
	app2 := Appointment{t2,1}
	app3 := Appointment{t2,3} //Defines time interval 6-9

	t3, _ := time.Parse(myFormat, "Dec 12 2016 08:00:00")
	app4 := Appointment{t3,2} //Defines time interval 8-10

	t4, _ := time.Parse(myFormat, "Dec 12 2016 09:00:00")
	app5 := Appointment{t4,2} //Defines time interval 9-11

	t5, _ := time.Parse(myFormat, "Dec 13 2016 09:00:00")
	app6 := Appointment{t5,2} //Defines time interval 9-11

	Convey("Subject: Test Appointments", t, func() {
		Convey("Should Be Equal if they define the exact same time interval", func() {
			So(app1.Equals(app2), ShouldBeTrue)
		})
		Convey("Should Not Be Equal if they don't define the exact same time interval", func() {
			So(app1.Equals(app3), ShouldBeFalse)
		})
		Convey("Should Not Be Equal if they don't occur on the same day", func() {
			So(app5.Equals(app6), ShouldBeFalse)
		})
		Convey("Should Not Conflict if their time intervals do not intersect", func() {
			So(app3.ConflictsWith(app5), ShouldBeFalse)
		})
		Convey("Should Conflict if their time intervals do intersect", func() {
			So(app3.ConflictsWith(app4), ShouldBeTrue)
		})
	})
}
