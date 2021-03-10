package utils

import(
	"golang.org/x/crypto/bcrypt"
)

func GeneratePassHash(password string) ([]byte, error) {
	// could vaildate password requirements, but since we are not
	// implementing a create account, leaving it out for now
	return bcrypt.GenerateFromPassword([]byte(password), 10)
}

func CheckPassword(hashedpass []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedpass, []byte(password))
}
