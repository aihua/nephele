package util

import (
	"runtime"
	"strconv"
)

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
