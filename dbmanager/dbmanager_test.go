package dbmanager

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestCheckUserCredentials inherently also tests getPasswordHash
func TestCheckUserCredentials(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v \n", err)
	}
	defer db.Close()

	PSQL, err := New("postgres", "myPassword", "iot_dashboard")
	if err != nil {
		t.Fatalf("Unable to initialize DB Manager: %v \n", err)
	}
	//overwrite db connection with mock
	PSQL.DB = db

	cases := []struct {
		email, password, hashedPassword string
		exists                          bool
	}{
		{"user@gmail.com", "S3cure3Pa$$", "$2a$10$cRmL5Rtm0bunl1uqYAP.8OfJE36RUkvMcX3.v0kJyY2JBhalX4KEG", true},
		{"user@gmail.com", "wrongpass", "$2a$10$cRmL5Rtm0bunl1uqYAP.8OfJE36RUkvMcX3.v0kJyY2JBhalX4KEG", false},
		{"1000@doesntexist.com", "bloop", "notarealhash", false},
	}

	for _, c := range cases {
		rows := sqlmock.NewRows([]string{"password"}).
			AddRow(c.hashedPassword)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT password from users WHERE email = $1")).WillReturnRows(rows)

		err := PSQL.CheckUserCredentials(c.email, c.password)
		if (err != nil && c.exists) || (err == nil && !c.exists) {
			t.Errorf("User credential validation failed (email: %s - pass: %s). Error: %v", c.email, c.password, err)
		}
	}
}

func TestAddNewUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v \n", err)
	}
	defer db.Close()

	PSQL, err := New("postgres", "myPassword", "iot_dashboard")
	if err != nil {
		t.Fatalf("Unable to initialize DB Manager: %v \n", err)
	}
	//overwrite db connection with mock
	PSQL.DB = db

	cases := []struct {
		email, pass, mockResponse string
		shouldSucceed             bool
	}{
		{"user@gmail.com", "newpass", "INSERT 0 1", true},
		{"user@gmail.com", "greatpass", `ERROR:  duplicate key value violates unique constraint "users_email_key"`, false},
	}

	for _, c := range cases {
		if c.shouldSucceed {
			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users(email,password) VALUES ($1 , $2);")).
				WillReturnResult(sqlmock.NewResult(1, 1)) //result not important since we only check for error
		} else {
			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users(email,password) VALUES ($1 , $2);")).
				WillReturnError(fmt.Errorf(c.mockResponse))
		}

		err := PSQL.AddNewUser(c.email, c.pass)
		if (err != nil && c.shouldSucceed) || (err == nil && !c.shouldSucceed) {
			t.Errorf("Add new user fails for: (email: %s - pass: %s). Error: %v", c.email, c.pass, err)
		}
	}
}
