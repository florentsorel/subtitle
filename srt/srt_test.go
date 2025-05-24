package srt

import (
	"testing"
	"time"

	"github.com/florentsorel/subtitle"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input             string
		expectedSubtitles []subtitle.Cue
	}{
		{
			input: `1
00:00:01,123 --> 00:00:02,456
First line
Second line

2
00:00:03,123 --> 00:00:04,456
Single line

3
00:00:05,123 --> 00:00:06,456
Last cue
Last line`,
			expectedSubtitles: []subtitle.Cue{
				{
					Index: 1,
					Start: parseTime("00:00:01,123"),
					End:   parseTime("00:00:02,456"),
					Text:  "First line\nSecond line",
				},
				{
					Index: 2,
					Start: parseTime("00:00:03,123"),
					End:   parseTime("00:00:04,456"),
					Text:  "Single line",
				},
				{
					Index: 3,
					Start: parseTime("00:00:05,123"),
					End:   parseTime("00:00:06,456"),
					Text:  "Last cue\nLast line",
				},
			},
		},
	}

	for i, tt := range tests {
		subtitles, err := Parse(tt.input)

		if err != nil {
			t.Fatalf("test[%d] - unexpected error: %v", i, err)
		}

		if len(subtitles) != len(tt.expectedSubtitles) {
			t.Fatalf("test[%d] - wrong number of subtitles. expected=%d, got=%d",
				i, len(tt.expectedSubtitles), len(subtitles))
		}

		for j, sub := range subtitles {
			expected := tt.expectedSubtitles[j]

			if sub.Index != expected.Index {
				t.Fatalf("test[%d][%d] - wrong index. expected=%d, got=%d",
					i, j, expected.Index, sub.Index)
			}

			if !sub.Start.Equal(expected.Start) {
				t.Fatalf("test[%d][%d] - wrong start time. expected=%v, got=%v",
					i, j, expected.Start, sub.Start)
			}

			if !sub.End.Equal(expected.End) {
				t.Fatalf("test[%d][%d] - wrong end time. expected=%v, got=%v",
					i, j, expected.End, sub.End)
			}

			if sub.Text != expected.Text {
				t.Fatalf("test[%d][%d] - wrong text. expected=%q, got=%q",
					i, j, expected.Text, sub.Text)
			}
		}
	}
}

func parseTime(s string) time.Time {
	normalized := s
	if len(s) > 0 {
		normalized = s[:8] + "." + s[9:]
	}
	t, _ := time.Parse("15:04:05.000", normalized)
	return t
}
