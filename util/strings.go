package util

import (
	"archive/zip"
	"bytes"
	cat "github.com/ctripcorp/cat.go"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//cover 字符串补长度
func Cover(s, converV string, length int) string {
	currentLen := len(s)
	for i := 0; i < length-currentLen; i++ {
		s = converV + s
	}
	return s
}

func JoinString(args ...string) string {
	var buf bytes.Buffer
	for _, v := range args {
		buf.WriteString(v)
	}
	return buf.String()
}

func GetPartitionKey(t time.Time) int16 {
	return int16((t.Year()-2015)*12 + int(t.Month()) - 1)
}

var localIP string = ""

func GetIP() string {
	if localIP != "" {
		return localIP
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		add := strings.Split(addr.String(), "/")[0]
		if add == "127.0.0.1" || add == "::1" {
			continue
		}
		first := strings.Split(add, ".")[0]
		if _, err := strconv.Atoi(first); err == nil {
			localIP = add
			return add
		}
	}
	return ""
}

func GetClientIP(req *http.Request) string {
	addr := req.Header.Get("X-Real-IP")
	if addr == "" {
		addr = req.Header.Get("X-Forwarded-For")
		if addr == "" {
			addr = req.RemoteAddr
		}
	}
	return addr
}

func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0
	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length
	if start > end {
		start, end = end, start
	}
	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

func Zip(files map[string][]byte) ([]byte, Error) {
	buffer := bytes.NewBuffer(nil)
	wzip := zip.NewWriter(buffer)

	errtype := "ZipFail"
	for fileName, content := range files {
		f, err := wzip.Create(fileName)
		if err != nil {
			return []byte{}, Error{IsNormal: false, Err: err, Type: errtype}
		}
		_, err = f.Write(content)
		if err != nil {
			return []byte{}, Error{IsNormal: false, Err: err, Type: errtype}
		}
	}
	if err := wzip.Close(); err != nil {
		return []byte{}, Error{IsNormal: false, Err: err, Type: errtype}
	}
	return buffer.Bytes(), Error{}
}

func LogErrorEvent(cat cat.Cat, name string, err string) {
	if cat == nil {
		return
	}
	event := cat.NewEvent("Error", name)
	event.AddData("detail", err)
	event.SetStatus("ERROR")
	event.Complete()
}

func LogEvent(cat cat.Cat, title string, name string, data map[string]string) {
	if cat == nil {
		return
	}
	event := cat.NewEvent(title, name)
	if data != nil {
		for k, v := range data {
			event.AddData(k, v)
		}
	}
	event.SetStatus("0")
	event.Complete()
}

// Alloc        uint64      bytes allocated and still in use // 已分配且仍在使用的字节数
// 	TotalAlloc   uint64      // bytes allocated (even if freed) // 已分配（包括已释放的）字节数
// 	Sys          uint64      // bytes obtained from system (sum of XxxSys below) // 从系统中获取的字节数（应当为下面 XxxSys 之和）
// 	Mallocs      uint64      // number of mallocs // malloc 数
// 	Frees        uint64      // number of frees // free 数
// 	HeapAlloc    uint64      // bytes allocated and still in use // 已分配且仍在使用的字节数
// 	HeapSys      uint64      // bytes obtained from system // 从系统中获取的字节数
// 	HeapIdle     uint64      // bytes in idle spans // 空闲区间的字节数
// 	HeapInuse    uint64      // bytes in non-idle span // 非空闲区间的字节数
// 	HeapReleased uint64      // bytes released to the OS // 释放给OS的字节数
// 	HeapObjects  uint64      // total number of allocated objects// 已分配对象的总数
// 	OtherSys     uint64      // other system allocations // 其它系统分配
// 	NextGC       uint64      // next run in HeapAlloc time (bytes) // 下次运行的 HeapAlloc 时间（字节）
// 	LastGC       uint64      // last run in absolute time (ns) // 上次运行的绝对时间（纳秒 ns）
// 	PauseNs      [256]uint64 // circular buffer of recent GC pause times, most recent at [(NumGC+255)%256]
// 	NumGC        uint32

func GetStatus() map[string]string {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	var second uint64 = 1000000000
	var memory uint64 = 1024 * 1024
	return map[string]string{"Alloc": strconv.FormatUint(mem.Alloc/memory, 10),
		"TotalAlloc":   strconv.FormatUint(mem.TotalAlloc/memory, 10),
		"Sys":          strconv.FormatUint(mem.Sys/memory, 10),
		"Mallocs":      strconv.FormatUint(mem.Mallocs, 10),
		"Frees":        strconv.FormatUint(mem.Frees, 10),
		"HeapAlloc":    strconv.FormatUint(mem.HeapAlloc/memory, 10),
		"HeapSys":      strconv.FormatUint(mem.HeapSys/memory, 10),
		"HeapIdle":     strconv.FormatUint(mem.HeapIdle/memory, 10),
		"HeapInuse":    strconv.FormatUint(mem.HeapInuse/memory, 10),
		"HeapReleased": strconv.FormatUint(mem.HeapReleased/memory, 10),
		"HeapObjects":  strconv.FormatUint(mem.HeapObjects, 10),
		"OtherSys":     strconv.FormatUint(mem.OtherSys, 10),
		"NextGC":       strconv.FormatUint(mem.NextGC/second, 10),
		"LastGC":       strconv.FormatUint(mem.LastGC/second, 10),
		"PauseNs":      strconv.FormatUint(mem.PauseNs[(mem.NumGC+255)%256]/second, 10),
		"NumGC":        strconv.Itoa(int(mem.NumGC)),
	}
}

func GetImageSizeDistribution(size int) string {
	switch {
	case size < 0:
		return "<0"
	case size == 0:
		return "0"
	case size > 0 && size <= 512*1024:
		return "1~512KB"
	case size > 512*1024 && size <= 1024*1024:
		return "512~1024KB"
	case size > 1024*1024 && size <= 2*1024*1024:
		return "1~2M"
	case size > 2*1024*1024 && size <= 4*1024*1024:
		return "2~4M"
	case size > 4*1024*1024 && size <= 6*1024*1024:
		return "4~6M"
	case size > 6*1024*1024 && size <= 10*1024*1024:
		return "6~10M"
	case size > 10*1024*1024 && size <= 20*1024*1024:
		return "10~20M"
	case size > 20*1024*1024 && size <= 30*1024*1024:
		return "20~30M"
	default:
		return ">30M"
	}
}
