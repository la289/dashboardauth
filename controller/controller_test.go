package controller

import(
	"testing"
	"time"
)

func TestGenerateRandomToken(t *testing.T) {
	n := 256
	key, err := GenerateRandomToken(n)
	if (err != nil || len(key) != n){
		t.Errorf("The signing key is wrong: %s", key)
	}
}

func TestGenerateRandomString(t *testing.T) {
	n := 5
	s, err := GenerateRandomString(n)
	if (err != nil || len(s) < n){
		t.Errorf("The signing key is wrong: %s", s)
	}
}


func TestCheckUserCredentials(t *testing.T) {
	cases := []struct {
		email, pass string; exists bool
	}{
		{"user@gmail.com", "S3cure3Pa$$", true},
		{"user@gmail.com", "wrongpass", false},
		{"1000@doesntexist.com", "bloop", false},
	}

	for _,c := range cases{
		err := CheckUserCredentials(c.email,c.pass)
		if (err != nil && c.exists) || (err == nil && !c.exists) {
			t.Errorf("User credential validation failed (email: %s - pass: %s). Error: %v", c.email, c.pass,err)
		}
	}
}

func TestCreateAndValidateJWT(t *testing.T) {
	token, err := CreateJWT(2)
	if err != nil{
		t.Errorf("JWT creation failed: %s", err)
	}
	//test validation of valid token
	_,err = ValidateJWT(token)
	if err != nil{
		t.Errorf("Validation failed when it should have succeeded: %v", err)
	}
	//test validation of garbage token
	_,err = ValidateJWT("please fail")
	if err == nil{
		t.Errorf("Validation succeeded when it should have failed")
	}
	//test expiration of token
	time.Sleep(3*time.Second)
	_,err = ValidateJWT(token)
	if err == nil{
		t.Errorf("Validation succeeded when it should have failed")
	}
}

func TestCreateAndBlocklistJWT(t *testing.T) {
	token, err := CreateJWT(60)
	if err != nil{
		t.Errorf("JWT creation failed: %s", err)
	}
	err = Logout(token)
	if err != nil{
		t.Errorf("Logout failed: %s", err)
	}
	//test validation of logged out token
	_,err = ValidateJWT(token)
	if err == nil{
		t.Errorf("Validation succeeded when it should have failed")
	}
}
