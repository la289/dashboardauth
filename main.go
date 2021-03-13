package main

import (
	"iotdashboard/router"
)

func main() {

	router.Start("server-cert.pem", "server-key.pem")
}
