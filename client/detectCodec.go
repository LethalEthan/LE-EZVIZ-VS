package client

const (
	CODEC_UNKNOWN = iota
	CODEC_H264
	CODEC_H265
)
const (
	TRANS_UNKNOWN = iota
	TRANS_MPEG4
	TRANS_MPEG_TS
	TRANS_MPEG_PS
	TRANS_RTP
)

func DetectCodec(buf []byte, transport int) int {
	// if transport ==
	h264 := buf[0] & 0x1f
	h265 := (buf[0] >> 1) & 0x3f
	switch h264 {
	case 7: //SPS
		return CODEC_H264
	case 8: //PPS
		return CODEC_H264
	}
	switch h265 {
	case 31: // SPS
		return CODEC_H265
	case 32: //PPS
		return CODEC_H265
	}
	return CODEC_UNKNOWN
}

func DetectTransport(buf []byte) int {
	if buf[0] == 0 && buf[1] == 0 && buf[2] == 1 && buf[3] == 0xba {
		return TRANS_MPEG_PS
	}
	if buf[0] == 0x47 {
		return TRANS_MPEG_TS
	}
	rtpVersion := buf[0] >> 6
	if rtpVersion == 2 {
		return TRANS_RTP
	}
	return TRANS_UNKNOWN
}
