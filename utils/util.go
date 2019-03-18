package utils

import (
	"github.com/cihub/seelog"
	"strings"
	"time"
)

func StrToTime(timeStr string) time.Time {
	times := strings.Replace(timeStr,"T"," ",-1)
	times = strings.Split(times, ".")[0]
	parse, err := time.Parse("2006-01-02 15:04:05", times)
	if err!=nil {
		seelog.Info(err)
	}
	return parse
}
