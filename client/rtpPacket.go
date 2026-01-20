package client

import (
	"encoding/binary"
	"errors"

	"go.uber.org/zap"
)

const RTPFixedHeaderLen = 12
const PaddingBit byte = 0x20
const ExtensionBit byte = 0x10
const ContributionBits byte = 0x0f
const Marker byte = 0x80

/*
 128 64  32   16  8 4 2 1
|Version| P | X |   CC   |
*/

// A rudimentary RTP header decode
func (LEZ *LE_EZVIZ_Client) DecodeRTP(buf []byte) ([]byte, error) {
	if buf[0] == 0 && buf[1] == 0 && buf[2] == 1 || buf[0] == 0 && buf[1] == 0 && buf[2] == 0 && buf[3] == 1 {
		log.Debug("This is not RTP, NAL header?")
		return nil, nil
	}
	Version := buf[0] >> 6
	if Version != 2 {
		log.Debug("RTP version is not 2, this is a different stream type not yet supported")
		log.Sugar().Debugf("Received: %x", buf)
		return nil, nil
	}
	Padding := (buf[0] & PaddingBit) != 0
	PaddingLength := 0
	Extension := (buf[0] & ExtensionBit) != 0
	ContributionCount := buf[0] & ContributionBits
	Marker := (buf[1] & 0x80) != 0
	PayloadType := buf[1] & 0x7F
	SequenceNumber := binary.BigEndian.Uint16(buf[2:4])
	TimeStamp := binary.BigEndian.Uint32(buf[4:8])
	SSI := binary.BigEndian.Uint32(buf[8:12])
	offset := RTPFixedHeaderLen
	if ContributionCount > 0 {
		ContributionLength := int(ContributionCount) * 4
		if len(buf) < ContributionLength+offset {
			return nil, errors.New("csrc is greater than buffer provided")
		}
		CSRC := make([]uint32, ContributionCount)
		for i := 0; i < int(ContributionCount); i++ {
			CSRC[i] = binary.BigEndian.Uint32(buf[offset+i*4 : offset+(i*4)+4])
			log.Debug("RTP CSRC", zap.Uint32("CSRC", CSRC[i]))
		}
		offset += ContributionLength
	}
	// var ExtensionData []byte
	if Extension {
		if len(buf) < offset+4 {
			return nil, errors.New("extensionc is greater than buffer provided")
		}
		ExtensionProfile := binary.BigEndian.Uint16(buf[offset : offset+2])
		ExtensionLength := binary.BigEndian.Uint16(buf[offset+2 : offset+4])
		log.Debug("RTP Extension", zap.Uint16("Profile", ExtensionProfile), zap.Uint16("Length", ExtensionLength*4), zap.Uint8("PayloadT", PayloadType))
		offset += 4
		ExtensionBytes := int(ExtensionLength) * 4
		if len(buf) < offset+ExtensionBytes {
			return nil, errors.New("extension is greater than buffer provided")
		}
		if ExtensionBytes > 0 {
			// ExtensionData = make([]byte, ExtensionBytes)
			// copy(ExtensionData, buf[offset:offset+ExtensionBytes])
			log.Sugar().Debugf("RTP Extension Data: %x", buf[offset:offset+ExtensionBytes])
		} else {
			// ExtensionData = nil
		}
		offset += ExtensionBytes
		// return nil, nil
	}
	if len(buf) < offset {
		return nil, errors.New("buffer is smaller than offset")
	}
	Payload := buf[offset:]
	if Padding {
		if len(Payload) == 0 {
			return nil, errors.New("no more bytes cannot be pad")
		}
		PaddingLength = int(Payload[len(Payload)-1])
		if PaddingLength == 0 || len(Payload) < PaddingLength || PaddingLength > 255 {
			return nil, errors.New("invalid padding length")
		}
		Payload = Payload[:len(Payload)-PaddingLength]
	}
	var RTPnal byte
	var NAL byte
	var Frame = make([]byte, 0, len(Payload))
	switch PayloadType {
	//Audio
	case 0:
	case 8:
	case 0xb:
	case 0xe:
	case 0x62:
	case 100:
	case 0x66:
	case 0x67:
	case 0x68:
	case 0x73:
		// RTPProcessAudio(buf)
	}
	if len(Payload) > 2 {
		if PayloadType == 96 {
			RTPnal = (Payload[0] >> 1) & 0x3F
			if RTPnal == 0x20 {
				log.Debug("VPS h265")
			}
			if RTPnal == 0x21 {
				log.Debug("SPS h265")
			}
			if RTPnal == 0x22 {
				log.Debug("PPS h265")
			}
			//Aggregated - TBD
			if RTPnal == 0x30 {
				log.Debug("Aggregated, not supported in progress this may cause a drop/skip")
				return nil, nil
			}
			if RTPnal == 0x31 {
				if Payload[2]&0x80 != 0 { // Start packet
					NAL = (Payload[2]*2^Payload[0])&0x7e ^ Payload[0]
					Frame = AppendAVCStartCode(Frame)
					Frame = append(Frame, []byte{NAL, 0x01}...)
				}
				Frame = append(Frame, Payload[3:]...)
			} else {
				if RTPnal != 0x32 {
					Frame = AppendAVCStartCode(Payload)
				} else {
					log.Debug("rtpnal is 0x32")
				}
			}
			log.Sugar().Debugf("RTPNAL: %x", RTPnal)
			log.Sugar().Debugf("NAL: %x", NAL)
			log.Sugar().Debugf("RTP Payload: %x", Payload)
		}
		log.Debug("RTP Header", zap.Uint8("Ver", Version), zap.Bool("Pad", Padding), zap.Int("PadLen", PaddingLength), zap.Bool("Ext", Extension), zap.Uint8("CC", ContributionCount), zap.Bool("Mark", Marker), zap.Uint8("PayloadT", PayloadType), zap.Uint16("Seq", SequenceNumber), zap.Uint32("Time", TimeStamp), zap.Uint32("SSI", SSI), zap.Int("PayloadLen", len(Payload)))

		// Payload = AppendAVCStartCode(Payload)
		return Frame, nil
	} else {
		log.Debug("RTP Header", zap.Uint8("Ver", Version), zap.Bool("Pad", Padding), zap.Int("PadLen", PaddingLength), zap.Bool("Ext", Extension), zap.Uint8("CC", ContributionCount), zap.Bool("Mark", Marker), zap.Uint8("PayloadT", PayloadType), zap.Uint16("Seq", SequenceNumber), zap.Uint32("Time", TimeStamp), zap.Uint32("SSI", SSI), zap.Int("PayloadLen", len(Payload)))
		log.Sugar().Debugf("RTP Unkown Payload: %x", Payload)
		return nil, nil
	}
	log.Debug("RTP Header", zap.Uint8("Ver", Version), zap.Bool("Pad", Padding), zap.Bool("Ext", Extension), zap.Uint8("CC", ContributionCount), zap.Bool("Mark", Marker), zap.Uint8("PayloadT", PayloadType), zap.Uint16("Seq", SequenceNumber), zap.Uint32("Time", TimeStamp), zap.Uint32("SSI", SSI), zap.Int("PayloadLen", len(Payload)))
	log.Sugar().Debugf("RTP Payload: %x", Payload)
	return nil, nil
}

// Currently unsure whether to replace or append the start code
// Added both so we can mess around with the stream

// AVCStartCode/NAL we replace the first four bytes after we check it is H264/H265
func ReplaceAVCStartCode(buf []byte) []byte {
	buf[0] = 0
	buf[1] = 0
	buf[2] = 0
	buf[3] = 1
	return buf
}

func AppendAVCStartCode(buf []byte) []byte {
	avc := make([]byte, 0, len(buf)+4)
	avc = append(avc, []byte{0, 0, 0, 1}...)
	avc = append(avc, buf...)
	return avc
}
