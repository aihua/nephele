package util

import (
	"github.com/ctripcorp/cat.go"
)

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
