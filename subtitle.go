package subtitle

import "time"

type Cue struct {
	Index int
	Start time.Time
	End   time.Time
	Text  string
}
