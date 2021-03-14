package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoginHandler(t *testing.T) {
	router,err := NewRouter(":8080",":9090")
	if err != nil{
		t.Errorf("Could not initialize router: %v \n", err)
	}

	cases := []struct {
		method, path, email, pass, csrfC, csrfB string
		status                                  int
	}{
		{"POST", "/login", "user@gmail.com", "S3cure3Pa$$", "123", "123", http.StatusOK},
		{"POST", "/login11", "user@gmail.com", "S3cure3Pa$$", "123", "123", http.StatusOK},
		{"POST", "/login", "user@gmail.com", "S3cure3Pa$$", "1234", "123", http.StatusUnauthorized},
		{"POST", "/login?1231", "user@gmail.com", "s3cure3Pa$$", "123", "123", http.StatusUnauthorized},
		{"GET", "/login", "user@gmail.com", "S3cure3Pa$$", "123", "123", http.StatusMethodNotAllowed},
	}

	for _, c := range cases {
		// Create test request
		bodyReader := strings.NewReader(`{"email":"` + c.email +
			`","password":"` + c.pass +
			`","csrf":"` + c.csrfB + `"}`)
		req, err := http.NewRequest(c.method, c.path, bodyReader)
		if err != nil {
			t.Errorf("Failed to make %v request %v", c.method, err)
		}
		req.AddCookie(&http.Cookie{Name: "CSRF", Value: c.csrfC})

		//Record test request through Login Handler
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(router.loginHandler)
		handler.ServeHTTP(rr, req)
		response := rr.Result()

		//Evaluate response for status code
		if rr.Code != c.status {
			t.Errorf("Handler returned wrong status code: got %v want %v. %v",
				rr.Code, c.status, req)
		}

		//Evaluate response for cookies
		if c.status == http.StatusOK {
			cookies := response.Cookies()

			if cookies[0].Name != "JWT" {
				t.Errorf("Handler returned unexpected token: got %v = %v with error: %v",
					cookies[0].Name, cookies[0].Value, err)
			}
		}

		//Evaluate response for security headers
		for key, value := range router.getSecHeaders() {
			if response.Header.Get(key) != value {
				t.Errorf("Response is missing headers %v, %v", key, value)
			}
		}
	}
}

func TestLogoutHandler(t *testing.T) {
	router,err := NewRouter(":8080",":9090")
	if err != nil{
		t.Errorf("Could not initialize router: %v \n", err)
	}

	token1, err1 := router.ctrlr.TokenUtil.CreateJWT(60)
	token2, err2 := router.ctrlr.TokenUtil.CreateJWT(60)
	if (err1 != nil || err2 != nil) {
		t.Errorf("Failed to generate a token. Errors: \n %v \n %v \n", err1, err2 )
	}
	cases := []struct {
		method, path, jwt, csrfC, csrfB string
		status                          int
	}{
		{"POST", "/logout", token1, "123", "123", http.StatusOK},
		//trying to log out an already logged out token
		{"POST", "/logout", token1, "123", "123", http.StatusUnauthorized},
		{"POST", "/logout?abc", token2 + "9", "123", "123", http.StatusUnauthorized},
		{"POST", "/logout", token2, "1234", "123", http.StatusUnauthorized},
		{"GET", "/logout", "user@gmail.com", "123", "123", http.StatusMethodNotAllowed},
	}

	for _, c := range cases {
		// Create test request
		bodyReader := strings.NewReader(`{"csrf":"` + c.csrfB + `"}`)

		req, err := http.NewRequest(c.method, c.path, bodyReader)
		if err != nil {
			t.Errorf("Failed to make %v request %v \n", c.method, err)
		}
		req.AddCookie(&http.Cookie{Name: "JWT", Value: c.jwt})
		req.AddCookie(&http.Cookie{Name: "CSRF", Value: c.csrfC})

		//Record test request through Logout Handler
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(router.logoutHandler)
		handler.ServeHTTP(rr, req)
		response := rr.Result()

		//Evaluate response for status code
		if rr.Code != c.status {
			t.Errorf("Handler returned wrong status code: got %v want %v. %v \n",
				rr.Code, c.status, req)
		}

		//Evaluate response for security headers
		for key, value := range router.getSecHeaders() {
			if response.Header.Get(key) != value {
				t.Errorf("Response is missing headers %v, %v \n", key, value)
			}
		}
	}
}

func TestCsrfHandler(t *testing.T) {
	router,err := NewRouter(":8080",":9090")
	if err != nil{
		t.Errorf("Could not initialize router: %v \n", err)
	}

	cases := []struct {
		method, path string
		status       int
	}{
		{"GET", "/csrf", http.StatusOK},
		{"GET", "/csrf?1231", http.StatusOK},
		{"POST", "/csrf", http.StatusMethodNotAllowed},
	}

	for _, c := range cases {
		// Create test request
		req, err := http.NewRequest(c.method, c.path, nil)
		if err != nil {
			t.Errorf("Failed to make %v request %v \n", c.method, err)
		}

		//Record test request through CSRF Handler
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(router.csrfHandler)
		handler.ServeHTTP(rr, req)
		response := rr.Result()

		//Evaluate response for status code
		if rr.Code != c.status {
			t.Errorf("Handler returned wrong status code: got %v want %v \n",
				rr.Code, http.StatusOK)
		}
		//Evaluate response for cookies
		if c.status == http.StatusOK {
			cookies := response.Cookies()

			if cookies[0].Name != "CSRF" {
				t.Errorf("Handler returned unexpected token: got %v = %v with error: %v \n",
					cookies[0].Name, cookies[0].Value, err)
			}
		}
		//Evaluate response for security headers
		for key, value := range router.getSecHeaders() {
			if response.Header.Get(key) != value {
				t.Errorf("Response is missing headers %v, %v \n", key, value)
			}
		}
	}
}

// func TestValidateCSRF(t *testing.T){} -> validated through request handler testing

func TestRedirectTLS(t *testing.T) {
	router,err := NewRouter(":8080",":9090")
	if err != nil{
		t.Errorf("Could not initialize router: %v \n", err)
	}

	cases := []struct {
		method, path, newPath string
		status       int
	}{
		{"GET", "http://192.0.0.1:8080", "https://192.0.0.1:9090", http.StatusPermanentRedirect},
		{"POST", "http://192.0.0.1:8080/login", "https://192.0.0.1:9090/login", http.StatusPermanentRedirect},
		{"GET", "http://192.1.1.1:8080/csrf?123", "https://192.1.1.1:9090/csrf?123", http.StatusPermanentRedirect},
	}

	for _, c := range cases {
		// Create test request
		req, err := http.NewRequest(c.method, c.path, nil)
		if err != nil {
			t.Errorf("Failed to make %v request %v \n", c.method, err)
		}
		//Record test request through redirect Handler
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(router.redirectTLS)
		handler.ServeHTTP(rr, req)
		response := rr.Result()

		//Evaluate response for status code
		if rr.Code != c.status {
			t.Errorf("Handler returned wrong status code: got %v want %v \n",
				rr.Code, http.StatusOK)
		}
		//Evaluate redirect url
		redirectURL, err := response.Location()
		if (redirectURL.String() != c.newPath || err != nil) {
			t.Errorf("Expected URL: %v - Received URL: %v \n", c.newPath, redirectURL)
		}

		//Evaluate response for security headers
		for key, value := range router.getSecHeaders() {
			if response.Header.Get(key) != value {
				t.Errorf("Response is missing headers %v, %v \n", key, value)
			}
		}
	}
}
