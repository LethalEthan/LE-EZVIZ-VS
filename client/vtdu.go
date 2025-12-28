package client

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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
	Codec         int
	Transport     int
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
		streamFile, err := os.Create("stream")
		if err != nil {
			panic(err)
		}
		go CreateKeepAliveTicker(VTDUstream.Conn, *Rsp.Streamssn)
		// var Packets = make([]VTMPacket, 0, 1600)
		var ATD bool
		for {
			Packet := new(VTMPacket)
			Packet.Header = make([]byte, 8)
			// VTDUstream.Conn.SetReadDeadline(time.Now().Add(5 * time.Second))
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
			Packet.Body = make([]byte, 0, Len)
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
					Packet.Body = append(Packet.Body[:n], Buf...)
					n += ReadByteCount
				}
				Packet.Body = Packet.Body[0:n]
			}
			if Chan == 0x01 {
				if VTDUstream.Transport == TRANS_UNKNOWN {
					VTDUstream.Transport = DetectTransport(Packet.Body)
				}
				if VTDUstream.Transport == TRANS_MPEG_PS {
					log.Debug("Transport is MPEG-PS")
					ATD = true
					// Find MPEG-PS, still under work
					// PS := FindPSwithinbuffer(Packet.Body)
					// if PS != nil {
					// }
					// Save packets to PACKS/ - not done so ignore this
					// Packets = append(Packets, *Packet)
					// if len(Packets) >= 1600 {
					// 	err := os.Mkdir("PACKS", os.ModeDir)
					// 	if err != nil {
					// 		panic(err)
					// 	}
					// }
					_, err := streamFile.Write(Packet.Body)
					if err != nil {
						panic(err)
					}
				}
				if VTDUstream.Transport == TRANS_RTP {
					log.Debug("Transport is RTP, this is still WIP, MPEG-PS is working currently")
					_, err := LEZ.DecodeRTP(Packet.Body)
					if err != nil {
						log.Error("Error decoding RTP", zap.Error(err))
					}
				}
			}
			if Chan == 0x00 && Seq == 0x00 && ATD {
				log.Info("Attempting ffmpeg re-encode to mp4")
				var eg errgroup.Group
				stream := ffmpeg.Input("./stream").
					Output("stream.mp4", ffmpeg.KwArgs{
						"vcodec": "copy",
						"f":      "mp4",
					})
				eg.Go(func() error {
					return stream.OverWriteOutput().Run()
				})
				if err := eg.Wait(); err != nil {
					log.Error("Error ffo: ", zap.Error(err))
					return err
				}
			}
			if len(Packet.Body) > 64 {
				log.Sugar().Debugf("packet 64b in hex: %x", Packet.Body[:64])
			} else {
				log.Sugar().Debugf("whole packet in hex: %x", Packet.Body)
			}
		}
		// streamFile.Sync()
		// streamFile.Close()
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

func CreateKeepAliveTicker(Sock net.Conn, StreamSSN string) {
	ticker := time.NewTicker(15 * time.Second)
	for t := range ticker.C {
		fmt.Println("Sending Keep Alive at", t)
		err := SendKeepAlive(Sock, StreamSSN)
		if err != nil {
			ticker.Stop()
			log.Error("Error sending keep alive", zap.Error(err))
			return
		}
	}
	ticker.Stop()
}
