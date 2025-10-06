package client

type MessageCode uint16

const (
	MSG_SIGNAL_MSG_TYPE_IDLE         = 0x00
	MSG_GET_VTDUINDO_REQ             = 0x12c
	MSG_GET_VTDUINFO_RSP             = 0x12d
	MSG_STARTSTREAM_REQ              = 0x12e
	MSG_STARTSTREAM_RSP              = 0x12f
	MSG_STOPSTREAM_REQ               = 0x130
	MSG_STOPSTREAM_RSP               = 0x131
	MSG_KEEPALIVE_REQ                = 0x132
	MSG_KEEPALIVE_RSP                = 0x133
	MSG_VTMSTREAM_NOTIFY             = 0x134
	MSG_PEERSTREAM_REQ               = 0x135
	MSG_PEERSTREAM_RSP               = 0x136
	MSG_GET_PLAYBACK_VTDU_REQ        = 0x137
	MSG_GET_PLAYBACK_VTDU_RSP        = 0x138
	MSG_START_PLAYBACK_REQ           = 0x139
	MSG_START_PLAYBACK_RSP           = 0x13a
	MSG_STREAMINFO_REQ               = 0x13b
	MSG_STREAMINFO_RSP               = 0x13c
	MSG_STREAMINFO_NOTIFY            = 0x13d
	MSG_SHAREDEV_TIMEOUT_NOTIFY      = 0x13e
	MSG_STREAM_MODIFY_SPEED_REQ      = 0x13f
	MSG_STREAM_MODIFY_SPEED_RSP      = 0x140
	MSG_STREAM_SEEK_REQ              = 0x141
	MSG_STREAM_SEEK_RSP              = 0x142
	MSG_STREAM_CONTINUE_REQ          = 0x143
	MSG_STREAM_CONTINUE_RSP          = 0x144
	MSG_STREAM_PAUSE_REQ             = 0x145
	MSG_STREAM_PAUSE_RSP             = 0x146
	MSG_STREAM_RESUME_REQ            = 0x147
	MSG_STREAM_RESUME_RSP            = 0x148
	MSG_LINKINFO_NOTIFY              = 0x149
	MSG_STREAM_VTMSTREAM_ECDH_NOTIFY = 0x14a
)
