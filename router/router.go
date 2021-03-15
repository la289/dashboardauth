package router

import (
	"encoding/json"
	"iotdashboard/controller"
	"log"
	"net"
	"net/http"
)


type RouterService struct{
	Ctrlr controller.ControllerService
	httpPort string
	httpsPort string
}

// Credentials is a struct that holds the email, password, and CSRF token of a request
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	CSRF     string `json:"csrf"`
}

func NewRouter(httpPort, httpsPort string) (RouterService, error) {
	Ctrlr,err := controller.NewController()
	if err != nil{
		return RouterService{}, err
	}
	return RouterService{Ctrlr, httpPort, httpsPort}, nil
}

func (rtr *RouterService) Start(certPath, keyPath string) error {
	log.Printf("Starting webserver ... \n")
	//start listening for http to redirect to https
	go http.ListenAndServe(rtr.httpPort, http.HandlerFunc(rtr.redirectTLS))

	//start listening for https and handle requests
	err := rtr.handleRequests(certPath, keyPath)
	if err != nil{
		log.Printf("HandleRequests Error: %v /n", err)
		return err
	}
	return nil
}

func (rtr *RouterService) handleRequests(certPath, keyPath string) error {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("iotdbfrontend/build/")))
	mux.HandleFunc("/login", rtr.loginHandler)
	mux.HandleFunc("/logout", rtr.logoutHandler)
	mux.HandleFunc("/csrf", rtr.csrfHandler)

	log.Printf("Running! \n")
	err := http.ListenAndServeTLS(rtr.httpsPort, certPath, keyPath, mux)
	if err != nil {
		log.Printf("ListenAndServeTLS Error: %v /n", err)
		return err
	}
	return nil
}

func (rtr *RouterService) loginHandler(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("LoginHandler Error: %v /n", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// CSRF Validation
	rtr.validateCSRF(w, r, creds)
	log.Printf("CSRF Validated")

	//Perform Login
	jwt, err := rtr.Ctrlr.Login(creds.Email, creds.Password)
	if err != nil {
		http.Error(w, "Email and Password do not match", http.StatusUnauthorized)
		return
	}
	log.Printf("Logged In")

	http.SetCookie(w, &http.Cookie{
		Name:     "JWT",
		Value:    jwt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	log.Printf("Responding")

}

func (rtr *RouterService) logoutHandler(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("LogoutHandler Error: %v /n", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// CSRF Validation
	rtr.validateCSRF(w, r, creds)

	jwtCookie, err := r.Cookie("JWT")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Unauthorized Request", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = rtr.Ctrlr.Logout(jwtCookie.Value)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
}

func (rtr *RouterService) csrfHandler(w http.ResponseWriter, r *http.Request) {
	rtr.addHeaders(w)
	log.Printf("Received GET/csrf request")
	if r.Method != "GET" {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	csrf, err := rtr.Ctrlr.TokenUtil.GenerateRandomString(128)
	if err != nil {
		log.Printf("CSRFHandler Error: %v /n", err)
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

func (rtr *RouterService) redirectTLS(w http.ResponseWriter, r *http.Request) {
	rtr.addHeaders(w)
	//discarding old port value
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		log.Printf("Redirect TLS Error: %v /n", err)
		return
	}
	u := r.URL
	u.Host = net.JoinHostPort(host, rtr.httpsPort [1:])
	u.Scheme = "https"
	target := u.String()

	log.Printf("redirect to: %s", target)
	http.Redirect(w, r, target,
		http.StatusPermanentRedirect)

}

// TODO: turn this into middleware:
func (rtr *RouterService) validateCSRF(w http.ResponseWriter, r *http.Request, creds Credentials) {
	csrfCookie, err := r.Cookie("CSRF")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
		log.Printf("ValidateCSRF Error: %v /n", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	if csrfCookie.Value != creds.CSRF {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

//TODO: turn this into middleware:
func (rtr *RouterService) addHeaders(w http.ResponseWriter) {
	for key, value := range rtr.getSecHeaders() {
		w.Header().Set(key, value)
	}
}


func (rtr *RouterService) getSecHeaders() map[string]string {
	return map[string]string{
		"Strict-Tarnsport-Security": "max-age=63072000; includeSubDomains;",
		"Content-Security-Policy":   "default-src 'self'",
		"X-Frame-Options":           "DENY",
		"X-Content-Type-Options":    "nosniff",
		"Cache-Control":             "no-store",
	}
}
