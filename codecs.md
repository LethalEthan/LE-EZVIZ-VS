# Codecs

The common codecs used are H.265 and H.264, the RTP payload that I have seen omits the start code so we have to add it in, it makes sense to not waste the bytes but does make differentiating different codecs and stream types harder. To determine whether it is HEVC or AVC we check if there is a valid NAL unit type, usually PPS or SPS. Some streams have the start code already there so we leave it and attempt to decode it as a h264/265 frame or as MPEG2-PS. There is more that we have to do with the RTP payload as the frame cannot be decoded as is, there are RTP extensions that I believe determine the codec and profile infornation. You can use ffprobe and get some information out of the frame such as framerate and bitrate but the other information is likely set in the RTP extension header when the stream begins.

## Container formats

The carrier of these codecs also vary between devices, we of course have RTP which is what most of the work has been done on to decode some of the data but there is: MPEG-4, MPEG-TS, MPEG-PS. I have seen MPEG2-PS and I am attempting work on decoding this too. RTP decodes the data without the start code, at first I thought I had to overwrite the first 4 bytes but in testing it doesn't decode right so we'll add it at the start.

Upon attempting a rudimentary mpeg2-ps detector and storing the stream into a file, I am able to get some video output albeit corrupted. ffmpeg is seeing invalid NAL unit types, what I am guessing is that there is non-standard NAL unit types being used or some other type of data being mixed in the video stream. What I am leaning towards is perhaps it's the unit types are shifted whether it's because of purposeful obfuscation or an actual purpose remains to be seen.

I have to work on more of the decoding to be sure and to start messing with the data to find out what is being changed. Encryption is also going to be a bigger ball ache to figure out if this is the case.
