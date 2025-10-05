package client

type V3_Auth_Login_Response struct {
	HcGvIsolate *bool `json:"hcGvIsolate"`
	Isolate     *bool `json:"isolate"`
	LoginArea   struct {
		ApiDomain *string `json:"apiDomain"`
		AreaId    *int    `json:"areaId"`
		AreaName  *string `json:"areaName"`
		WebDomain *string `json:"webDomain"`
	} `json:"loginArea"`
	LoginSession struct {
		RfSessionId *string `json:"rfSessionId"`
		SessionId   *string `json:"sessionId"`
	} `json:"loginSession"`
	LoginTerminalStatus struct {
		TerminalBinded *string `json:"terminalBinded"`
		TerminalOpened *string `json:"terminalOpened"`
	} `json:"loginTerminalStatus"`
	LoginUser struct {
		AreaId                         *int    `json:"areaId"`
		AvatarPath                     *string `json:"avatarPath"`
		Category                       *int    `json:"category"`
		ConfusedEmail                  *string `json:"confusedEmail"`
		ConfusedUsername               *string `json:"confusedUsername"`
		ConfusedPhone                  *string `json:"confusedPhone"`
		Contact                        *string `json:"contact"`
		Customno                       *string `json:"customno"`
		Email                          *string `json:"email"`
		HomeTitle                      *string `json:"homeTitle"`
		IsSecurityBind                 *bool   `json:"isSecurityBind"`
		LangType                       *string `json:"langType"`
		Location                       *string `json:"location"`
		MsgStatus                      *int    `json:"msgStatus"`
		NeedTrans                      *bool   `json:"needTrans"`
		Phone                          *string `json:"phone"`
		RegDate                        *string `json:"regDate"`
		TransferringToStandaloneRegion *bool   `json:"transferringToStandaloneRegion"`
		UserCode                       *string `json:"userCode"`
		UserId                         *string `json:"userId"`
		Username                       *string `json:"username"`
	} `json:"loginUser"`
	Meta struct {
		Code     *int               `json:"code"`
		Message  *string            `json:"message"`
		MoreInfo *map[string]string `json:"moreInfo"`
	} `json:"meta"`
	TelephoneCode *string `json:"telephoneCode"`
}
