# HTTP -> RTMP Streaming Proxy
[中文文档](README-CN.md)

An RTMP streaming proxy developed based on HTTP and ffmpeg, solving the problem of not being able to stream normal HTTP audio/video directly to an RTMP server.

ffmpeg installation is required for proper usage.

## Problem solved

1. Manual execution of shell commands for ffmpeg streaming.
2. Client uploads data in base64-encoded binary format.
3. Security issues with nginx-rtmp-module, currently only using "allow publish" to restrict access.
## Program parameters
```bash
ffmpeg-publisher -h
```
## Example

**Client request**
```
URL: htt://server_addr/base64
Method: POST
Content-Type: text/plain
Body: audio/video base64Data
```

**Server response**
```json
{
"errmsg": "ok",
"errcode": 0
}
```
