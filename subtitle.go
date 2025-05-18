package subtitle

import "time"

type Subtitle struct {
	Index int
	Start time.Time
	End   time.Time
	Text  string
}
