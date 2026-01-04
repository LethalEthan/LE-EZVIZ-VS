package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	client "le-ezviz-vs/client"
	logging "le-ezviz-vs/logging"

	"go.uber.org/zap"
)

var email = flag.String("email", "", "EZVIZ user e-mail")
var password = flag.String("password", "", "EZVIZ user password")
var region = flag.String("region", "Europe", "Manually set the region you are in: Europe, Africa, India, Oceania, NorthAmercia, Russia, SouthAmerica")
var terminalName = flag.String("terminalName", "LE-EZ", "Optional: Set the name of what your device is called when viewing registered devices (terminals)")
var preservefc = flag.Bool("preserveFeatureCode", true, "Preserve the featurecode, it is essentially a random ID to identify the terminal")
var deviceSerial = flag.String("deviceSerial", "", "The device serial you want to connect to")
var preserveSession = flag.Bool("preserveSession", false, "Preserve session token") //TBD
var stdout = flag.Bool("stdout", true, "Print log to the terminal")
var logFile = flag.Bool("logFile", true, "Print log to lez.log")
var log *zap.Logger

func main() {
	flag.Parse()
	if *email == "" {
		panic("email empty")
	}
	if *password == "" {
		panic("password empty")
	}
	logging.CreateLogger(*logFile, *stdout)
	client.CurrentRegion = *region
	client.TerminalName = *terminalName
	log = logging.Log
	if _, ok := client.Regions[*region]; !ok {
		log.Error("Invalid region", zap.String("Valid Values", "Europe|Africa|India|Oceania|NorthAmerica|Russia|SouthAmerica"))
	}
	LEZ, err := client.NewLE_EZVIZ_Client(*email, *password, *region, "00000000000000000000000000000000", *terminalName, "shipin7", 15)
	if err != nil {
		panic(err)
	}
	if *preservefc {
		LEZ.LoadFeatureCode("featurecode")
	} else {
		LEZ.SetFeatureCode(hex.EncodeToString(client.GenerateFeatureCode()))
	}
	if _, err = LEZ.V3_Login(); err != nil {
		panic(err)
	}
	_, err = LEZ.GetServerInfo()
	if err != nil {
		panic(err)
	}

	PageList, err := LEZ.GetPageList()
	if err != nil {
		panic(err)
	}
	_, err = LEZ.GetVTDUv2Token()
	if err != nil {
		panic(err)
	}
	for _, v := range *PageList.DeviceInfos {
		log.Info(v.Name, zap.String("Serial", v.DeviceSerial))
	}
	// fmt.Println(PageList.VTM)
	fmt.Println("!!!WARNING: This library is in beta, only use for development/testing until it is stable, things will change as development continues!!!")
	fmt.Println("!!!Encryption is not yet available including E2EE with stream servers, your streams will be unencrypted until encryption is implemented!!!")
	if *deviceSerial != "" {
		var RI client.Resource
		var DI client.DeviceInfos
		for _, v := range *PageList.ResourceInfos {
			if *deviceSerial == v.DeviceSerial {
				RI = v
			}
		}
		if v, ok := PageList.VTM[RI.ResourceID]; ok {
			VS, err := LEZ.ConnectVTM(v.ExternalIP, v.Port, RI, DI, "", v.PublicKey.Key)
			if err != nil {
				log.Error("Error connecting VTM", zap.Error(err))
				return
			}
			Tokens := *LEZ.VTDUTokens.Tokens
			URL := LEZ.BuildVtmUrl(VS.VTMIP, VS.VTMPort, RI.DeviceSerial, RI.StreamBizUrl, Tokens[0], DI.ChannelNumber, LEZ.ClientType)
			RStreamInfoRsp, err := LEZ.StartVTMStream(VS, URL)
			if err != nil {
				panic(err)
			}

			IP, Port, _, _, err := LEZ.ParseVtmUrl(*RStreamInfoRsp.Streamurl)
			if err != nil {
				panic(err)
			}
			VTDUStream, err := LEZ.ConnectVTDU(IP, Port, *RStreamInfoRsp.Vtmstreamkey, v.PublicKey.Key)
			if err != nil {
				panic(err)
			}
			err = LEZ.StartVTDUStream(VTDUStream, URL)
			if err != nil {
				panic(err)
			}
		}
	}

}
