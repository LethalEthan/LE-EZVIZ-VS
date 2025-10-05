package client

import (
	"bytes"
	"encoding/json"
	"io"
	"le-ezviz-vs/api"
	"strconv"

	"go.uber.org/zap"
)

type ServerInfoGetResponse struct {
	ResultCode string `json:"resultCode"`
	ServerResp *struct {
		Stun1IP        *string `json:"stun1Addr"`
		Stun1Port      *int    `json:"stun1Port"`
		Stun2IP        *string `json:"stun2Addr"`
		Stun2Port      *int    `json:"stun2Port"`
		TtsIP          *string `json:"ttsAddr"`
		TtsPort        *int    `json:"ttsPort"`
		VtmIP          *string `json:"vtmAddr"`
		VtmPort        *int    `json:"vtmPort"`
		AuthAddr       *string `json:"authAddr"`
		PushAddr       *string `json:"pushAddr"`
		PushHttpPort   *int    `json:"pushHttpPort"`
		PushHttpsPort  *int    `json:"pushHttpsPort"`
		CloutManager   *string `json:"cloutManager"` // manage your clout
		NodeJsAddr     *string `json:"nodeJsAddr"`
		NodeJsHttpPort *int    `json:"nodeJsHttpPort"`
		PmsAddr        *string `json:"pmsAddr"`
		PmsPort        *int    `json:"pmsPort"`
		DcLogAddr      *string `json:"cdLogAddr"`
		DcLogPort      *int    `json:"dcLogPort"`
	} `json:"serverResp"`
}

func (LEZ *LE_EZVIZ_Client) GetServerInfo() (*ServerInfoGetResponse, error) {
	resp, err := LEZ.FormURLEncodedAPIRequest("POST", api.API_SERVER_INFO_GET, USE_API_URL, map[string]string{"sessionId": *LEZ.LoginResponse.LoginSession.SessionId, "clientType": strconv.Itoa(LEZ.ClientType)})
	if err != nil {
		log.Error("GetServerInfo Request error", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading response body", zap.Error(err))
		return nil, err
	}
	LEZ.APIServerInfo = new(ServerInfoGetResponse)
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(LEZ.APIServerInfo); err != nil {
		log.Error("Error decoding JSON", zap.Error(err))
		LEZ.APIServerInfo = nil
		return nil, err
	}
	LEZ.AUTH_URL = *LEZ.APIServerInfo.ServerResp.AuthAddr
	return LEZ.APIServerInfo, nil
}
