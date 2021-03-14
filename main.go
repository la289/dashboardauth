package main

import (
	"iotdashboard/router"
)

func main() {
	router,err := router.NewRouter(":8080",":9090")
	if err != nil{
		panic(err)
	}


	err = router.Start("server-cert.pem", "server-key.pem")
	if err != nil {
		panic(err)
	}
}
