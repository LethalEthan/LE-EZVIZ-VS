package client

import (
	"le-ezviz-vs/ezproto"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func CreateStreamInfoReq(URL, VTMStreamKey string) (*[]byte, error) {
	SIR := new(ezproto.StreamInfoReq)
	SIR.Streamurl = &URL
	if VTMStreamKey != "" {
		SIR.Vtmstreamkey = &VTMStreamKey
	}
	version := "v3.6.3.20221124"
	var pt int32 = 0
	SIR.Proxytype = &pt
	SIR.Clnversion = &version
	SIR.Useragent = &version
	ba, err := proto.Marshal(SIR)
	if err != nil {
		log.Error("Error StreamInfoReq Marshal", zap.Error(err))
		return nil, err
	}
	return &ba, nil
}

func ParseStreamInfoRsp(Data []byte) (*ezproto.StreamInfoRsp, error) {
	StreamInfoRsp := new(ezproto.StreamInfoRsp)
	err := proto.Unmarshal(Data, StreamInfoRsp)
	if err != nil {
		log.Error("Error unmarshaling StreamInfoRsp", zap.Error(err))
		return nil, err
	}
	return StreamInfoRsp, nil
}

func CreateStreamKeepaliveReq(StreamSSN string) (*[]byte, error) {
	KeepAliveReq := new(ezproto.StreamKeepAliveReq)
	KeepAliveReq.Streamssn = []byte(StreamSSN)
	ba, err := proto.Marshal(KeepAliveReq)
	if err != nil {
		log.Error("Error StreamInfoReq Marshal", zap.Error(err))
		return nil, err
	}
	return &ba, nil
}
