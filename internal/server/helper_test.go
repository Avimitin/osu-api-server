package server

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestGetFormValue(t *testing.T) {
	form := url.Values{"key": {"val"}}
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	val := getFormValue(req, "key")
	if val != "val" {
		t.Errorf("got %s want val", val)
	}
}
