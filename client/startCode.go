package client

import "go.uber.org/zap"

var packetBuf = make([]byte, 0, 4096)

func DecodeMPEG2PS(buf []byte) {
	hasStartCode := false
	if buf[0] == 0 && buf[1] == 0 && buf[2] == 1 && buf[3] == 0xba {
		// payload = buf[3:]
		log.Debug("Attmepting decode of startcode")
		hasStartCode = true
	}
	// if buf[0] == 0 && buf[1] == 0 && buf[2] == 0 && buf[3] == 1 {
	// 	payload = buf[3:]
	// 	log.Debug("Attmepting decode of startcode")
	// 	hasStartCode = true
	// }
	if hasStartCode && len(buf) > 4 {
		log.Debug("MPEG2-PS")
		PS := FindPSwithinbuffer(buf)
		if PS != nil {

		}
	}
}

// Find MPEG2-PS within a buffer, we return once another startcode is found
func FindPSwithinbuffer(buf []byte) []byte {
	if len(buf) > 10 {
		if len(packetBuf) == 0 {
			packetBuf = buf
			idx := findStartCode(packetBuf[4:])
			// if len(idx) >= 1 {
			if idx != -1 {
				log.Debug("KMPIDX", zap.Int("idx", idx))
				// log.Sugar().Debug("Next 4 bytes from idx: %x", packetBuf[idx[0]:5])
				// log.Sugar().Debugf("KMPIDX: %x", idx)
				rb := packetBuf[:idx]
				packetBuf = packetBuf[idx:]
				return rb
			}
			return nil
		}
		// idx := KMPSearch(packetBuf[4:], []byte{0, 0, 1, 0xBA})
		idx := findStartCode(packetBuf[4:])
		// if len(idx) >= 1 {
		if idx != -1 {
			log.Debug("KMPIDX", zap.Int("idx", idx))
			// log.Sugar().Debug("Next 4 bytes from idx: %x", packetBuf[idx[0]:5])
			// log.Sugar().Debugf("KMPIDX: %x", idx)
			rb := packetBuf[:idx]
			packetBuf = packetBuf[idx:]
			return rb
		} else {
			log.Debug("No KMP idx")
			packetBuf = append(packetBuf, buf...)
		}
		return nil
		// } else {
		// 	packetBuf = append(packetBuf, buf...)
		// 	idx := KMPSearch([]byte{0, 0, 1, 0xba}, packetBuf[4:])
		// 	if idx != -1 {
		// 		rb := packetBuf[:idx]
		// 		packetBuf = packetBuf[idx:]
		// 		return rb
		// 	}
		// }
	}
	return nil
	// if lastStartCodeIndex == 0 {
	// 	idx := KMPSearch([]byte{0, 0, 1, 0xba}, packetBuf[4:])
	// 	if idx != -1 {
	// 		return
	// 	}
	// }
}

func findStartCode(buf []byte) int {
	for i := 0; i < len(buf)-3; i++ {
		if buf[i] == 0 && buf[i+1] == 0 && buf[i+2] == 1 && buf[i+3] == 0xBA {
			return i
		}
	}
	return -1
}
