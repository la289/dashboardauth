package controller

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	// "fmt"
)

func TestLoginAndLogout(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("failed to open sqlmock database: %v \n", err)
	}
	defer db.Close()

	controller, err := NewController()
	if err != nil {
		t.Errorf("Failed to initialize new controller: %v \n", err)
	}
	//overwrite db connection with mock
	controller.PSQL.DB = db

	cases := []struct {
		email, password, hashedPassword string
		success                         bool
	}{
		{"user@gmail.com", "S3cure3Pa$$", "$2a$10$cRmL5Rtm0bunl1uqYAP.8OfJE36RUkvMcX3.v0kJyY2JBhalX4KEG", true},
		{"user@gmail.com", "wrongpass", "$2a$10$cRmL5Rtm0bunl1uqYAP.8OfJE36RUkvMcX3.v0kJyY2JBhalX4KEG", false},
		{"1000@doesntexist.com", "bloop", "anotherfakehash", false},
	}

	for _, c := range cases {
		rows := sqlmock.NewRows([]string{"password"}).
			AddRow(c.hashedPassword)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT password from users WHERE email = $1")).
			WillReturnRows(rows)

		jwt, err := controller.Login(c.email, c.password)
		if (err != nil && c.success) || (err == nil && !c.success) {
			t.Errorf("Login failed (email: %s - pass: %s). Error: %v", c.email, c.password, err)
		}

		// if successful login, test logout
		if c.success == true {
			err = controller.Logout(jwt)
			if err != nil {
				t.Errorf("Logout failed (email: %s - pass: %s). Error: %v", c.email, c.password, err)
			}
			//test second logout on same JWT
			err = controller.Logout(jwt)
			if err == nil {
				t.Errorf("Logout succeded when it should have failed (email: %s - pass: %s).", c.email, c.password)
			}
		}
	}

}
