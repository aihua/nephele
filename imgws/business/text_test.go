package business

import "testing"
import "io/ioutil"

func TestTextImage(t *testing.T) {
	bts := getTextImage("你好", 20)
	ioutil.WriteFile("1.png", bts, 0644)
}
