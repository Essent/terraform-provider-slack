package tb

import "time"

func RandBool() bool {
	return time.Now().UnixNano()%2 == 0
}
