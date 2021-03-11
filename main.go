package main

import(
	"iotdashboard/router"
	// "iotdashboard/controller"
)


func init() {
}

func main() {

	router.Start("server-cert.pem","server-key.pem")
}
