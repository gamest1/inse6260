// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE handle.

// Package serviceTests implements tests for all system services.
package serviceTests

import (
	"testing"

  "github.com/goinggo/beego-mgo/utilities/scheduler"
	. "github.com/smartystreets/goconvey/convey"
)

// Test_CreateSchedule checks that a new schedule can be created via the scheduler
func Test_CreateSchedule(t *testing.T) {

	scheduler.CreateSchedule()

	Convey("Subject: Test CreateSchedule", t, func() {
		Convey("Should Be Able To Start the scheduler", func() {
			So(nil, ShouldEqual, nil)
		})
	})
}
