package models

import (
	"strconv"
	"testing"
)

func TestParseName(t *testing.T) {
	imgname := "100208000000000117D51"
	imgIndex := ImageIndex{}
	e := imgIndex.ParseName(imgname)
	if e.Err != nil {
		t.Error(e.Err)
	}
	if imgIndex.Idx != 11 {
		t.Error("idx is invalid!" + strconv.Itoa(int(imgIndex.Idx)))
	}
	if imgIndex.Channel != "10" {
		t.Error("channel is invalid!" + imgIndex.Channel)
	}
	if imgIndex.TableZone != 2 {
		t.Error("tablezone is invalid!" + strconv.Itoa(imgIndex.TableZone))
	}
	if imgIndex.PartitionKey != 8 {
		t.Error("partitionkey is invalid!" + strconv.Itoa(int(imgIndex.PartitionKey)))
	}

	if imgIndex.Version != "0" {
		t.Error("version is invalid!" + imgIndex.Version)
	}
}

func TestGetStorage(t *testing.T) {

}
