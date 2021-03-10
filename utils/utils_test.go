package utils

import (
	"testing"
)

func TestGeneratePassHash(t *testing.T) {
	cases := []struct {
		pass string
	}{
		{"Hello, world"},
		{"Hello, 世界"},
		{""},
		{"Password!23"},
	}

	for _, c := range cases {
		hp, err := GeneratePassHash(c.pass)
		if (string(hp[0:6]) != "$2a$10" || len(hp) != 60 || err != nil) {
		t.Errorf("Password hashing failed: %s", hp)
		}
    }
}

func TestCheckPasswordPass(t *testing.T) {
	cases := []struct {
		pass string
	}{
		{"Hello, world"},
		{"Hello, 世界"},
		{""},
		{"Password!23"},
	}

	for _, c := range cases {
		hp, err := GeneratePassHash(c.pass)
		if err != nil {
			t.Errorf("Generating password hash failed.")
		}
		err = CheckPassword(hp, c.pass)
		if err != nil {
			t.Errorf("Password Comparison failed.")
		}
	}
}

func TestCheckPasswordFail(t *testing.T) {
	cases := []struct {
		pass1, pass2 string
	}{
		{"Hello, world", "Hello  world"},
		{"Hello, 世界", "hello, 世界"},
		{"", " "},
		{"Password!23", "Password123"},
	}

	for _, c := range cases {
		hp,err := GeneratePassHash(c.pass1)
		if err != nil {
			t.Errorf("Generating password hash failed.")
		}

		err = CheckPassword(hp, c.pass2)
		if err == nil {
			t.Errorf("Password Comparison Succeeded but should have failed.")
		}
	}
}
