package stats

import (
	"runtime"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type MemStats struct {
	Alloc      string
	TotalAlloc string
	Sys        string
}

func GetMemStats() MemStats {
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)

	return MemStats{
		Alloc:      Bytes(memStat.Alloc),
		TotalAlloc: Bytes(memStat.TotalAlloc),
		Sys:        Bytes(memStat.Sys),
	}
}

// Bytes represent number only MB
func Bytes(i uint64) string {
	p := message.NewPrinter(language.English)

	return p.Sprintf("%d %s", i/1000, "MB")
}
