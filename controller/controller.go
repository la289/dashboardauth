package controller

import (
	"iotdashboard/dbmanager"
	"iotdashboard/utils"
	"time"
)

type ControllerService struct {
	PSQL      *dbmanager.DBManager
	TokenUtil *utils.TokenUtil
}

func NewController() (*ControllerService, error) {
	psql, err := dbmanager.New("postgres", "postgres", "iot_dashboard")
	if err != nil {
		return &ControllerService{}, err
	}
	tokenUtil, err := utils.NewTokenUtil()
	if err != nil {
		return &ControllerService{}, err
	}

	return &ControllerService{psql, tokenUtil}, nil
}

func (ct *ControllerService) Login(email, password string) (string, error) {
	// validate basic auth
	err := ct.PSQL.CheckUserCredentials(email, password)
	if err != nil {
		return "", err
	}

	// create and return JWT
	token, err := ct.TokenUtil.CreateJWT(time.Second * 60)
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
