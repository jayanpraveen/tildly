package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	m "github.com/jayanpraveen/tildly/entity"
	"github.com/jayanpraveen/tildly/service"
)

type HandleTester func(method string, url string, params url.Values) *httptest.ResponseRecorder

func generateHandleTester(t *testing.T, handleFunc http.Handler) HandleTester {
	return func(method string, url string, params url.Values) *httptest.ResponseRecorder {

		if url == "" {
			url = "/"
		}

		r, err := http.NewRequest(method, url, strings.NewReader(params.Encode()))

		data, _ := ioutil.ReadAll(r.Body)
		t.Log(string(data))

		if err != nil {
			t.Errorf("%v", err)
		}

		w := httptest.NewRecorder()

		handleFunc.ServeHTTP(w, r)

		return w

	}
}

func assertEquals(comp string, exp interface{}, act interface{}) error {
	if exp != act {
		return fmt.Errorf(comp, "Expected: %v, Actual: %v", exp, act)
	}
	return nil
}

type MockService struct {
	ch service.UrlCache
}

func (m *MockService) SaveUrl(longUrl string) error {
	return nil
}

func (m *MockService) GetUrlByHash(hash string) (*m.Url, error) {
	return nil, nil
}

func Test_handleLongUrl(t *testing.T) {

	t.Run("Can register valid url", func(t *testing.T) {
		// url := m.Url{Hash: "XikHs", LongUrl: "https://go.dev", CreatedAt: "NOW"}

		actual := "Url created!"

		us := &MockService{ch: &service.CacheRepo{}}
		uh := NewUrlHandler(us)

		req := httptest.NewRequest(http.MethodPost, "/api/longUrl", strings.NewReader(`{"longUrl": "http://go.dev/docs"}`))
		res := httptest.NewRecorder()

		h := uh.handleLongUrl()
		h(res, req)

		err := assertEquals("longUrl Hand", res.Code, http.StatusCreated)
		if err != nil {
			t.Error(err)
		}

		if res.Body.String() != actual {
			t.Errorf("expected %q got %q", res.Body.String(), actual)
		}

	})

	// todo: perform invalid test for this handler

}
