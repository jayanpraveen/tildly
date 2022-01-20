package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	m "github.com/jayanpraveen/tildly/entity"
)

type HandleTester func(method string, url string, body string) *httptest.ResponseRecorder

func generateHandleTester(t *testing.T, handleFunc http.Handler, expSC int, expBody string) HandleTester {
	return func(method string, url string, body string) *httptest.ResponseRecorder {

		if url == "" {
			url = "/"
		}

		if method == "" {
			method = http.MethodGet
		}

		res := httptest.NewRecorder()
		req := httptest.NewRequest(method, url, strings.NewReader(body))

		// Calling the given handler
		handleFunc.ServeHTTP(res, req)

		// Check http status
		assertEquals(t, expSC, res.Code)

		// Check output
		assertEquals(t, expBody, res.Body.String())

		return res
	}
}

/*
 * This functions follows the convention similar to JUnit's `assertEquals` method.
 * The exp or expected is the value the user expects from the program (the known one).
 * The act or actual is the value the program produces (the unknown one).
 */
func assertEquals(t *testing.T, exp interface{}, act interface{}) {
	if exp != act {
		t.Errorf("Expected: %v, Actual: %v", exp, act)
	}
}

type MockService struct{}

func (m *MockService) SaveUrl(longUrl string) error {
	return nil
}

func (m *MockService) GetUrlByHash(hash string) (*m.Url, error) {
	return nil, nil
}

func Test_handleLongUrl(t *testing.T) {

	t.Run("register valid url", func(t *testing.T) {
		srv := &MockService{}
		uh := NewUrlHandler(srv)
		h := uh.handleLongUrl()

		expSC := http.StatusCreated
		expBody := "Url created!"

		method := http.MethodPost
		url := "/api/longUrl"
		body := `{"longUrl": "http://go.dev/docs"}`

		gh := generateHandleTester(t, h, expSC, expBody)
		gh(method, url, body)

	})
}
