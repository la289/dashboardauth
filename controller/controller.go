package controller


import(
	"iotdashboard/dbmanager"
	"github.com/dgrijalva/jwt-go"
	"crypto/rand"
	"time"
	"iotdashboard/utils"
	"errors"
	"encoding/base64"
)

//TODO: these probably shouldn't be global vars
var PSQL dbmanager.DBManager
var jwtKey []byte

// key is the JWT, value is the expiration epoch time
var blocklist map[string]int64

type Claims struct {
	NotBefore int64 `json:"nbf"`
	jwt.StandardClaims
}

////////Errors
var ErrExpiredToken = errors.New("Token is expired")

func init() {
	var err error
	jwtKey, err = GenerateRandomToken(256)
	if err != nil {
		panic(err)
	}
	PSQL, err = dbmanager.New("postgres","myPassword","iot_dashboard")
	if err != nil {
		panic(err)
	}

	blocklist = make(map[string]int64)
}


func Login(email, password string) (string, error) {
	// use user struct instead of email and password?

	// validate CSRF

	// validate basic auth
	err := CheckUserCredentials(email, password)
	if err != nil {
		return "", err
	}

	// create and return JWT
	var token string
	token, err = CreateJWT(60)
	if err != nil {
		return "", err
	}
	return token, nil

}

func Logout(token string) error{
	// validate CSRF

	// validate JWT
	exp, err := ValidateJWT(token)
	if err != nil{
		return err
	}
	// blocklist JWT
	blocklist[token] = exp
	return nil

}



// returns an Error if the credentials are not valid
func CheckUserCredentials(email, password string) error {
	hash, err := PSQL.GetUserPassHash(email)
	if err != nil {
		return err
	}
	return utils.CheckPassword(hash, password)
}


func CreateJWT(validPeriod int64) (string, error) {
		claims := Claims{
			NotBefore: time.Now().Unix(),
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: time.Now().Add(time.Second * time.Duration(validPeriod)).Unix(),
				Issuer: "iot-dash",
			},
		}

		// Create the token
		token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"),claims)
		// token.Claims["exp"] =
		// Sign and get the complete encoded token as a string
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			return "", err
		}
		return tokenString, nil
}




func ValidateJWT(rawToken string) (int64, error) {
	exp, ok := blocklist[rawToken]
	if ok {
		return exp, ErrExpiredToken
	}
	token, err := jwt.ParseWithClaims(
		rawToken,
		&Claims{},
		func(rawToken *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
	if err != nil{
		return 0, err
	}
	// TODO: make this a subfunction so that validateJWT returns just the err
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, errors.New("Couldn't Parse Token Claims")
	}
	now := time.Now().UTC().Unix()
	if (claims.ExpiresAt < now || now < claims.NotBefore) {
		return claims.ExpiresAt, ErrExpiredToken
	}

	return claims.ExpiresAt, nil
}



func GenerateRandomToken(n int) ([]byte, error) {
	key := make([]byte, n)
    _, err := rand.Read(key)
    if err != nil {
        return nil, err
    }
    // fmt.Println(key)
	return key, nil
}

func GenerateRandomString(n int) (string, error) {
	s, err := GenerateRandomToken(n)
	return base64.URLEncoding.EncodeToString(s), err
}
