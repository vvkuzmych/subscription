package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"subscription/data"
	"testing"

	"github.com/stretchr/testify/assert"
)

var pageTests = []struct {
	name               string
	url                string
	expectedStatusCode int
	handler            http.HandlerFunc
	sessionData        map[string]any
	expectedHTML       string
}{
	{
		name:               "home",
		url:                "/",
		expectedStatusCode: http.StatusOK,
		handler:            testApp.HomePage,
	},
	{
		name:               "login",
		url:                "/login",
		expectedStatusCode: http.StatusOK,
		handler:            testApp.LoginPage,
		expectedHTML:       `<h1 class="mt-5">Login</h1>`,
	},
	{
		name:               "logout",
		url:                "/logout",
		expectedStatusCode: http.StatusOK,
		handler:            testApp.LoginPage,
		expectedHTML:       `<h1 class="mt-5">Login</h1>`,
		sessionData: map[string]any{
			"userID": 1,
			"user":   data.User{},
		},
	},
	{
		name:               "logout",
		url:                "/logout",
		expectedStatusCode: http.StatusSeeOther,
		handler:            testApp.Logout,
		sessionData: map[string]any{
			"userID": 1,
			"user":   data.User{},
		},
	},
}

func TestConfig_Pages(t *testing.T) {
	pathToTemplates = "./templates"

	for _, e := range pageTests {
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", e.url, nil)

		ctx := getCtx(req)
		req = req.WithContext(ctx)

		if len(e.sessionData) > 0 {
			for key, value := range e.sessionData {
				testApp.Session.Put(ctx, key, value)
			}
		}
		//handler := http.HandlerFunc(testApp.HomePage)
		e.handler.ServeHTTP(rr, req)

		assert.Equal(t, e.expectedStatusCode, rr.Code)
		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s handler returned wrong status code: got %v want %v", e.name, rr.Code, e.expectedStatusCode)
		}

		if len(e.expectedHTML) > 0 {
			html := rr.Body.String()
			assert.Contains(t, html, e.expectedHTML)
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("handler returned unexpected body: got %v want %v", html, e.expectedHTML)
			}
		}
	}
}
