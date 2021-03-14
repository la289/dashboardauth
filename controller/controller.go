package controller

import (
	"iotdashboard/dbmanager"
	"iotdashboard/utils"
)

type Controller struct {
	psql   dbmanager.DBManager
	TokenUtil utils.TokenUtil
  }

func NewController() (Controller, error) {
	PSQL, err := dbmanager.New("postgres", "myPassword", "iot_dashboard")
	if err != nil {
		return Controller{}, err
	}
	TokenUtil, err := utils.NewTokenUtil()

	return Controller{PSQL, TokenUtil}, nil
}


func (ct *Controller) Login(email, password string) (string, error) {
	// validate basic auth
	err := ct.psql.CheckUserCredentials(email, password)
	if err != nil {
		return "", err
	}

	// create and return JWT
	var token string
	token, err = ct.TokenUtil.CreateJWT(60)
	if err != nil {
		return "", err
	}
	return token, nil

}

func (ct *Controller) Logout(token string) error {
	// validate CSRF

	// validate JWT
	exp, err := ct.TokenUtil.GetJWTExpiry(token)
	if err != nil {
		return err
	}
	// blocklist JWT
	ct.TokenUtil.BlockListToken(token, exp)
	return nil

}

