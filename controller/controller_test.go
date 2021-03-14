package controller

import (
	"testing"
)


func TestLoginAndLogout(t *testing.T) {
	controller, err := NewController()
	if err != nil {
		t.Errorf("Failed to initialize new controller: %v \n", err)
	}
	cases := []struct {
		email, pass string
		success      bool
	}{
		{"user@gmail.com", "S3cure3Pa$$", true},
		{"user@gmail.com", "wrongpass", false},
		{"1000@doesntexist.com", "bloop", false},
	}

	for _, c := range cases {
		jwt, err := controller.Login(c.email, c.pass)
		if (err != nil && c.success) || (err == nil && !c.success) {
			t.Errorf("Login failed (email: %s - pass: %s). Error: %v", c.email, c.pass, err)
		}

	// if successful login, test logout
	if c.success == true {
		err = controller.Logout(jwt)
		if err != nil {
			t.Errorf("Logout failed (email: %s - pass: %s). Error: %v", c.email, c.pass, err)
		}
		//test second logout on same JWT
		err = controller.Logout(jwt)
		if err == nil {
			t.Errorf("Logout succeded when it should have failed (email: %s - pass: %s).", c.email, c.pass)
		}
	}
}



}

