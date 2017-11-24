package main

import (
	"github.com/alecthomas/log4go"
	"flag"
	"net/http"
	"time"
	"encoding/json"
	"io/ioutil"
	"encoding/base64"
	"os"
	"os/exec"
	"strings"
)

var (
	rtmp            string // rtmp 服务端地址
	ffmpeg          string // ffmpeg命令地址
	listen          string // 监听地址
	uploadKey       string // 上传key
	uploadDirectory string // 本地上传目录
	start           = time.Now()
)

const (
	VERSION = "1.0.0"
)

func init() {
	flag.StringVar(&rtmp, "rtmp", "rtmp://localhost:1935/hls", "rtmp upstream address")
	flag.StringVar(&ffmpeg, "ffmpeg", "ffmpeg", "ffmpeg executable")
	flag.StringVar(&listen, "listen", ":8081", "http server address")
	flag.StringVar(&uploadKey, "key", "testkey", "http upload key")
	flag.StringVar(&uploadDirectory, "upload-directory", "upload", "local upload directory")
}

func main() {
	flag.Parse()
	log4go.LoadConfiguration("log4go.xml")
	if rtmp == "" {
		panic("rtmp upstream not set!")
	}
	// 创建目录
	os.Mkdir(uploadDirectory, os.ModePerm)
	log4go.Info("start ffmpeg-publisher, version %s", VERSION)
	log4go.Info("RTMP UpStream %s", rtmp)
	log4go.Info("HTTP Upload Directory %s", uploadDirectory)
	log4go.Info("HTTP Listen ON %s", listen)
	log4go.Info("HTTP Upload Key %s", uploadKey)

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/base64", handleBase64Publish)
	panic(http.ListenAndServe(listen, nil))
}

func sendResponse(errmsg string, statusCode int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	if statusCode == 0 {
		statusCode = 200
	}
	w.WriteHeader(statusCode)
	ret := map[string]interface{}{
		"errmsg":  errmsg,
		"errcode": statusCode,
	}
	json.NewEncoder(w).Encode(&ret)
	if statusCode == 500 {
		log4go.Error(errmsg)
	}
}

func handleIndex(w http.ResponseWriter, _ *http.Request) {
	ret := map[string]interface{}{
		"name":    "ffpmeg-publisher",
		"version": VERSION,
		"uptime":  time.Since(start).String(),
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	json.NewEncoder(w).Encode(&ret)
}

func handleBase64Publish(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		sendResponse("Method Not Allowed", 405, w)
		return
	}
	// 获取直播ID和key
	if err := r.ParseForm(); err != nil {
		sendResponse(err.Error(), 500, w)
		return
	}
	// 检测直播Key
	if key := r.Form.Get("key"); key != uploadKey {
		sendResponse("Forbidden", 403, w)
		return
	}
	// 检测直播ID
	id := r.Form.Get("id")
	if id == "" {
		sendResponse("Bad Request", 400, w)
		return
	}
	// 读取base64
	base64Buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendResponse(err.Error(), 500, w)
		return
	}
	buf, err := base64.StdEncoding.DecodeString(string(base64Buf))
	if err != nil {
		sendResponse(err.Error(), 500, w)
		return
	}
	log4go.Trace("live %s receive %d bytes", id, len(buf))
	if err := publishVoice(buf, id); err != nil {
		sendResponse(err.Error(), 500, w)
	} else {
		sendResponse("ok", 0, w)
	}
}

// 保存音频到临时文件目录,1个直播一个目录
func saveVoice(buf []byte, id string) (string, error) {
	directory := uploadDirectory + "/" + id
	os.Mkdir(directory, os.ModePerm)
	name := time.Now().Format("2006-01-02-15-04-05") + ".mp3"
	filename := directory + "/" + name
	log4go.Trace("live %s save chunk %s", id, filename)
	err := ioutil.WriteFile(filename, buf, os.ModePerm)
	return filename, err
}

// 发布音频至rtmp
func publishVoice(buf []byte, id string) error {
	mp3filename, err := saveVoice(buf, id)
	if err != nil {
		log4go.Error("live %s %q", id, err)
		return err
	}
	// 执行系统命令
	cmd := exec.Command(ffmpeg, "-i", mp3filename, "-acodec", "aac", "-f", "flv", rtmp+"/"+id)
	log4go.Trace("live %s execute %s", id, strings.Join(cmd.Args, " "))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log4go.Error("live %s execute %q", id, err)
		return err
	}
	return nil
}
