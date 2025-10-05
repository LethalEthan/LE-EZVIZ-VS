# Encryption analysis

Connections to VTM (Video Transmission Manager) and VTDU (Video Transmission Data Unit) use E2EE when available i.e. ECDHE support is true and other encryption support checks pass, still yet to be hashed out as I've been looking more into the protocol than the device specific values. We generate ECDHE keypair with scep256r1 curve and send our publickey to the VTM/VTDU, we get the VTM/VTDU publickey from the pagelist or presumeably another API endpoint is available too.

The stream recieved from VTDU is encrypted with AES-128 which I believe uses the password the user provides with Image/Video Encryption or derives from it.

I see mentions of des_ecb3_cbc in stack traces, not sure if this is being used, or is part of another function.

Devices have a field called PermanentKey, I am not sure what this is for yet but I believe it can be decoded into an AES-128 key, since it is called permanent key I presume this isn't used directly for decrypting video streams but may be used to derive one.

Image encyption is pretty simple to do, look for hikencode and the 2 round md5 hash of the password, compare with your 2 round md5 hash of password if it is correct. Discard the first 48 bytes which is the hikencode and hash password head, copy the password into a 16 byte array and decrypt using AES CBC with the static IV and then unpad pkcs5. Whether or not a similar method is used for video encryption is not known, SRTP has a masterkey and session key where session keys are derived from the master key as I understand and SRTP uses either f8 or CTR mode, I am presuming the sequence in the header is used as a counter for CTR mode and that if f8 was used with unreliable transmission it'd be more difficult to use. If CTR mode is used and assuming it is SRTP then the IV/nonce is negotiated in some way, it shouldn't be static like it is with image encryption.
