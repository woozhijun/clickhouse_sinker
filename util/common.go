package util

import (
	"github.com/wswz/go_commons/log"
	"time"
)

const LayoutDatetime = "2006-01-02 15:04:05"
const LayoutDate = "2006-01-02"

func StringContains(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

func StringParseToDate(str string, layout string) time.Time {
	p, err := time.Parse(layout, str)
	if err != nil {
		log.Error("StringParseToDate error. value: " + str, err)
	}
	return p
}
