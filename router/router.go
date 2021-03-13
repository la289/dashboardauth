package router

import (
	"encoding/json"
	"fmt"
	"iotdashboard/controller"
	"log"
	"net"
	"net/http"
)

//TODO: maybe these ports should live in a config file?
var httpPort = ":8080"
var httpsPort = ":9090"

var secHeaders = map[string]string{
	"Strict-Tarnsport-Security": "max-age=63072000; includeSubDomains;",
	"Content-Security-Policy":   "default-src 'self'",
	"X-Frame-Options":           "DENY",
	"X-Content-Type-Options":    "nosniff",
	"Cache-Control":             "no-store",
	"Access-Control-Allow-Origin": "*", //TODO: remove this. Only used for dev while react is served from different spot
}

// Create a struct to read the email and pass from the request
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	CSRF     string `json:"csrf"`
}

func Start(certPath, keyPath string) {
	fmt.Printf("Starting webserver ... \n")
	//start listening for http to redirect to https
	go http.ListenAndServe(httpPort, http.HandlerFunc(redirectTLS))
	//start listening for https and handle requests
	handleRequests(certPath, keyPath)
}

func handleRequests(certPath, keyPath string) {
	mux := http.NewServeMux()
	//TODO: move the finished react code into a local folder
	mux.Handle("/", http.FileServer(http.Dir("../../../frontend-root/iot-dashboard/build/")))
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/logout", logoutHandler)
	mux.HandleFunc("/csrf", csrfHandler)

	log.Fatal(http.ListenAndServeTLS(httpsPort, certPath, keyPath, mux))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received POST/login request")
	addHeaders(w)
	if r.Method != "POST" {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}
	var creds Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// CSRF Validation
	errMsg, errCode := validateCSRF(r, creds)
	if errCode != 0 {
		http.Error(w, errMsg, errCode)
		return
	}

	//Perform Login
	jwt, err := controller.Login(creds.Email, creds.Password)
	if err != nil {
		http.Error(w, "Email and Password do not match", http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "JWT",
		Value:    jwt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrict,
	})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	addHeaders(w)
	log.Printf("Received POST/logout request")
	if r.Method != "POST" {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// CSRF Validation
	errMsg, errCode := validateCSRF(r, creds)
	if errCode != 0 {
		http.Error(w, errMsg, errCode)
		return
	}

	jwtCookie, err := r.Cookie("JWT")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Unauthorized Request", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = controller.Logout(jwtCookie.Value)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
}

func csrfHandler(w http.ResponseWriter, r *http.Request) {
	addHeaders(w)
	log.Printf("Received GET/csrf request")
	if r.Method != "GET" {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}
	// Get the JSON body and decode into credentials

	csrf, err := controller.GenerateRandomString(128)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "CSRF",
		Value:    csrf,
		Secure:   true,
		SameSite: http.SameSiteStrict,
	})
}

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	addHeaders(w)
	host, _, _ := net.SplitHostPort(r.Host)
	u := r.URL
	u.Host = net.JoinHostPort(host, httpsPort[1:])
	u.Scheme = "https"
	target := u.String()

	log.Printf("redirect to: %s", target)
	http.Redirect(w, r, target,
		http.StatusPermanentRedirect)
}

// TODO: turn this into middleware:
func validateCSRF(r *http.Request, creds Credentials) (string, int) {
	csrfCookie, err := r.Cookie("CSRF")
	if err != nil {
		if err == http.ErrNoCookie {
			return "Unauthorized Request", http.StatusUnauthorized
		}
		return "Bad Request", http.StatusBadRequest
	}
	if csrfCookie.Value != creds.CSRF {
		return "Unauthorized Request", http.StatusUnauthorized
	}
	return "", 0
}

//TODO: turn this into middleware:
func addHeaders(w http.ResponseWriter) {
	for key, value := range secHeaders {
		w.Header().Set(key, value)
	}
}
