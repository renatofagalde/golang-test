package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Has(t *testing.T) {
	form := NewForm(nil)

	has := form.Has("whatever")
	if has {
		t.Error("form shows has field when it should not")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")

	form = NewForm(postedData)

	has = form.Has("a")
	if !has {
		t.Error("shows form dows not have field when it should")
	}
}

func TestForm_Required(t *testing.T) {

	request := httptest.NewRequest("POST", "/whatever", nil)

	form := NewForm(request.PostForm)
	form.Required("a", "b", "c")

	if form.Valid() {
		t.Error("Form show valid when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	request, _ = http.NewRequest("POST", "/whatever", nil)
	request.PostForm = postedData

	form = NewForm(request.PostForm)
	form.Required("a", "b", "c")

	if !form.Valid() {
		t.Error("shows post does not have required fields, when it does")
	}

}

func TestForm_check(t *testing.T) {
	form := NewForm(nil)
	form.Check(false, "password", "password is required")

	if form.Valid() {
		t.Error("Valid() returns false,and it should be true when calling Check()")
	}
}

func TestForm_ErrorGet(t *testing.T) {
	form := NewForm(nil)
	form.Check(false, "password", "password is required")

	s := form.Errors.Get("password")

	if len(s) == 0 {
		t.Error("should have an error returned from Get, but do not")
	}

	s = form.Errors.Get("whaterver")
	if len(s) != 0 {
		t.Error("should not have an error, but got one")
	}
}
