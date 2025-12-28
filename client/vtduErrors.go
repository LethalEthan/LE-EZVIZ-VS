package client

import (
	"errors"
	"strconv"
)

const (
	ERR_SERVER_EXCEPTION                        = "5000: server exception"
	ERR_CLIENT_WRONG_PARAMETERS                 = "5400: wrong client parameters"
	ERR_RECORDING_FILE_CANNOT_BE_FOUND          = "5402: the recording file cannot be found on the device"
	ERR_OPCODE_SIGNAL                           = "5403: opcode signaling key does not match device"
	ERR_DEVICE_OFFLINE                          = "5404: device offline"
	ERR_DEVICE_SIGNAL_TIMEOUT                   = "5405: signaling timeout/CAS response timeout after 10s"
	ERR_TOKEN_INVALID                           = "5406: token invalid"
	ERR_URL_MALFORMED                           = "5407: url is malformed"
	ERR_DEVICE_PRIVACY_PROTECTION_ON            = "5409: device has privacy protection on"
	ERR_TOKEN_UNAUTHORISED                      = "5411: token/user does not permission"
	ERR_SESSION_NOT_EXIST                       = "5412: session does not exist"
	ERR_TOKEN_VALIDATION                        = "5413: token validation failed"
	ERR_DEVICE_WRONG_CHANNEL                    = "5415: device wrong channel"
	ERR_RESOURCES_LIMITED                       = "5416: resources limited"
	ERR_STREAM_UNSUPPORTED                      = "5451: stream unsupported"
	ERR_DEVICE_LINK_STREAM_SERVER               = "5452: device link to stream server failed"
	ERR_NO_SESSION_ABOUT_DEVICE_STREAMING       = "5454: no session about device streaming"
	ERR_DEVICE_CHANNEL_UNASSOCIATED             = "5455: device channel is not associated"
	ERR_DEVICE_CHANNEL_ASSOCIATED_OFFLINE       = "5456: device channel associated device is offline"
	ERR_CLIENT_NOT_SUPPORT_E2EE                 = "5457: client doesn't support E2EE"
	ERR_DEVICE_NOT_SUPPORT_CONCURRENT_ECDH_PASS = "5458: device does not support concurrent ECDH password"
	ERR_VTDU_PROCESS_ECDH                       = "5459: VTDU failed to process ECDH encryption"
	ERR_SAME_REQUEST_PROCESSED                  = "5491: The same request is being processed and will now be rejected"
	ERR_DEVICE_COMMAND_UNSUPPORTED              = "5492: commands not support by device"
	ERR_SERVER_PROCESSING                       = "5500: server processing failed"
	ERR_VTM_ALLOCATE_VTDU                       = "5503: VTM failed to allocate VTDU"
	ERR_VTDU_STREAMING_LIMIT_REACHED            = "5504: streaming VTDUs limit reached for user"
	ERR_DEVICE_NO_VIDEO_SOURCE                  = "5544: device returns no video source"
	ERR_VIDEO_SHARING_TIME_ENDED                = "5545: video sharing time has ended"
	ERR_VTDU_CONCURRENT_CHANNEL_LIMIT           = "5546: VTDU concurrent channels limit reached"
	ERR_DEVICE_PACKET_TOO_LARGE                 = "6518: device packet too large"
	ERR_DEVICE_NETWORK_LINK_UNSTABLE            = "6519: device network link is unstable"
	ERR_NETWORK_UNSTABLE                        = "6520: unstable network"
	ERR_VTDU_DISCONNECTED                       = "7005: VTDU disconnected"
)

func (LEZ *LE_EZVIZ_Client) CheckRetcode(ret int32) error {
	switch ret {
	case 5000:
		return errors.New(ERR_SERVER_EXCEPTION)
	case 5400:
		return errors.New(ERR_CLIENT_WRONG_PARAMETERS)
	case 5402:
		return errors.New(ERR_RECORDING_FILE_CANNOT_BE_FOUND)
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
	case 5409:
		return errors.New(ERR_DEVICE_PRIVACY_PROTECTION_ON)
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
	case 5491:
		return errors.New(ERR_SAME_REQUEST_PROCESSED)
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
	return errors.New("Response result not 0 instead got: " + strconv.Itoa(int(ret)))
}
