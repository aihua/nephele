package util

import (
	"archive/zip"
	"bytes"
	"net"
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

type MemoryWriter struct {
	Content []byte
}

func (this *MemoryWriter) Write(p []byte) (n int, err error) {
	this.Content = p
	n = len(p)
	err = nil
	return
}
func Zip(files map[string][]byte) ([]byte, Error) {
	memory := new(MemoryWriter)
	wzip := zip.NewWriter(memory)
	defer func() {
		if err := wzip.Close(); err != nil {
			//todo
		}
	}()
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
	return memory.Content, Error{}
}
