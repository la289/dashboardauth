package main

import (
	"iotdashboard/router"
	"log"
)

func main() {
	router, err := router.NewRouter(":8080", ":9090")
	if err != nil {
		panic(err)
	}

	err = router.Ctrlr.PSQL.AddNewUser("e@g.c", "test")
	if err != nil {
		log.Printf("Admin user already exists in DB")
	}

	err = router.Start("server-cert.pem", "server-key.pem")
	if err != nil {
		panic(err)
	}
}
