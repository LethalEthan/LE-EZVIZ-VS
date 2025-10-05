package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"le-ezviz-vs/api"
	"strconv"

	"go.uber.org/zap"
)

type VTDU_TokenV2 struct {
	Message *string   `json:"msg"`
	Tokens  *[]string `json:"tokens"`
	Retcode *int      `json:"retcode"`
}

func (LEZ *LE_EZVIZ_Client) GetVTDUv2Token() (*VTDU_TokenV2, error) {
	claims, err := LEZ.GetSessionIDClaims()
	if err != nil {
		log.Error("Error getting VTDU TokenV2 sign claim", zap.Error(err))
		return nil, err
	}

	resp, err := LEZ.QueryEncodedAPIRequest("GET", api.VTDUTOKEN2, USE_AUTH_URL, map[string]string{"ssid": *LEZ.LoginResponse.LoginSession.SessionId, "sign": claims["s"].(string)})
	if err != nil {
		log.Error("Error Get VTDU TokenV2", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading response body", zap.Error(err))
		return nil, err
	}
	LEZ.VTDUTokens = new(VTDU_TokenV2)
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(LEZ.VTDUTokens); err != nil {
		log.Error("Error decoding JSON", zap.Error(err))
		return nil, err
	}
	if *LEZ.VTDUTokens.Retcode != 0 {
		log.Error("VTDU TokenV2 non 0 return", zap.Int("Retcode", *LEZ.VTDUTokens.Retcode))
		return nil, errors.New("non 0 return code: " + strconv.Itoa(*LEZ.VTDUTokens.Retcode))
	}
	return LEZ.VTDUTokens, nil
}
