## Video Streaming

Video streaming is the continuous transmission of video files from a server to a client, allowing users to watch video content without downloading the entire file first.

## Table of Contents

1. [What is Video Streaming?](#what-is-video-streaming)
2. [How Does Video Streaming Work?](#how-does-video-streaming-work)
3. [Understanding the Code](#understanding-the-code)
4. [Further Reading](#further-reading)

### What is Video Streaming?

Video streaming delivers video content over the internet in a continuous flow, allowing users to watch content as it's being delivered without waiting for the complete download.

**Streaming Protocols**:
- **HTTP Live Streaming (HLS)**: Apple's adaptive bitrate protocol
- **MPEG-DASH**: Dynamic Adaptive Streaming over HTTP
- **RTMP**: Real-Time Messaging Protocol (Adobe)
- **WebRTC**: Real-time communication for browsers

**Adaptive Bitrate Streaming**:
- Multiple quality levels (360p, 720p, 1080p, 4K)
- Automatically adjusts based on network conditions
- Provides smooth playback experience

### How Does Video Streaming Work?

1. **Content Preparation**: Encode video in multiple qualities
2. **Segmentation**: Split into small chunks (2-10 seconds)
3. **Manifest File**: Create playlist (M3U8 for HLS, MPD for DASH)
4. **Client Request**: Player requests manifest
5. **Adaptive Selection**: Choose appropriate quality based on bandwidth
6. **Chunk Download**: Download and buffer video segments
7. **Playback**: Play video while downloading next segments
8. **Quality Switching**: Adjust quality as network changes

**Buffering Strategy**:
- Pre-buffer initial segments before playback
- Maintain buffer to handle network fluctuations
- Balance between startup time and smooth playback

### Understanding the Code

This implementation provides HTTP-based video streaming with viewer tracking.

**Features**:
- Multiple concurrent streams
- Viewer counting
- Stream management API
- Chunked transfer encoding
- Connection keep-alive

**Endpoints**:
- `/`: Main page with stream list
- `/stream/{id}`: Watch a specific stream
- `/api/streams`: List all streams (JSON)
- `/api/stream/create`: Create new stream (POST)

### Further Reading

- [HTTP Live Streaming](https://en.wikipedia.org/wiki/HTTP_Live_Streaming)
- [MPEG-DASH](https://en.wikipedia.org/wiki/Dynamic_Adaptive_Streaming_over_HTTP)
- [Adaptive Bitrate Streaming](https://en.wikipedia.org/wiki/Adaptive_bitrate_streaming)
- [Video Codecs](https://en.wikipedia.org/wiki/Video_codec)
