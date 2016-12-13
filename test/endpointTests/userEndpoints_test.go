// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE handle.

// Package endpointTests implements tests for the buoy endpoints.
package endpointTests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	//"encoding/json"

	"github.com/astaxie/beego"
	log "github.com/goinggo/tracelog"
	. "github.com/smartystreets/goconvey/convey"
)

// We support 2 AJAX get endPoints:
// /requests/:userId  and
// /user/display/:userId
// 5832b0cd1a4c35385a885b19

// TestRequestsForUser is a sample to run an endpoint test
func TestRequestsForUser(t *testing.T) {
	r, _ := http.NewRequest("GET", "/requests/egarro@transcriptics.com", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	log.Trace("testing", "TestRequestsForUser", "Code[%d]\n%s", w.Code, w.Body.String())

	// var response struct {
	// 	RequestID string `json:"req_id"`
	// }
	// json.Unmarshal(w.Body.Bytes(), &response)

	//Not sure why this test returns 500?!
	Convey("Subject: Test Requests Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
	})
}

// TestRequestsForInvalidUser is a sample to run an endpoint for an invalid user
// an empty result set
// func TestRequestsForInvalidUser(t *testing.T) {
// 	r, _ := http.NewRequest("GET", "/requests/nonexistent@transcriptics.com", nil)
// 	w := httptest.NewRecorder()
// 	beego.BeeApp.Handlers.ServeHTTP(w, r)
//
// 	log.Trace("testing", "TestStation", "Code[%d]\n%s", w.Code, w.Body.String())
//
// 	var err struct {
// 		Errors []string `json:"errors"`
// 	}
// 	json.Unmarshal(w.Body.Bytes(), &err)
//
// 	Convey("Subject: Test Requests Endpoint\n", t, func() {
// 		Convey("Status Code Should Be 409", func() {
// 			So(w.Code, ShouldEqual, 409)
// 		})
// 		Convey("The Result Should Not Be Empty", func() {
// 			So(w.Body.Len(), ShouldBeGreaterThan, 0)
// 		})
// 		Convey("The Should Be An Error In The Result", func() {
// 			So(len(err.Errors), ShouldEqual, 1)
// 		})
// 	})
// }
