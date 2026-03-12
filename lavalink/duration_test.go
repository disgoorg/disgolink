package lavalink

import (
	"testing"
)

func TestDuration_Milliseconds(t *testing.T) {
	milliseconds := Millisecond * 2
	equal(t, int64(2), milliseconds.Milliseconds())
}

func TestDuration_MillisecondsPart(t *testing.T) {
	seconds := Duration(2345)
	equal(t, int64(345), seconds.MillisecondsPart())
}

func TestDuration_Seconds(t *testing.T) {
	seconds := Second * 2
	equal(t, int64(2), seconds.Seconds())
}

func TestDuration_SecondsPart(t *testing.T) {
	seconds := Duration(2345)
	equal(t, int64(2), seconds.SecondsPart())
}

func TestDuration_Minutes(t *testing.T) {
	minutes := Minute * 2
	equal(t, int64(2), minutes.Minutes())
}

func TestDuration_MinutesPart(t *testing.T) {
	minutes := Duration(123456)
	equal(t, int64(2), minutes.MinutesPart())
}

func TestDuration_Hours(t *testing.T) {
	hours := Hour * 2
	equal(t, int64(2), hours.Hours())
}

func TestDuration_HoursPart(t *testing.T) {
	hours := Duration(7234567)
	equal(t, int64(2), hours.HoursPart())
}

func TestDuration_Days(t *testing.T) {
	days := Day * 2
	equal(t, int64(2), days.Days())
}

func equal(t *testing.T, expected int64, actual int64) {
	if expected != actual {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}
