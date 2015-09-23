package business

import (
	_ "github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"testing"
)

func TestFDFSUpload(t *testing.T) {
	InitDBForTest()
	bts, err := ioutil.ReadFile("111.png")
	if err != nil {
		t.Error(err)
	}

	fdfs := FdfsStorage{}
	path, e := fdfs.Upload(bts, "png")
	if e.Err != nil {
		t.Error(e)
	}
	println("image storage path->>>>" + path)
}

func TestFdfsConvertFilePath(t *testing.T) {
	fdfs := FdfsStorage{Path: "/fd/hotel/group1/M01/6C/36/CgIZH1YAwemAUhVrAAMwrelaA5k766.png"}
	e := fdfs.ConvertFilePath()
	if e.Err != nil {
		t.Error(e.Err)
	}
	if fdfs.Path != "group1/M01/6C/36/CgIZH1YAwemAUhVrAAMwrelaA5k766.png" {
		t.Error("TestFdfsConvertFilePath fail")
	}
}

func TestNfsConvertFilePath(t *testing.T) {
	nfs := NfsStorage{Path: "/t1/headphoto/057/777/943/0b93f8268d5546308915f4f9fcaa9483.jpg"}
	e := nfs.ConvertFilePath()
	if e.Err != nil {
		t.Error(e.Err)
	}
	path := "/home/gct/target/"
	if nfs.Path != path+"headphoto/057/777/943/0b93f8268d5546308915f4f9fcaa9483.jpg" {
		t.Error("TestNfsConvertFilePath fail")
	}

	nfs1 := NfsStorage{Path: "/t1/headphoto/057/777/943/0b93f8268d5546308915f4f9fcaa9483.jpg"}
	e = nfs1.ConvertFilePath()
	if e.Err != nil {
		t.Error(e.Err)
	}
	if nfs1.Path != path+"headphoto/057/777/943/0b93f8268d5546308915f4f9fcaa9483.jpg" {
		t.Error("TestNfsConvertFilePath fail")
	}
}
