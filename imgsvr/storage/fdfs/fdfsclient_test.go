package fdfs

import (
	cat "github.com/ctripcorp/nephele/Godeps/_workspace/src/github.com/ctripcorp/cat.go"
	"io/ioutil"
	"os"
	"testing"
	//	"time"
)

func TestDownloadToBuffer(t *testing.T) {
	trackerIp := []string{"10.2.25.31"}
	trackerPort := "22122"
	fdfsClient, err := NewFdfsClient(trackerIp, trackerPort)
	if err != nil {
		t.Error(err)
	}
	catInstance := cat.Instance()
	b, err := fdfsClient.DownloadToBuffer("group1/M00/18/91/CgIZH1RyormAP7xTAAA1U2y_hqk858.jpg", catInstance)
	if err != nil {
		t.Error(err)
	}
	ioutil.WriteFile("download.jpg", b, os.ModePerm)

}

func TestUploadByBuffer(t *testing.T) {
	trackerIp := []string{"10.2.25.31"}
	trackerPort := "22122"
	fdfsClient, err := NewFdfsClient(trackerIp, trackerPort)
	if err != nil {
		t.Error(err)
	}
	b, err := ioutil.ReadFile("1.jpg")
	if err != nil {
		t.Error(err)
	}
	fileId, err := fdfsClient.UploadByBuffer("group1", b, "jpg")
	if err != nil {
		t.Error(err)
	}
	println(fileId)
}

func TestUploadSlaveByBuffer(t *testing.T) {
	trackerIp := []string{"10.2.25.31"}
	trackerPort := "22122"
	fdfsClient, err := NewFdfsClient(trackerIp, trackerPort)
	if err != nil {
		t.Error(err)
	}
	b, err := ioutil.ReadFile("1.jpg")
	if err != nil {
		t.Error(err)
	}
	fileId, err := fdfsClient.UploadSlaveByBuffer(b, "group1/M00/6B/C3/CgIZH1X5ViWAMEtRAAxHnmB5PVo920.jpg", "yyy1", "jpg")
	if err != nil {
		t.Error(err)
	}
	println(fileId)
}

func TestDeleteFile(t *testing.T) {
	trackerIp := []string{"10.2.25.31"}
	trackerPort := "22122"
	fdfsClient, err := NewFdfsClient(trackerIp, trackerPort)
	if err != nil {
		t.Error(err)
	}
	fdfsClient.DeleteFile("group1/M00/6B/C3/CgIZH1X5ViWAMEtRAAxHnmB5PVo920yyy1.jpg")
	if err != nil {
		t.Error(err)
	}
}
