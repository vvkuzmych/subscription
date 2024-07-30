package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConfig_AddDefaultData(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	testApp.Session.Put(ctx, "flash", "flash")
	testApp.Session.Put(ctx, "warning", "warning")
	testApp.Session.Put(ctx, "error", "error")

	td := testApp.AddDefaultData(&TemplateData{}, req)

	assert.Equal(t, "flash", td.Flash)
	assert.Equal(t, "warning", td.Warning)
	assert.Equal(t, "error", td.Error)

	if td.Flash != "flash" {
		t.Errorf("Flash expected: %s, got: %s", "flash", td.Flash)
	}
	if td.Warning != "warning" {
		t.Errorf("Warning expected: %s, got: %s", "warning", td.Warning)
	}
	if td.Error != "error" {
		t.Errorf("Error expected: %s, got: %s", "error", td.Error)
	}
}

func TestConfig_IsAuthenticated(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	auth := testApp.isAuthenticated(req)
	assert.False(t, auth)
	if auth {
		t.Errorf("Authentication expected: %t, got: %t", true, auth)
	}
	testApp.Session.Put(ctx, "userID", 1)

	auth = testApp.isAuthenticated(req)
	assert.True(t, auth)
	if !auth {
		t.Errorf("Authentication expected: %t, got: %t", false, auth)
	}
}

func TestConfig_render(t *testing.T) {
	pathToTemplates = "./templates"
	rr := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	testApp.render(rr, req, "home.page.gohtml", &TemplateData{})

	assert.Equal(t, http.StatusOK, rr.Code)
	if rr.Code != http.StatusOK {
		t.Errorf("Wrong status code. Expected: %d, got: %d", http.StatusOK, rr.Code)
	}
}
