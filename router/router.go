package router

import (
	"encoding/json"
	"fmt"
	"iotdashboard/controller"
	"log"
	"net"
	"net/http"
)

const defaultHTTPport  = ":8080"
const defaultHTTPSport  = ":9090"


type Router struct{
	ctrlr controller.Controller
	httpPort string
	httpsPort string
}

// Credentials is a struct that holds the email, password, and CSRF token of a request
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	CSRF     string `json:"csrf"`
}

func NewRouter(httpPort, httpsPort string) (Router, error) {
	ctrlr,err := controller.NewController()
	if err != nil{
		return Router{}, err
	}
	return Router{ctrlr, httpPort, httpsPort}, nil
}

func (rtr *Router) Start(certPath, keyPath string) error {
	fmt.Printf("Starting webserver ... \n")
	//start listening for http to redirect to https
	go http.ListenAndServe(rtr.httpPort, http.HandlerFunc(rtr.redirectTLS))

	//start listening for https and handle requests
	err := rtr.handleRequests(certPath, keyPath)
	if err != nil{
		return err
	}
	return nil
}

func (rtr *Router) handleRequests(certPath, keyPath string) error {
	mux := http.NewServeMux()
	//TODO: move the finished react code into a local folder
	mux.Handle("/", http.FileServer(http.Dir("../../../frontend-root/iot-dashboard/build/")))
	mux.HandleFunc("/login", rtr.loginHandler)
	mux.HandleFunc("/logout", rtr.logoutHandler)
	mux.HandleFunc("/csrf", rtr.csrfHandler)

	err := http.ListenAndServeTLS(rtr.httpsPort, certPath, keyPath, mux)
	if err != nil {
		return err
	}
	return nil
}

func (rtr *Router) loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received POST/login request")
	rtr.addHeaders(w)
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
	errMsg, errCode := rtr.validateCSRF(r, creds)
	if errCode != 0 {
		http.Error(w, errMsg, errCode)
		return
	}

	//Perform Login
	jwt, err := rtr.ctrlr.Login(creds.Email, creds.Password)
	if err != nil {
		http.Error(w, "Email and Password do not match", http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "JWT",
		Value:    jwt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func (rtr *Router) logoutHandler(w http.ResponseWriter, r *http.Request) {
	rtr.addHeaders(w)
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
	errMsg, errCode := rtr.validateCSRF(r, creds)
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

	err = rtr.ctrlr.Logout(jwtCookie.Value)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
}

func (rtr *Router) csrfHandler(w http.ResponseWriter, r *http.Request) {
	rtr.addHeaders(w)
	log.Printf("Received GET/csrf request")
	if r.Method != "GET" {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}
	// Get the JSON body and decode into credentials

	csrf, err := rtr.ctrlr.TokenUtil.GenerateRandomString(128)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "CSRF",
		Value:    csrf,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func (rtr *Router) redirectTLS(w http.ResponseWriter, r *http.Request) {
	rtr.addHeaders(w)
	host, _, _ := net.SplitHostPort(r.Host)
	u := r.URL
	u.Host = net.JoinHostPort(host, defaultHTTPSport [1:])
	u.Scheme = "https"
	target := u.String()

	log.Printf("redirect to: %s", target)
	http.Redirect(w, r, target,
		http.StatusPermanentRedirect)
}

// TODO: turn this into middleware:
func (rtr *Router) validateCSRF(r *http.Request, creds Credentials) (string, int) {
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
func (rtr *Router) addHeaders(w http.ResponseWriter) {
	for key, value := range rtr.getSecHeaders() {
		w.Header().Set(key, value)
	}
}


func (rtr *Router) getSecHeaders() map[string]string {
	return map[string]string{
		"Strict-Tarnsport-Security": "max-age=63072000; includeSubDomains;",
		"Content-Security-Policy":   "default-src 'self'",
		"X-Frame-Options":           "DENY",
		"X-Content-Type-Options":    "nosniff",
		"Cache-Control":             "no-store",
		"Access-Control-Allow-Origin": "*", //TODO: remove this. Only used for dev while react is served from different spot
	}
}
