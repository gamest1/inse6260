// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE handle.

// Package serviceTests implements tests for all system services.
package serviceTests

import (
	"testing"

	"github.com/goinggo/beego-mgo/services/userService"
	. "github.com/smartystreets/goconvey/convey"
)

// Test_FindUsersOfKind checks that we can fetch all users of certain type
func Test_FindUsersOfKind(t *testing.T) {
	service := Prepare()
	defer Finish(service)

	userType := "cg"

	careGivers, err := userService.FindUsersOfKind(service, userType)

	Convey("Subject: Test FindUsersOfKind", t, func() {
		Convey("Should Be Able To Perform A Search", func() {
			So(err, ShouldEqual, nil)
		})
		Convey("Should Have At Least One Care Giver", func() {
			So(len(careGivers), ShouldBeGreaterThan, 0)
		})
	})
}

// Test_Region checks the region service call is working
func Test_FetchAllLanguagesForKind(t *testing.T) {
	service := Prepare()
	defer Finish(service)

	kind := "cg"

	languages, err := userService.FetchAllLanguagesForKind(service, kind)

	Convey("Subject: FetchAllLanguagesForKind", t, func() {
		Convey("Should Be Able To Perform A Search", func() {
			So(err, ShouldEqual, nil)
		})
		Convey("Should Have Language Data", func() {
			So(len(languages), ShouldBeGreaterThan, 0)
		})
	})
}
