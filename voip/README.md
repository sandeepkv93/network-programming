## Voice over IP (VoIP)

Voice over IP is a technology that allows voice communications and multimedia sessions over Internet Protocol networks.

## Table of Contents

1. [What is VoIP?](#what-is-voip)
2. [How Does VoIP Work?](#how-does-voip-work)
3. [Understanding the Code](#understanding-the-code)
4. [Further Reading](#further-reading)

### What is VoIP?

VoIP enables voice communication over IP networks by converting analog audio signals into digital data packets that can be transmitted over the internet.

**Key Components**:
- **Codec**: Encodes/decodes audio (G.711, Opus, etc.)
- **RTP**: Real-time Transport Protocol for audio delivery
- **SIP/H.323**: Signaling protocols for call setup
- **Jitter Buffer**: Compensates for network delays

**Advantages**:
- Lower costs compared to traditional telephony
- Integration with other services
- Scalability and flexibility
- Advanced features (voicemail, call forwarding, etc.)

### How Does VoIP Work?

1. **Audio Capture**: Microphone captures analog audio
2. **Digitization**: Convert to digital samples
3. **Encoding**: Compress using codec
4. **Packetization**: Split into RTP packets
5. **Transmission**: Send over IP network
6. **Reception**: Receive packets at destination
7. **Decoding**: Decompress audio
8. **Playback**: Convert to analog and play through speaker

**Protocol Stack**:
```
Application (VoIP)
    |
RTP/RTCP
    |
UDP
    |
IP
```

### Understanding the Code

This is a simplified VoIP implementation using UDP for low-latency audio transmission.

**Server**: Receives audio from clients and broadcasts to others
**Client**: Sends/receives audio packets

**Packet Format**:
- Timestamp (4 bytes)
- Sequence number (2 bytes)
- Reserved (2 bytes)
- Audio data (variable)

### Further Reading

- [VoIP - Wikipedia](https://en.wikipedia.org/wiki/Voice_over_IP)
- [RTP - Real-time Transport Protocol](https://en.wikipedia.org/wiki/Real-time_Transport_Protocol)
- [SIP - Session Initiation Protocol](https://en.wikipedia.org/wiki/Session_Initiation_Protocol)
