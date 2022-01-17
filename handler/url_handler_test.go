package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type HandleTester func(method string, params url.Values) *httptest.ResponseRecorder

func GenerateHandleTester(t *testing.T, handleFunc http.Handler) HandleTester {
	return func(method string, params url.Values) *httptest.ResponseRecorder {

		req, err := http.NewRequest(method, "/", strings.NewReader(params.Encode()))

		if err != nil {
			t.Errorf("%v", err)
		}

		w := httptest.NewRecorder()

		handleFunc.ServeHTTP(w, req)
		return w

	}
}

func Test_handleIndex(t *testing.T) {

	uh := UrlHandler{}
	hi := uh.handleIndex()
	test := GenerateHandleTester(t, hi)

	w := test("GET", url.Values{})

	if w.Code != http.StatusOK {
		t.Errorf("Expected: %v, Actual: %v", http.StatusOK, w.Code)
	}

}
