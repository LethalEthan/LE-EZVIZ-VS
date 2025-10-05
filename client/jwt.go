package client

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (LEZ *LE_EZVIZ_Client) GetSessionIDClaims() (jwt.MapClaims, error) {
	if LEZ.LoginResponse.LoginSession.SessionId == nil {
		log.Error("You need to login first before getting sessionID claims")
		return nil, errors.New("you need to login first")
	}
	claims, err := parseNoVerify(*LEZ.LoginResponse.LoginSession.SessionId)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func parseNoVerify(tokenString string) (jwt.MapClaims, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid claims type")
}

// CheckSessionID - Check if the sessionID is still valid or exists
func (LEZ *LE_EZVIZ_Client) CheckSessionID() bool {
	if LEZ.LoginResponse.LoginSession.SessionId == nil {
		log.Error("You need to login first before getting sessionID claims")
		return false
	}
	claims, err := parseNoVerify(*LEZ.LoginResponse.LoginSession.SessionId)
	if err != nil {
		return false
	}
	if v, ok := claims["exp"]; ok {
		if time.Now().After(time.Unix(v.(int64), 0)) {
			return false
		} else {
			return true
		}
	} else {
		return false
	}
}
