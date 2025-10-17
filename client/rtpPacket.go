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
		return nil, nil
	}
	Padding := (buf[0] & PaddingBit) != 0
	Extension := (buf[0] & ExtensionBit) != 0
	ContributionCount := buf[0] & ContributionBits
	Marker := (buf[1] & 0x80) != 0
	PayloadType := buf[1] & 0x7F
	SequenceNumber := binary.BigEndian.Uint16(buf[2:4])
	TimeStamp := binary.BigEndian.Uint32(buf[4:8])
	SSI := binary.BigEndian.Uint32(buf[8:12])
	offset := RTPFixedHeaderLen
	if ContributionCount > 0 {
		ContributionLength := len(buf) * 4
		if len(buf) < ContributionLength+offset {
			return nil, errors.New("csrc is greater than buffer provided")
		}
		CSRC := make([]uint32, ContributionCount)
		for i := 0; i < int(ContributionCount); i++ {
			CSRC[i] = binary.BigEndian.Uint32(buf[offset+i*4 : offset+i*4+4])
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
		log.Debug("RTP Extension", zap.Uint16("Profile", ExtensionProfile), zap.Uint16("Length", ExtensionLength))
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
	}
	if len(buf) < offset {
		return nil, errors.New("buffer is smaller than offset")
	}
	Payload := buf[offset:]
	if Padding {
		if len(Payload) == 0 {
			return nil, errors.New("no more bytes cannot be pad")
		}
		PadLen := int(Payload[len(Payload)-1])
		if PadLen == 0 || len(Payload) < PadLen || PadLen > 255 {
			return nil, errors.New("invalid padding length")
		}
		Payload = Payload[:len(Payload)-PadLen]
	}
	if len(Payload) > 10 {
		h265 := (Payload[0] >> 1) & 0x3f
		if h265 == 0x30 {
			log.Debug("H265/HEVC frame more work needed")
		}
		if h265 == 0x31 {
			log.Debug("H265/HEVC frame")
			if Payload[2]&0x80 != 0 {
				Payload = AddAVCStartCode(Payload)
				log.Debug("RTP Header", zap.Uint8("Ver", Version), zap.Bool("Pad", Padding), zap.Bool("Ext", Extension), zap.Uint8("CC", ContributionCount), zap.Bool("Mark", Marker), zap.Uint8("PayloadT", PayloadType), zap.Uint16("Seq", SequenceNumber), zap.Uint32("Time", TimeStamp), zap.Uint32("SSI", SSI), zap.Int("PayloadLen", len(Payload)))
				log.Sugar().Debugf("RTP Payload: %x", Payload)
				return Payload, nil
			} else {
				// missing more work needed on processing frame
			}
		} //!=0x32, need to extract profile and codecinfo first
		h264 := Payload[0] & 0x1f
		if h264 != 9 {
			switch h264 {
			case 0x18:
				log.Debug("H264/AVC frame more work needed")
			case 0x1c:
				log.Debug("H264/AVC frame")
				if (Payload[1] & 0xc0) == 0x80 {
					Payload = AddAVCStartCode(Payload)
					log.Debug("RTP Header", zap.Uint8("Ver", Version), zap.Bool("Pad", Padding), zap.Bool("Ext", Extension), zap.Uint8("CC", ContributionCount), zap.Bool("Mark", Marker), zap.Uint8("PayloadT", PayloadType), zap.Uint16("Seq", SequenceNumber), zap.Uint32("Time", TimeStamp), zap.Uint32("SSI", SSI), zap.Int("PayloadLen", len(Payload)))
					log.Sugar().Debugf("RTP Payload: %x", Payload)
					return nil, nil
				}
			default:
				log.Debug("H264/AVC frame with NAL")
				log.Sugar().Debugf("RTP Payload: %x", Payload)
				Payload = AddAVCStartCode(Payload)
				return nil, nil
				// log.Sugar().Debugf("first 10b: %x", Payload[0:10])

			}
		} else {
			log.Debug("Unknown H264/AVC")
		}
	}
	log.Debug("RTP Header", zap.Uint8("Ver", Version), zap.Bool("Pad", Padding), zap.Bool("Ext", Extension), zap.Uint8("CC", ContributionCount), zap.Bool("Mark", Marker), zap.Uint8("PayloadT", PayloadType), zap.Uint16("Seq", SequenceNumber), zap.Uint32("Time", TimeStamp), zap.Uint32("SSI", SSI), zap.Int("PayloadLen", len(Payload)))
	log.Sugar().Debugf("RTP Payload: %x", Payload)
	return nil, nil
}

// AVCStartCode/NAL we replace the first four bytes after we check it is H264/H265
func AddAVCStartCode(buf []byte) []byte {
	buf[0] = 0
	buf[1] = 0
	buf[2] = 0
	buf[3] = 1
	return buf
}
