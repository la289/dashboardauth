package router

import (
	"iotdashboard/controller"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoginHandler(t *testing.T) {
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
		handler := http.HandlerFunc(loginHandler)
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
		for key, value := range secHeaders {
			if response.Header.Get(key) != value {
				t.Errorf("Response is missing headers %v, %v", key, value)
			}
		}
	}
}

func TestLogoutHandler(t *testing.T) {
	token1, err := controller.CreateJWT(60)
	token2, err := controller.CreateJWT(60)

	if err != nil {
		t.Errorf("Failed to generate a token: %v", err)

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
			t.Errorf("Failed to make %v request %v", c.method, err)
		}
		req.AddCookie(&http.Cookie{Name: "JWT", Value: c.jwt})
		req.AddCookie(&http.Cookie{Name: "CSRF", Value: c.csrfC})

		//Record test request through Logout Handler
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(logoutHandler)
		handler.ServeHTTP(rr, req)
		response := rr.Result()

		//Evaluate response for status code
		if rr.Code != c.status {
			t.Errorf("Handler returned wrong status code: got %v want %v. %v",
				rr.Code, c.status, req)
		}

		//Evaluate response for security headers
		for key, value := range secHeaders {
			if response.Header.Get(key) != value {
				t.Errorf("Response is missing headers %v, %v", key, value)
			}
		}
	}
}

func TestCsrfHandler(t *testing.T) {
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
			t.Errorf("Failed to make %v request %v", c.method, err)
		}

		//Record test request through CSRF Handler
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(csrfHandler)
		handler.ServeHTTP(rr, req)
		response := rr.Result()

		//Evaluate response for status code
		if rr.Code != c.status {
			t.Errorf("Handler returned wrong status code: got %v want %v",
				rr.Code, http.StatusOK)
		}
		//Evaluate response for cookies
		if c.status == http.StatusOK {
			cookies := response.Cookies()

			if cookies[0].Name != "CSRF" {
				t.Errorf("Handler returned unexpected token: got %v = %v with error: %v",
					cookies[0].Name, cookies[0].Value, err)
			}
		}
		//Evaluate response for security headers
		for key, value := range secHeaders {
			if response.Header.Get(key) != value {
				t.Errorf("Response is missing headers %v, %v", key, value)
			}
		}
	}
}

// func TestValidateCSRF(t *testing.T){} -> validated through request handler testing

func TestRedirectTLS(t *testing.T) {
	cases := []struct {
		method, path string
		status       int
	}{
		{"GET", "/", http.StatusPermanentRedirect},
		{"POST", "/login", http.StatusPermanentRedirect},
		{"GET", "/csrf?123", http.StatusPermanentRedirect},
	}

	for _, c := range cases {
		// Create test request
		req, err := http.NewRequest(c.method, c.path, nil)
		if err != nil {
			t.Errorf("Failed to make %v request %v", c.method, err)
		}
		//Record test request through redirect Handler
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(redirectTLS)
		handler.ServeHTTP(rr, req)
		response := rr.Result()

		//Evaluate response for status code
		if rr.Code != c.status {
			t.Errorf("Handler returned wrong status code: got %v want %v",
				rr.Code, http.StatusOK)
		}
		//Evaluate response for security headers
		for key, value := range secHeaders {
			if response.Header.Get(key) != value {
				t.Errorf("Response is missing headers %v, %v", key, value)
			}
		}
	}
}
