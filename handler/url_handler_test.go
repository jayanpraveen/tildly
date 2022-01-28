package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	m "github.com/jayanpraveen/tildly/entity"
)

type HandleTester func(method string, url string, body string) *httptest.ResponseRecorder

func generateHandleTester(t *testing.T, handleFunc http.Handler, expSC int) HandleTester {
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

		return res
	}
}

/*
 * This functions follows the convention similar to JUnit's `assertEquals` method.
 * The exp or expected is the value the user expects from the program (the known one).
 * The act or actual is the value the program produces (the unknown one).
 */
func assertEquals(t *testing.T, exp interface{}, act interface{}) {
	if !cmp.Equal(exp, act) {
		t.Errorf(" Expected: %v, Actual: %v string", exp, act)
	}
}

type MockService struct {
	SaveUrlFunc      func(longUrl string, expireAt int64) error
	GetUrlByHashFunc func(hash string) (*m.Url, error)
}

func (m *MockService) SaveUrl(longUrl string, exipreAt int64) error {
	return m.SaveUrlFunc(longUrl, 1257894000)
}

func (m *MockService) GetUrlByHash(hash string) (*m.Url, error) {
	return m.GetUrlByHashFunc(hash)
}

func Test_handleIndex(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	uh := UrlHandler{}
	h := uh.handleIndex()
	h(res, req)

	// Check http status
	assertEquals(t, http.StatusOK, res.Code)

	// Check response output
	assertEquals(t, []byte("tildly !\n"), res.Body.Bytes())
}

func Test_handleLongUrl(t *testing.T) {

	t.Run("post valid url", func(t *testing.T) {
		srv := &MockService{
			SaveUrlFunc: func(longUrl string, exipreAt int64) error {
				return nil
			},
		}
		uh := NewUrlHandler(srv)
		h := uh.handleLongUrl()

		expSC := http.StatusCreated
		expBody := "tildly url created!"

		method := http.MethodPost
		url := "/api/longUrl"
		body := `{"longUrl": "http://go.dev/docs"}`

		gh := generateHandleTester(t, h, expSC)
		res := gh(method, url, body)

		// Check response output
		assertEquals(t, expBody, res.Body.String())

	})

	t.Run("post invalid url", func(t *testing.T) {
		srv := &MockService{
			SaveUrlFunc: func(longUrl string, exipreAt int64) error {
				return nil
			},
		}
		uh := NewUrlHandler(srv)
		h := uh.handleLongUrl()

		expSC := http.StatusBadRequest
		expBody := fmt.Sprintln(`{"status":400,"msg":"Not a valid URL"}`) // Using `ln` becuase Json.NewDecoder adds \n at eof while marshlling

		method := http.MethodPost
		url := "/api/longUrl"
		body := `{"longUrl": "htt/go.dev/docs"}`

		gh := generateHandleTester(t, h, expSC)
		res := gh(method, url, body)

		// Check response output
		assertEquals(t, expBody, res.Body.String())

		// when given empty url
		body = `{"longUrl": ""}`
		res = gh(method, url, body)
		assertEquals(t, expBody, res.Body.String())

	})

	t.Run("post invalid JSON", func(t *testing.T) {
		srv := &MockService{
			SaveUrlFunc: func(longUrl string, exipreAt int64) error {
				return nil
			},
		}
		uh := NewUrlHandler(srv)
		h := uh.handleLongUrl()

		expSC := http.StatusBadRequest
		expBody := fmt.Sprintln(`{"status":400,"msg":"Not a proper JSON format"}`)

		method := http.MethodPost
		url := "/api/longUrl"
		body := `{"longUrl" "htt/go.dev/docs"}` // Invalid JSON: missing ':' in longUrl

		gh := generateHandleTester(t, h, expSC)
		res := gh(method, url, body)

		// Check response output
		assertEquals(t, expBody, res.Body.String())
	})

	t.Run("returns interal server error", func(t *testing.T) {
		srv := &MockService{
			SaveUrlFunc: func(longUrl string, exipreAt int64) error {
				return fmt.Errorf("Failed to save url")
			},
		}

		uh := NewUrlHandler(srv)
		h := uh.handleLongUrl()

		expSC := http.StatusInternalServerError
		expBody := fmt.Sprintln(`{"status":500,"msg":"Internal server error"}`)

		method := http.MethodPost
		url := "/api/longUrl"
		body := `{"longUrl": "http://go.dev/docs"}`

		gh := generateHandleTester(t, h, expSC)
		res := gh(method, url, body)

		// Check response output
		assertEquals(t, expBody, res.Body.String())
	})

}

func Test_handleShortUrl(t *testing.T) {

	t.Run("short url exist", func(t *testing.T) {
		srv := &MockService{
			GetUrlByHashFunc: func(hash string) (*m.Url, error) {
				return &m.Url{
					Hash:      "XikHsqW",
					LongUrl:   "https://go.dev",
					CreatedAt: 1257894000,
				}, nil
			},
		}

		uh := NewUrlHandler(srv)
		h := uh.handleShortUrl()

		expSC := http.StatusPermanentRedirect

		method := http.MethodGet
		url := "/XikHsqW"
		body := ""

		gh := generateHandleTester(t, h, expSC)
		gh(method, url, body)

	})

	t.Run("short url doesn't exist", func(t *testing.T) {
		srv := &MockService{
			GetUrlByHashFunc: func(hash string) (*m.Url, error) {
				return nil, fmt.Errorf("short url not found")
			},
		}

		uh := NewUrlHandler(srv)
		h := uh.handleShortUrl()

		expSC := http.StatusNotFound

		method := http.MethodGet
		url := "/XikHsqW"
		body := ""

		gh := generateHandleTester(t, h, expSC)
		gh(method, url, body)
	})
}

// test fixture
func init() {
	if err := os.Chdir("./testdata"); err != nil {
		panic(err)
	}
}
