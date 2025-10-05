package client

import (
	"encoding/binary"
	"errors"
)

type VTMPacket struct {
	Header []byte
	Body   []byte
}

// This only encodes message unencrypted channel
func EncodeVTMPacket(data []byte) []byte {
	header := make([]byte, 8)
	header[0] = 0x24
	header[1] = 0
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], uint16(len(data)))
	header[2] = b[0]
	header[3] = b[1]
	header[4] = 0
	header[5] = 0
	var m [2]byte
	binary.BigEndian.PutUint16(m[:], uint16(0x13b))
	header[6] = m[0]
	header[7] = m[1]
	// log.Sugar().Debugf("ENCVTMP: %x", header)
	return append(header, data...)
}

func (p *VTMPacket) DecodeHeader() (Length uint16, Channel byte, Sequence uint16, Message uint16, err error) {
	if p.Header[0] != 0x24 {
		return 0, 0, 0, 0, errors.New("Magic not found")
	}
	switch p.Header[1] {
	case 0x00:
	case 0x01:
		break
	case 0x0a:
		return 0, 0, 0, 0, errors.New("encrypted channel currently unsupported")
	case 0x0b:
		return 0, 0, 0, 0, errors.New("encrypted channel currently unsupported")
	default:
		return 0, 0, 0, 0, errors.New("Unknown channel")
	}
	Length = binary.BigEndian.Uint16(p.Header[2:4])
	Sequence = binary.BigEndian.Uint16(p.Header[4:6])
	Message = binary.BigEndian.Uint16(p.Header[6:8])
	return Length, p.Header[1], Sequence, Message, nil
}
