package util

import (
	"math/rand"
	"strconv"
	"time"
)


func getTimestamp() int {
	return int(time.Now().UnixNano() / 1000000)
}

func GetTimestamp() string {
	return strconv.Itoa(getTimestamp())
}

func GetR() string {
	return strconv.Itoa(-getTimestamp() / 1579)
}

func SleepSec(sec int) {
	time.Sleep(time.Duration(sec) * time.Second)
}

func GetRandomID(n int) string {
	rand.Seed(time.Now().Unix())
	return "e" + strconv.FormatFloat(rand.Float64(), 'f', n, 64)[2:]
}