package controller

import (
	"iotdashboard/dbmanager"
	"iotdashboard/utils"
)

type ControllerService struct {
	PSQL      dbmanager.DBManager
	TokenUtil utils.TokenUtil
}

func NewController() (ControllerService, error) {
	PSQL, err := dbmanager.New("postgres", "myPassword", "iot_dashboard")
	if err != nil {
		return ControllerService{}, err
	}
	TokenUtil, err := utils.NewTokenUtil()

	return ControllerService{PSQL, TokenUtil}, nil
}

func (ct *ControllerService) Login(email, password string) (string, error) {
	// validate basic auth
	err := ct.PSQL.CheckUserCredentials(email, password)
	if err != nil {
		return "", err
	}

	// create and return JWT
	var token string
	token, err = ct.TokenUtil.CreateJWT(15)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (ct *ControllerService) Logout(token string) error {
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
