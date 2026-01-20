package client

import (
	"errors"
	"le-ezviz-vs/ezproto"
	"net"
	"strconv"

	"go.uber.org/zap"
)

type VTMStream struct {
	Conn           net.Conn
	VTMIP          string
	VTMPort        int
	VTMPublicKey   string
	DevicePassword string
}

func (LEZ *LE_EZVIZ_Client) ConnectVTM(vtmIP string, vtmPort int, deviceResource Resource, deviceInfos DeviceInfos, devicePwd, vtmPublicKey string) (*VTMStream, error) {
	sock, err := net.Dial("tcp", vtmIP+":"+strconv.Itoa(vtmPort))
	if err != nil {
		log.Error("Error dialing VTM", zap.Error(err))
		return nil, err
	}
	VS := &VTMStream{Conn: sock, VTMIP: vtmIP, VTMPort: vtmPort, VTMPublicKey: vtmPublicKey}
	return VS, nil
}

/* breakdown:
We send StreamInfoReq to VTM
We parse StreamInfoRsp from VTM
Then use the information to create VTDUStream
*/

func (LEZ *LE_EZVIZ_Client) StartVTMStream(VS *VTMStream, StreamURL string) (*ezproto.StreamInfoRsp, error) {
	StreamReq, err := CreateStreamInfoReq(StreamURL, "")
	if err != nil {
		log.Error("Error creating StreamInfoReq", zap.Error(err))
		return nil, err
	}
	EncodedPacket := EncodeVTMPacket(*StreamReq, CHAN_MSG, MSG_STREAMINFO_REQ)
	log.Sugar().Debugf("bytes: %x", EncodedPacket)
	// return nil
	_, err = VS.Conn.Write(EncodedPacket)
	if err != nil {
		log.Error("Error writing to TCP sock", zap.Error(err))
		return nil, err
	}
	Packet := new(VTMPacket)
	Packet.Header = make([]byte, 8)
	if _, err = VS.Conn.Read(Packet.Header); err != nil {
		log.Error("Error reading from TCP sock", zap.Error(err))
		return nil, err
	}
	Len, _, _, Msg, err := Packet.DecodeHeader()
	if err != nil {
		log.Error("Error decoding packet header", zap.Error(err))
		return nil, err
	}
	if Msg != MSG_STREAMINFO_RSP {
		return nil, errors.New("Unexpected message type in buffer area")
	}
	Packet.Body = make([]byte, Len)
	if _, err = VS.Conn.Read(Packet.Body); err != nil {
		log.Error("Error reading from TCP sock", zap.Error(err))
		return nil, err
	}
	if Msg == MSG_STREAMINFO_RSP {
		log.Debug("0x13c recieved StreamInfoRsp")
		log.Sugar().Debugf("RspBytes:%x", Packet.Body)
		Rsp, err := ParseStreamInfoRsp(Packet.Body)
		if err != nil {
			log.Error("Closing VTM Stream, error on StreamInfoRsp", zap.Error(err))
			return nil, err
		}
		log.Debug("StreamInfoResponse", zap.Int32p("Result", Rsp.Result),
			zap.Int32p("datakey", Rsp.Datakey),
			zap.Stringp("StreamHead", Rsp.Streamhead),
			zap.Stringp("StreamSSN", Rsp.Streamssn),
			zap.Stringp("VTMStreamKey", Rsp.Vtmstreamkey),
			zap.Stringp("ServerInfo", Rsp.Serverinfo),
			zap.Stringp("StreamURl", Rsp.Streamurl),
			zap.Stringp("SrvInfo", Rsp.Srvinfo),
			zap.Stringp("AesMD5", Rsp.Aesmd5),
			zap.Stringp("UDPtransinfo", Rsp.Udptransinfo),
			zap.Stringp("Peerpbkey", Rsp.Peerpbkey),
		)
		log.Sugar().Info("PDS", Rsp.Pdslist)
		return Rsp, nil
	} else {
		log.Error("Closing connection invalid sequence of messages", zap.Uint16("Recieved", Msg), zap.Uint16("Expected", MSG_STREAMINFO_RSP))
		VS.Conn.Close()
		return nil, errors.New("Invalid message code sequence")
	}
}
