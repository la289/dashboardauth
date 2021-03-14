package dbmanager

import (
	"testing"
)


// TestCheckUserCredentials inherently also tests getPasswordHash
func TestCheckUserCredentials(t *testing.T) {
	PSQL, err := New("postgres", "myPassword", "iot_dashboard")
	if err != nil{
		t.Errorf("Unable to initialize DB Manager: %v \n", err)
	}

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
	PSQL, err := New("postgres", "myPassword", "iot_dashboard")
	if err != nil{
		t.Errorf("Unable to initialize DB Manager: %v \n", err)
	}

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
