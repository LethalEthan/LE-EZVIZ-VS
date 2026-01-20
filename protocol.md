# Protocol Analysis

The protocol is custom and there is no documentation on it, I have made assumptions and guesses and they have worked out reasonably well so far but I did make many wrong assumptions during my analysis. It was no easy task, lots and lots of research went into this.

Analysing the network captures, I kept trying to find patterns, protocol headers and magic values, the packets always started with 0x24 which was my first clue I researched into many different protocol formats trying to see if it was a known standard that pointed me to RTSP interleaved which is similar in format but misses the sequence and message code field and I also haven't seen any RTCP messages. It's kind of like RTSP but instead of RTCP it's protobuf, RTCP might be used in different kinds of streams, there are many connection modes other than VTDU but I haven't seen them in use yet.


| Offset | Length  | Value            |
| ------ | ------- | ---------------- |
| 0      | 1 byte  | Magic Value 0x24 |
| 1      | 1 byte  | Channel          |
| 2      | 2 bytes | Length           |
| 4      | 2 bytes | Sequence         |
| 6      | 2 bytes | Message Code     |

It is a nice well formatted 8 byte header, the message codes that are first used is the StreamInfoReq (0x13b) and StreamInfoRsp (0x13c). What about the body, the body was more difficult for me to find out as I've never used protobufs before, I've re-made some of the protobuf structures, not all of them are included at this time.

Once the VTDU recieves streaminforeq and responds it then begins sending the stream, the stream I believe can be either MPEG-PS and RTP.
