package business

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func TestFDFSUpload(t *testing.T) {
	InitDBForTest()

}
