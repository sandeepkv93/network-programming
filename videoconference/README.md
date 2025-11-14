## Video Conferencing

Video conferencing enables real-time audio and video communication between multiple participants over the internet using WebRTC technology.

## What is Video Conferencing?

Video conferencing allows multiple users to have face-to-face meetings over the internet. Modern implementations use WebRTC for peer-to-peer communication with a signaling server for connection establishment.

**Key Components**:
- **WebRTC**: Peer-to-peer audio/video communication
- **Signaling Server**: Coordinates connection establishment
- **STUN/TURN Servers**: NAT traversal and relay
- **Media Servers**: For multi-party conferences (SFU/MCU)

**Architectures**:
- **Mesh**: Each peer connects to every other peer
- **SFU**: Selective Forwarding Unit (routes streams)
- **MCU**: Multipoint Control Unit (mixes streams)

## Further Reading

- [WebRTC](https://webrtc.org/)
- [Video Conferencing - Wikipedia](https://en.wikipedia.org/wiki/Videotelephony)
