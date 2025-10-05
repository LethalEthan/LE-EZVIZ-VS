package client

import (
	"errors"
	"strconv"
)

const (
	ERR_SERVER_EXCEPTION                        = "server exception"
	ERR_CLIENT_WRONG_PARAMETERS                 = "wrong client parameters"
	ERR_OPCODE_SIGNAL                           = "opcode signaling key does not match device"
	ERR_DEVICE_OFFLINE                          = "device offline"
	ERR_DEVICE_SIGNAL_TIMEOUT                   = "signaling timeout/CAS response timeout"
	ERR_TOKEN_INVALID                           = "token invalid"
	ERR_URL_MALFORMED                           = "url is malformed"
	ERR_TOKEN_UNAUTHORISED                      = "token/user does not permission"
	ERR_SESSION_NOT_EXIST                       = "session does not exist"
	ERR_TOKEN_VALIDATION                        = "token validation failed"
	ERR_DEVICE_WRONG_CHANNEL                    = "device wrong channel"
	ERR_RESOURCES_LIMITED                       = "resources limited"
	ERR_STREAM_UNSUPPORTED                      = "stream unsupported"
	ERR_DEVICE_LINK_STREAM_SERVER               = "deveice link to stream server failed"
	ERR_NO_SESSION_ABOUT_DEVICE_STREAMING       = "no session about device streaming"
	ERR_DEVICE_CHANNEL_UNASSOCIATED             = "device channel is not associated"
	ERR_DEVICE_CHANNEL_ASSOCIATED_OFFLINE       = "device channel associated device is offline"
	ERR_CLIENT_NOT_SUPPORT_E2EE                 = "client doesn't support E2EE"
	ERR_DEVICE_NOT_SUPPORT_CONCURRENT_ECDH_PASS = "device does not support concurrent ECDH password"
	ERR_VTDU_PROCESS_ECDH                       = "VTDU failed to process ECDH encryption"
	ERR_DEVICE_COMMAND_UNSUPPORTED              = "commands not support by device"
	ERR_SERVER_PROCESSING                       = "server processing failed"
	ERR_VTM_ALLOCATE_VTDU                       = "VTM failed to allocate VTDU"
	ERR_VTDU_STREAMING_LIMIT_REACHED            = "streaming VTDUs limit reached for user"
	ERR_DEVICE_NO_VIDEO_SOURCE                  = "device returns no video source"
	ERR_VIDEO_SHARING_TIME_ENDED                = "video sharing time has ended"
	ERR_VTDU_CONCURRENT_CHANNEL_LIMIT           = "VTDU concurrent channels limit reached"
	ERR_DEVICE_PACKET_TOO_LARGE                 = "device packet too large"
	ERR_DEVICE_NETWORK_LINK_UNSTABLE            = "device network link is unstable"
	ERR_NETWORK_UNSTABLE                        = "unstable network"
	ERR_VTDU_DISCONNECTED                       = "VTDU disconnected"
)

func (LEZ *LE_EZVIZ_Client) CheckRetcode(ret int32) error {
	switch ret {
	case 5000:
		return errors.New(ERR_SERVER_EXCEPTION)
	case 5400:
		return errors.New(ERR_CLIENT_WRONG_PARAMETERS)
	case 5403:
		return errors.New(ERR_OPCODE_SIGNAL)
	case 5404:
		return errors.New(ERR_DEVICE_OFFLINE)
	case 5405:
		return errors.New(ERR_DEVICE_SIGNAL_TIMEOUT)
	case 5406:
		return errors.New(ERR_TOKEN_INVALID)
	case 5407:
		return errors.New(ERR_URL_MALFORMED)
	case 5411:
		return errors.New(ERR_TOKEN_UNAUTHORISED)
	case 5412:
		return errors.New(ERR_SESSION_NOT_EXIST)
	case 5413:
		return errors.New(ERR_TOKEN_VALIDATION)
	case 5415:
		return errors.New(ERR_DEVICE_WRONG_CHANNEL)
	case 5416:
		return errors.New(ERR_RESOURCES_LIMITED)
	case 5451:
		return errors.New(ERR_STREAM_UNSUPPORTED)
	case 5452:
		return errors.New(ERR_DEVICE_LINK_STREAM_SERVER)
	case 5454:
		return errors.New(ERR_NO_SESSION_ABOUT_DEVICE_STREAMING)
	case 5455:
		return errors.New(ERR_DEVICE_CHANNEL_UNASSOCIATED)
	case 5456:
		return errors.New(ERR_DEVICE_CHANNEL_ASSOCIATED_OFFLINE)
	case 5457:
		return errors.New(ERR_CLIENT_NOT_SUPPORT_E2EE)
	case 5458:
		return errors.New(ERR_DEVICE_NOT_SUPPORT_CONCURRENT_ECDH_PASS)
	case 5459:
		return errors.New(ERR_VTDU_PROCESS_ECDH)
	case 5492:
		return errors.New(ERR_DEVICE_COMMAND_UNSUPPORTED)
	case 5500:
		return errors.New(ERR_SERVER_PROCESSING)
	case 5503:
		return errors.New(ERR_VTM_ALLOCATE_VTDU)
	case 5504:
		return errors.New(ERR_VTDU_STREAMING_LIMIT_REACHED)
	case 5544:
		return errors.New(ERR_DEVICE_NO_VIDEO_SOURCE)
	case 5545:
		return errors.New(ERR_VIDEO_SHARING_TIME_ENDED)
	case 5546:
		return errors.New(ERR_VTDU_CONCURRENT_CHANNEL_LIMIT)
	case 6518:
		return errors.New(ERR_DEVICE_PACKET_TOO_LARGE)
	case 6519:
		return errors.New(ERR_DEVICE_NETWORK_LINK_UNSTABLE)
	case 6520:
		return errors.New(ERR_NETWORK_UNSTABLE)
	case 7005:
		return errors.New(ERR_VTDU_DISCONNECTED)
	}
	return errors.New("Response result not 0 " + strconv.Itoa(int(ret)))
}
