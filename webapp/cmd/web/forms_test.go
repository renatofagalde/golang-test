package main

import (
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
