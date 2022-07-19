package store

import "time"

func ExportSetTimeNow(t time.Time) {
	timeNow = func() *time.Time {
		return &t
	}
}
