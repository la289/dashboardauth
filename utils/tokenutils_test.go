package utils

import (
	"testing"
	"time"
)

func TestGenerateRandomToken(t *testing.T) {
	n := 256
	key, err := GenerateRandomToken(n)
	if err != nil || len(key) != n {
		t.Errorf("The signing key is wrong: %s", key)
	}
}

func TestGenerateRandomString(t *testing.T) {
	tu, err := NewTokenUtil()
	if err != nil {
		t.Fatalf("Not able to create TokenUtil: %v \n", err)
	}
	n := 5
	s, err := tu.GenerateRandomString(n)
	if err != nil || len(s) != n {
		t.Errorf("The signing key is wrong: %s", s)
	}
}

func TestCreateAndGetJWTExpiry(t *testing.T) {
	tu, err := NewTokenUtil()
	if err != nil {
		t.Fatalf("Not able to create TokenUtil: %v \n", err)
	}
	token, err := tu.CreateJWT(time.Second * 1)
	if err != nil {
		t.Errorf("JWT creation failed: %s", err)
	}
	//test validation of valid token
	_, err = tu.GetJWTExpiry(token)
	if err != nil {
		t.Errorf("Validation failed when it should have succeeded: %v", err)
	}
	//test validation of garbage token
	_, err = tu.GetJWTExpiry("please fail")
	if err == nil {
		t.Errorf("Validation succeeded when it should have failed")
	}
	//test expiration of token
	//TODO: Replace time.sleep with https://github.com/jonboulle/clockwork
	time.Sleep(3 * time.Second)
	_, err = tu.GetJWTExpiry(token)
	if err == nil {
		t.Errorf("Validation succeeded when it should have failed")
	}
}

func TestCreateAndBlocklistJWT(t *testing.T) {
	tu, err := NewTokenUtil()
	if err != nil {
		t.Fatalf("Not able to create TokenUtil: %v \n", err)
	}
	token, err := tu.CreateJWT(time.Second * 60)
	if err != nil {
		t.Errorf("JWT creation failed: %s", err)
	}

	tu.BlockListToken(token, time.Now().UTC().Add(time.Second*60))

	//test validation of logged out token
	_, err = tu.GetJWTExpiry(token)
	if err == nil {
		t.Errorf("Validation succeeded when it should have failed")
	}
}
