package fdfs

import (
	"testing"
	//  "time"
)

func TestFixString(t *testing.T) {
	s := "helloworld"
	r := fixString(s, 5)
	if r != "hello" {
		t.Error("test fix string fail")
	}
	r = fixString(s, 11)
	if r != "helloworld"+string(0) {
		t.Error("test fix string fail")
	}
}
