package util

import (
	"fmt"
	"strconv"
	"testing"
)

func Test_Cover(t *testing.T) {
	s := "1"
	s1 := Cover(s, "0", 5)
	if s1 != "00001" {
		t.Error("cover error |" + s1)
	}
}
