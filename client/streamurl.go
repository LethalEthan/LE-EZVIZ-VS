package client

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func (LEZ *LE_EZVIZ_Client) BuildVtmUrl(VTMip string, VTMport int, DeviceSerial, Biz, VTDUToken string, Channel, ClientType int) string {
	URL := fmt.Sprintf("ysproto://%s:%d/live?", VTMip, VTMport)
	// Yes we should use smth like url.Value{} but it automatically sorts the keys which messes the layout
	Params := fmt.Sprintf("dev=%s&chn=%d&stream=1&cln=%d&isp=0&auth=1&ssn=%s&%s&vip=0&timestamp=%d", DeviceSerial, Channel, ClientType, VTDUToken, Biz, time.Now().UnixMilli())
	return URL + Params
}

func (LEZ *LE_EZVIZ_Client) ParseVtmUrl(URL string) (string, int, string, map[string]string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		log.Error("Error parsing vtmurl", zap.Error(err))
		return "", 0, "", nil, err
	}
	if u.Scheme != "ysproto" {
		return "", 0, "", nil, errors.New("invalid scheme")
	}
	if u.Host == "" {
		return "", 0, "", nil, errors.New("host empty")
	}
	// fmt.Println("Scheme :", u.Scheme) // ex: ysproto
	// fmt.Println("Host   :", u.Host)   // ex: 100.100.100.100:8554
	// fmt.Println("Path   :", u.Path)   // ez: /live
	Params := make(map[string]string)
	for key, vals := range u.Query() {
		if len(vals) > 0 {
			Params[key] = vals[0]
		}
	}
	ip, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		log.Error("Error splitting host", zap.Error(err))
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.Error("Error converting port to int", zap.Error(err))
		return "", 0, "", nil, err
	}

	// fmt.Println("\nQuery options:")
	// for k, v := range Params {
	// 	fmt.Printf("%s=%s\n", k, v)
	// }
	return ip, portInt, u.Path, Params, nil
}

/* stream params
&serial=
&streamtag=
&chn= - Channel
&cln= - clientType
&isp=
&auth=
&ssn= - session auth0 opaque token
&weakstream=1 - I have seen some streams use it mainly when encrypted, is it signifying a weakly secured stream or something else, I am leaning towards it being something else
&isretry=
&lid
&e2ee=   - set to 1 when using encrypted channel
// the are used when using /playback path I believe this is when playing back cloud videos or perhaps even from the device, not 100% known
&begin=
&end=
&seg=
*/
//&a=1 - ? audio? another auth flag?

// There is also ysudp:// currently not looking into that atm but it does have &linkid= we'll look into UDP streams later
// It does make more sense to use UDP for streaming but most of the time it uses TCP, I am assuming to prevent dropouts when malformed packets come in or maybe issues relating to encyption>
