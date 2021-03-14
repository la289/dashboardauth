package dbmanager_test

import (
	"testing"
	"iotdashboard/dbmanager"
)

var PSQL dbmanager.DBManager
var PSQLerr error

func init() {
	PSQL, PSQLerr = dbmanager.New("postgres", "myPassword", "iot_dashboard")
}

func TestNewDBManager(t *testing.T) {
	//this inherently also tests ConnectToPSQL()
	if PSQLerr != nil {
		t.Errorf("Not Able to connect to database")
	}
}

// TestCheckUserCredentials inherently also tests getPasswordHash
func TestCheckUserCredentials(t *testing.T) {
	cases := []struct {
		email, pass string
		exists      bool
	}{
		{"user@gmail.com", "S3cure3Pa$$", true},
		{"user@gmail.com", "wrongpass", false},
		{"1000@doesntexist.com", "bloop", false},
	}

	for _, c := range cases {
		err := PSQL.CheckUserCredentials(c.email, c.pass)
		if (err != nil && c.exists) || (err == nil && !c.exists) {
			t.Errorf("User credential validation failed (email: %s - pass: %s). Error: %v", c.email, c.pass, err)
		}
	}
}

func TestAddNewUser(t *testing.T) {
	// this tests that we cann't add the same email twice
	cases := []struct {
		email, pass string
		exists      bool
	}{
		{"user@gmail.com", "newpass", true},
		{"108@gmail.com", "greatpass", false},
	}

	for _, c := range cases {
		err := PSQL.AddNewUser(c.email, c.pass)
		if (err == nil && c.exists) || (err != nil && !c.exists) {
			t.Errorf("Add new user fails for: (email: %s - pass: %s). Error: %v", c.email, c.pass, err)
		}
	}
}
