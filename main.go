package main

import (
	"iotdashboard/router"
	"log"
)

func main() {
	router, err := router.NewRouter(":8080", ":9090", "server-cert.pem", "server-key.pem")
	if err != nil {
		log.Fatal(err)
	}

	err = router.Ctrlr.PSQL.AddNewUser("e@g.c", "test")
	if err != nil {
		log.Printf("Admin user already exists in DB")
	}

	err = router.Start()
	if err != nil {
		panic(err)
	}
}
