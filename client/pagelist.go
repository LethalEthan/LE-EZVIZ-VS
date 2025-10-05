package client

import (
	"bytes"
	"encoding/json"
	"io"
	"le-ezviz-vs/api"
	"strconv"

	"go.uber.org/zap"
)

type PageListResponse struct {
	Meta struct {
		Code     *int               `json:"code"`
		Message  *string            `json:"message"`
		Moreinfo *map[string]string `json:"moreInfo"`
	} `json:"meta"`
	Page struct {
		Offset       *int  `json:"offset"`
		Limit        *int  `json:"limit"`
		TotalResults *int  `json:"totalResults"`
		HasNext      *bool `json:"hasNext"`
	}
	ResourceInfos *[]Resource            `json:"resourceInfos"`
	VTM           map[string]VTMResource `json:"VTM"`
	DeviceInfos   *[]DeviceInfos         `json:"deviceInfos"`
}

type VTMResource struct {
	Domain          string `json:"domain"`
	ExternalIP      string `json:"externalIp"`
	InternalIP      string `json:"internalIp"`
	Port            int    `json:"port"`
	Memo            string `json:"memo"`
	ForceStreamType int    `json:"forceStreamType"`
	IsBackup        int    `json:"isBackup"`
	PublicKey       struct {
		Key     string `json:"key"`
		Version int    `json:"version"`
	} `json:"publicKey"`
}

type Resource struct {
	ResourceID        string `json:"resourceId"`
	ResourceName      string `json:"resourceName"`
	DeviceSerial      string `json:"deviceSerial"`
	SuperDeviceSerial string `json:"superDeviceSerial"`
	LocalIndex        string `json:"localIndex"`
	ShareType         int    `json:"shareType"`
	Permission        int    `json:"permission"`
	ResourceType      int    `json:"resourceType"`
	ResourceCover     string `json:"resourceCover"`
	IsShow            int    `json:"isShow"`
	VideoLevel        int    `json:"videoLevel"`
	StreamBizUrl      string `json:"streamBizUrl"`
	GroupId           int    `json:"groupId"`
	CustomSetTag      int    `json:"customSetTag"`
	Conceal           int    `json:"conceal"`
	GlobalState       int    `json:"globalState"`
	Child             bool   `json:"child"`
}

type DeviceInfos struct {
	Name                 string  `json:"name"`
	DeviceSerial         string  `json:"deviceSerial"`
	FullSerial           string  `json:"fullSerial"`
	DeviceType           string  `json:"deviceType"`
	DevicePicPrefix      string  `json:"devicePicPrefix"`
	Version              string  `json:"version"`
	SupportExt           string  `json:"supportExt"`
	Status               int     `json:"status"`
	UserDeviceCreateTime string  `json:"userDeviceCreateTime"`
	ChannelNumber        int     `json:"channelNumber"`
	Hik                  bool    `json:"hik"`
	DeviceCategory       string  `json:"deviceCategory"`
	DeviceSubCategory    string  `json:"deviceSubCategory"`
	EZDeviceCapability   string  `json:"ezDeviceCapability"`
	CustomType           string  `json:"customType"`
	OfflineTime          string  `json:"offlineTime"`
	OfflineNotify        int     `json:"offlineNotify"`
	InstructionBook      string  `json:"instructionBook"`
	AuthCode             string  `json:"authCode"`
	UserName             string  `json:"userName"`
	RiskLevel            int     `json:"riskLevel"`
	OfflineTimestamp     int     `json:"offlineTimeStamp"`
	Mac                  string  `json:"mac"`
	ExtStatus            int     `json:"extStatus"`
	Classify             int     `json:"classify"`
	Tags                 *string `json:"tags"`
}

// ToDo - add custom params
func (LEZ *LE_EZVIZ_Client) GetPageList() (*PageListResponse, error) {
	resp, err := LEZ.QueryEncodedAPIRequest("GET", api.V3_USERDEVICES_V1_RESOURCES_PAGELIST, USE_API_URL, map[string]string{"sessionId": *LEZ.LoginResponse.LoginSession.SessionId, "clientType": strconv.Itoa(LEZ.ClientType), "clientNo": LEZ.ClientNo, "clientVersion": "2,5,1,2109068", "groupId": "-1", "limit": "50", "offset": "0", "filter": "VTM"})
	if err != nil {
		log.Error("GetPageList Request error", zap.Error(err))
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading response body", zap.Error(err))
		return nil, err
	}
	PageList := new(PageListResponse)
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(PageList); err != nil {
		log.Error("Error decoding JSON", zap.Error(err))
		return nil, err
	}
	return PageList, nil
}
