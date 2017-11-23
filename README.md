# HTTP->RTMP推流Proxy
基于http+ffmpeg开发的RTMP推流代理，解决普通HTTP音/视频无法直接推流到RTMP服务器问题

需要安装ffmpeg才能正常使用

## 解决的问题
1. ffmpeg推流需要手动执行shell命令
2. 客户端上传的数据为base64编码过的二进制数据
3. nginx-rtmp-module安全问题，目前只会使用`allow publish`来限制

## 程序参数
```bash
ffmpeg-publisher -h
```

## HTTP推流
## 客户端请求
```
Url: htt://server_addr/base64
Method: POST
Content-Type: text/plain
Body: audio/video base64Data
```

## 服务端响应
```json
{
  "errmsg":"ok",
  "errcode":0
}
```