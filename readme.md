# LE-EZVIZ Video Stream

LE-EZVIZ-VS, this is a piece of my wider project of creating a fully fledged program to control and connect to EZVIZ cameras in my own implementation. I have many modules as it makes testing easier and more will be released when ready. I am releasing this to hopefully spur on more development and get more help and hands on with the streaming implementation as not much is out there. Modern EZVIZ cameras have locked down their management ports, which is a shame and is disappointing as I myself want to record these streams or use them for home automation and using AI to detect faces and objects as I am sure many others want to do as well. It would not only be cheaper and more versatile for them and us to have the option of re-enabling those ports and having local control, it would save them costs of them having to stream and re-encode video streams and would also give us higher bitrate streams rather than cutting it down for internet transmission. It would also allow the cameras to be used in the event of an internet outage, I find it puzzling for why they decided to revoke local access as it benefits them and us, but I digress. I do have a little idea on perhaps how we can unlock the RTSP port again but I have no concrete proof or PoC atm and needs further investigation.

After many months of reverse engineering (2024 August - December, and 2025 September - now) First stage was general looking around and seeing what I could find, second stage was more specific on looking into the video streaming flow as my personal goal is to retrieve the video streams of my own devices and do whatever the hell I want with it.

I have developed an example of connecting to EZVIZ Video Management System. I'll delve into a little on how it works, during the login flow it gives a VTM address you connect to for video streams, VTM I believe stands for Video Transmission Manangement. VTM as I observe is a load balancer that then points you to a VTDU (I believe this means Video Transmission Data/Distribution Unit) that is the server that streams the video feed.

## Limitations

Currently encryption is unsupported, I have a decent idea on how it works and how it should be implemented but I do have bits of information missing. My expertise in using x64dbg is limited and I'll be honest I haven't read into it and it was my first time using it this month, nonetheless it at least gave me more insight into what functions are executed and what the keys may or may not be used for.

As seen with the VTM packet format the channel depicts whether it is an encrypted message or stream (0x00 for unencrypted message and 0x01 for unecrypted stream, 0x0a for encrypted message and 0x0b for encrypted stream), most of the time it is encrypted but some do not use it which is mainly older cameras. There are effectively 2 layers of encryption as you may guess from the Image/Video Encryption option. There is E2EE with VTM/VTDU which is still yet to be implemented and then the stream is encrypted with the password you set. The E2EE is ECDHE with Prime-256 curve, the stream I believe is AES-128 I am unsure which cipher mode it is but we'll cross the encryption bridge when we have a good foundation.

As I am only one person looking into this and not seeing much else online, information is quite scarce and with a limited amount of tools and knowledge I cannot provide a 100% implementation yet. That is why I am releasing this in it's current state to share my knowledge and induce collaboration and work on this so we can reimplement video streaming capabilities that we should have access to.

## What currently works
Currently as of 2026-01-04, MPEG-PS streaming works and so does H.265 RTP streams with no encryption enabled.

## If you want to help

If you want to help contribute, provide information or reimplement in another language it would be highly appreciated, feel free to join my [discord](https://discord.gg/3k5eBGsBkN) and if you need any pointers or help on where to look that's the best place to go to so we can give updates and share information. If you are creating a reimplementation let me know on github or discord and I'll create a list and add it to the repo or if you want to become a collaborator and create a branch. Remember to not include any sensitive information such as your keys, serials, tokens, password and any other kind of account information.

## To compile and run

Download [golang](https://go.dev/dl/) and install, open terminal and run `go build`. To use and run `./le-ezviz-vs -email="your@email.com" -password="password"` with no deviceSerial provided it will list all your devices, to have it begin the stream which for now doesn't do much but is a PoC for now add -deviceSerial like so `./le-ezviz-vs -email="your@email.com" -password="password" -deviceSerial="BC0123456"`

This will begin the stream and start spitting out the first 64 bytes of packets and their sequence. The stream may end after a while as I haven't setup the keepalive timer yet, the code isn't perfect and things will change but it does illustrate the things I have researched.

## My methods of reverse engineering

I have used a multitude of programs to help aid my efforts as listed below, research into different types of protocols and video formats also helped. What helped the most was packet analysis and finding bits and pieces of information and patterns along with looking at headers, structs and enums in ghidra.

* Ghidra
* x64dbg
* IEChooser
* Wireshark
* Charles Proxy
* MITMProxy
* Coffee, lots and lots of it
* Sleepless nights
* Shower thoughts
* Proxifier
* IDA
