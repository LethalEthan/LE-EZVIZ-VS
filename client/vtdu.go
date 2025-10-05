package client

import (
	"errors"
	"net"
	"strconv"

	"go.uber.org/zap"
)

type VTDUStream struct {
	Conn          net.Conn
	VTDUIP        string
	VTDUPort      int
	VTDUStreamSSN string
	VTMPublicKey  string
	VTMStreamKey  string
	SessionKey    string
	MasterKey     string
}

func (LEZ *LE_EZVIZ_Client) ConnectVTDU(vtduIP string, vtduPort int, vtmStreamKey, vtmPublicKey string) (*VTDUStream, error) {
	sock, err := net.Dial("tcp", vtduIP+":"+strconv.Itoa(vtduPort))
	if err != nil {
		log.Error("Error dialing VTM", zap.Error(err))
		return nil, err
	}
	VS := &VTDUStream{Conn: sock, VTDUIP: vtduIP, VTDUPort: vtduPort, VTMStreamKey: vtmStreamKey, VTMPublicKey: vtmPublicKey}
	return VS, nil
}

func (LEZ *LE_EZVIZ_Client) StartVTDUStream(VTDUstream *VTDUStream, StreamURL string) error {
	StreamReq, err := CreateStreamInfoReq(StreamURL, VTDUstream.VTMStreamKey)
	if err != nil {
		log.Error("Error creating StreamInfoReq", zap.Error(err))
		return err
	}
	EncodedPacket := EncodeVTMPacket(*StreamReq)
	log.Sugar().Debugf("bytes: %x", EncodedPacket)
	// return nil
	_, err = VTDUstream.Conn.Write(EncodedPacket)
	if err != nil {
		log.Error("Error writing to TCP sock", zap.Error(err))
		return err
	}
	Packet := new(VTMPacket)
	Packet.Header = make([]byte, 8)
	if _, err = VTDUstream.Conn.Read(Packet.Header); err != nil {
		log.Error("Error reading from TCP sock", zap.Error(err))
		return err
	}
	Len, _, _, Msg, err := Packet.DecodeHeader()
	if err != nil {
		log.Error("Error decoding packet header", zap.Error(err))
		return err
	}
	if Msg != MSG_STREAMINFO_RSP {
		return errors.New("Unexpected message type in buffer area")
	}
	Packet.Body = make([]byte, Len)
	if _, err = VTDUstream.Conn.Read(Packet.Body); err != nil {
		log.Error("Error reading from TCP sock", zap.Error(err))
		return err
	}
	if Msg == MSG_STREAMINFO_RSP {
		log.Debug("0x13c recieved StreamInfoRsp")
		log.Sugar().Debugf("RspBytes:%x", Packet.Body)
		Rsp, err := ParseStreamInfoRsp(Packet.Body)
		if err != nil {
			log.Error("Closing VTDU Stream, error on StreamInfoRsp", zap.Error(err))
			return err
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
		if Rsp.Result == nil {
			return errors.New("Result is nil")
		}
		if *Rsp.Result != 0 {
			err := LEZ.CheckRetcode(*Rsp.Result)
			return err
		}
		if Rsp.Streamssn == nil {
			return errors.New("Streamssn nil")
		}
		err = SendKeepAlive(VTDUstream.Conn, *Rsp.Streamssn)
		if err != nil {
			log.Error("Error sending keepalive to VTDU", zap.Error(err))
			return err
		}
		log.Debug("Print first 64 bytes of packets as an example, further handling is needed. I believe these are like RTSP interleaved packets over TCP")
		for {
			Packet := new(VTMPacket)
			Packet.Header = make([]byte, 8)
			if _, err = VTDUstream.Conn.Read(Packet.Header); err != nil {
				log.Error("Error reading from TCP sock", zap.Error(err))
				return err
			}
			Len, Chan, Seq, Msg, err := Packet.DecodeHeader()
			if err != nil {
				if err.Error() != "Magic not found" {
					log.Error("Error decoding packet header", zap.Error(err))
					return err
				} else {
					log.Debug("Current packet does not have magic, this is another protocol or continuation from previous TCP segment")
				}
			}
			log.Debug("Sequence", zap.Uint16("Seq", Seq))
			if Chan == 0x00 {
				log.Debug("Message channel Msgcode", zap.Uint16("code", Msg))
			}
			if Chan == 0x01 {
				log.Debug("Stream channel Msgcode", zap.Uint16("code", Msg))
			}
			if Msg == MSG_KEEPALIVE_REQ {
				log.Debug("Keepalive responded")
			}
			Packet.Body = make([]byte, Len)
			ReadByteCount, err := VTDUstream.Conn.Read(Packet.Body)
			if err != nil {
				log.Error("Error reading from TCP sock", zap.Error(err))
				return err
			}
			if ReadByteCount != int(Len) { // read until we fill the Len
				n := ReadByteCount
				// log.Debug("ReadByteCount is not equal to Len", zap.Int("ReadByteCount", ReadByteCount), zap.Uint16("Len", Len))
				for n != int(Len) {
					Buf := make([]byte, int(Len)-n)
					ReadByteCount, err := VTDUstream.Conn.Read(Buf)
					if err != nil {
						log.Error("Error reading from TCP sock", zap.Error(err))
						return err
					}
					Packet.Body = append(Packet.Body, Buf...)
					n += ReadByteCount
				}
			}
			if Chan == 0x01 {
				err = LEZ.DecodeRTP(Packet.Body)
				if err != nil {
					log.Error("Error decoding RTP", zap.Error(err))
				}
			}
			if len(Packet.Body) > 64 {
				log.Sugar().Debugf("packet 64b in hex: %x", Packet.Body[:64])
			} else {
				log.Sugar().Debugf("whole packet in hex: %x", Packet.Body)
			}

		}
	} else {
		log.Error("Closing connection invalid sequence of messages", zap.Uint16("Recieved", Msg), zap.Uint16("Expected", MSG_STREAMINFO_RSP))
		VTDUstream.Conn.Close()
		return errors.New("invalid message code sequence")
	}
	return nil
}

func SendKeepAlive(Sock net.Conn, StreamSSN string) error {
	bytes, err := CreateStreamKeepaliveReq(StreamSSN)
	if err != nil {
		return err
	}
	EncodedPacket := EncodeVTMPacket(*bytes)
	_, err = Sock.Write(EncodedPacket)
	if err != nil {
		return err
	}
	return nil
}
