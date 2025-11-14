## Video over IP

Video over IP transmits video signals over IP networks, enabling video surveillance, broadcasting, and communication systems over standard network infrastructure.

## Table of Contents

1. [What is Video over IP?](#what-is-video-over-ip)
2. [How Does Video over IP Work?](#how-does-video-over-ip-work)
3. [Understanding the Code](#understanding-the-code)
4. [Further Reading](#further-reading)

### What is Video over IP?

Video over IP transmits digital video over an IP network. It's used in IP cameras, video conferencing, IPTV, and video surveillance systems.

**Key Technologies**:
- **Video Codecs**: H.264, H.265/HEVC, VP9, AV1
- **Container Formats**: MP4, WebM, MPEG-TS
- **Streaming Protocols**: RTP, RTSP, HLS, DASH
- **Network Protocols**: UDP (low latency) or TCP (reliability)

**Applications**:
- IP surveillance cameras
- Video conferencing
- IPTV and streaming
- Remote monitoring

### How Does Video over IP Work?

1. **Video Capture**: Camera captures video frames
2. **Encoding**: Compress using codec (H.264, etc.)
3. **Packetization**: Split into network packets
4. **Transmission**: Send over IP network via UDP/TCP
5. **Reception**: Receive packets
6. **Buffering**: Handle jitter and packet loss
7. **Decoding**: Decompress video
8. **Display**: Show on screen

**Frame Types**:
- **I-frame**: Independent frame (keyframe)
- **P-frame**: Predicted frame (delta from previous)
- **B-frame**: Bi-directional predicted frame

### Understanding the Code

Simplified video transmission over UDP.

**Packet Format**:
- Timestamp (4 bytes)
- Sequence number (2 bytes)
- Frame type (1 byte): I/P/B frame
- Fragment ID (1 byte)
- Reserved (2 bytes)
- Video data (variable)

### Further Reading

- [Video over IP - Wikipedia](https://en.wikipedia.org/wiki/Video_over_IP)
- [H.264/AVC](https://en.wikipedia.org/wiki/Advanced_Video_Coding)
- [RTP - Real-time Transport Protocol](https://en.wikipedia.org/wiki/Real-time_Transport_Protocol)
