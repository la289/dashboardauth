package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	NotBefore int64 `json:"nbf"`
	jwt.StandardClaims
}

type TokenUtil struct {
	jwtKey []byte
	//Blocklist key is the JWT and the value is the expiration epoch time
	blocklist *sync.Map
}

var ErrExpiredToken = errors.New("Token is expired")

func NewTokenUtil() (TokenUtil, error) {
	jwtKey, err := GenerateRandomToken(256)
	if err != nil {
		return TokenUtil{}, err
	}
	return TokenUtil{jwtKey, new(sync.Map)}, nil
}

func (tu *TokenUtil) CreateJWT(validPeriod time.Duration) (string, error) {
	claims := Claims{
		NotBefore: time.Now().UTC().Unix(),
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: time.Now().UTC().Add(time.Second * validPeriod).Unix(),
			Issuer:    "iot-dash",
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(tu.jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (tu *TokenUtil) GetJWTExpiry(rawToken string) (time.Time, error) {
	exp, ok := tu.blocklist.Load(rawToken)
	if ok {
		return exp.(time.Time), ErrExpiredToken
	}
	token, err := jwt.ParseWithClaims(
		rawToken,
		&Claims{},
		func(rawToken *jwt.Token) (interface{}, error) {
			return tu.jwtKey, nil
		})
	if err != nil {
		return time.Time{}, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return time.Time{}, errors.New("Couldn't Parse Token Claims")
	}

	now := time.Now().UTC().Unix()
	expiresAt := time.Unix(claims.ExpiresAt, 0)
	if claims.ExpiresAt < now || now < claims.NotBefore {
		return expiresAt, ErrExpiredToken
	}

	return expiresAt, nil
}

func (tu *TokenUtil) BlockListToken(jwt string, expiration time.Time) {
	tu.blocklist.Store(jwt, expiration)
}

func (tu *TokenUtil) GenerateRandomString(n int) (string, error) {
	s, err := GenerateRandomToken(n)
	return base64.URLEncoding.EncodeToString(s)[:n], err
}

func GenerateRandomToken(n int) ([]byte, error) {
	key := make([]byte, n)
	_, err := rand.Read(key)
	return key, err
}
