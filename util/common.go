package util

import "time"

func SliceStringContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func NowInt() int {
	return int(time.Now().Unix())
}