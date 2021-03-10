package dbmanager

import (
	"testing"
)

var PSQL DBManager
var PSQLerr error

func init(){
	PSQL, PSQLerr = New("postgres","myPassword","iot_dashboard")
}


func TestNewDBManager(t *testing.T){
	//this inherently also tests ConnectToPSQL()
	if PSQLerr != nil {
		t.Errorf("Not Able to connect to database")
	}
}

func TestVerifyUserExists(t *testing.T) {
	cases := []struct {
		email string; exists bool
	}{
		{"user@gmail.com", true},
		{"1@gmail.com", false},
		{"", false},
	}

	for _, c := range cases {
		err := PSQL.VerifyUserExists(c.email)
		if (err == nil && !c.exists) || (err != nil && c.exists) {
			t.Errorf("user %s exists is returning the wrong result", c.email)
		}
	}
}

func TestAddNewUser(t *testing.T) {
	// this tests that we cann't add the same email twice
	cases := []struct {
		email, pass string; exists bool
	}{
		{"user@gmail.com","newpass", true},
		{"101@gmail.com","greatpass", false},
	}

	for _, c := range cases {
		err := PSQL.AddNewUser(c.email, c.pass)
		if (err == nil && c.exists) || (err != nil && !c.exists) {
			t.Errorf("Add new user not working properly for: (email: %s - pass: %s)", c.email, c.pass)
		}
	}
}

