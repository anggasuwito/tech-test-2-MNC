package utils

import "time"

const (
	DefaultLoc = "Asia/Jakarta"
)

var (
	loc, _ = time.LoadLocation(DefaultLoc)
)

func TimeNow() time.Time {
	return time.Now().In(loc)
}
