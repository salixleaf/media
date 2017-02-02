package video

import (
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Video struct {
	Path     string
	Length   float64 // 时长(s)
	Bitrate  float64 // 播放速率(kb/s)
	Size     int64   // 文件大小(byte)
	Width    int64   // 视频分辨率宽度
	Height   int64   // 视频分辨率高度
	Fps      float64 // 视频帧率(帧/s)
	Vbitrate float64 // 视频比特率(kb/s)
	Abitrate float64 // 音频比特率(kb/s)
	Ahz      float64 // 音频采集率(Hz)
}

// Stat 通过调用 ffmpeg命令 使用正则获取视频信息,
// 部分视频无法正常获取时长或比特率等，则使用0表示;
// 如果多个属性无法获取，则可能是正则匹配不全，
// 请手动执行 ffmpeg -i file_path 参照输出信息来确认问题
func Stat(video_path string) *Video {

	cmd := exec.Command("ffmpeg", "-i", video_path)
	r, _ := cmd.CombinedOutput()

	// sample1
	//  Duration: 00:00:00.00, start: 0.000000, bitrate: N/A
	//    Stream #0:0: Video: rv40 (RV40 / 0x30345652), yuv420p, 640x480, 25 fps, 25 tbr, 1k tbn, 1k tbc
	//    Stream #0:1: Audio: cook (cook / 0x6B6F6F63), 44100 Hz, mono, fltp, 64 kb/s

	// sample2
	//  Duration: 00:10:23.13, start: 0.000000, bitrate: 1741 kb/s
	//    Stream #0:0: Video: h264 (High) (H264 / 0x34363248), yuv420p(progressive), 352x288 [SAR 1:1 DAR 11:9], 1604 kb/s, 30 fps, 30 tbr, 30 tbn, 60 tbc
	//    Stream #0:1: Audio: mp3 (U[0][0][0] / 0x0055), 44100 Hz, stereo, s16p, 128 kb/s

	// sample3
	//  Duration: 00:17:57.43, start: 0.000000, bitrate: 383 kb/s
	//    Stream regexp.MustCompile(`.*Duration:\s(.*?),.*bitrate:\s(\S+)`)#0:0: Audio: cook (cook / 0x6B6F6F63), 44100 Hz, stereo, fltp, 64 kb/s
	//    Stream #0:1: Video: rv40 (RV40 / 0x30345652), yuv420p, 640x480, 308 kb/s, 23.98 fps, 23.98 tbr, 1k tbn, 1k tbc

	str_r := string([]byte(r))
	length, bitrate := parse_duration(str_r)
	width, height, v_bitrate, fps := parse_video(str_r)
	a_hz, a_bitrate := parse_audio(str_r)

	v := &Video{
		Path:     video_path,
		Length:   length,
		Bitrate:  bitrate,
		Size:     get_size(video_path),
		Width:    width,
		Height:   height,
		Fps:      fps,
		Vbitrate: v_bitrate,
		Abitrate: a_bitrate,
		Ahz:      a_hz,
	}

	return v
}

func get_size(path string) int64 {
	file, _ := os.Stat(path)
	return file.Size()
}

var reg_duration = regexp.MustCompile(`.*Duration:\s(.*?),.*bitrate:\s(\S+)`)

// 解析Duration行
func parse_duration(str string) (float64, float64) {
	s := reg_duration.FindStringSubmatch(str)
	if len(s) != 3 {
		return 0, 0
	}

	t := strings.Split(s[1], ":")
	length := atof64(t[0])*3600 + atof64(t[1])*60 + atof64(t[2])
	return length, atof64(s[2])
}

var reg_video = regexp.MustCompile(`Stream.*Video.*\s(\d+)x(\d+)(?:.*?(\S+)\skb/s)?.*?(\S+)\sfps`)

// 解析Video行
func parse_video(str string) (int64, int64, float64, float64) {
	s := reg_video.FindStringSubmatch(str)

	if len(s) != 5 {
		return 0, 0, 0, 0
	}
	return atoi64(s[1]), atoi64(s[2]), atof64(s[3]), atof64(s[4])
}

var reg_audio = regexp.MustCompile(`Stream.*Audio.*?(\d+)\sHz.*\s(\S+)\skb/s`)

// 解析Audio行
func parse_audio(str string) (float64, float64) {
	s := reg_audio.FindStringSubmatch(str)

	if len(s) != 3 {
		return 0, 0
	}
	return atof64(s[1]), atof64(s[2])
}

func atoi64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func atof64(s string) float64 {
	i, _ := strconv.ParseFloat(s, 64)
	return i
}
