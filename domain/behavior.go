package domain

type WatchAction struct {
	VideoID   uint
	SegmentID uint
	Event     string // pause|replay
	Duration  int
}

type WatchProgress struct {
	VideoID    uint
	CurrentSec int
}

type VideoProgress struct {
	VideoID         uint
	Title           string
	Status          string
	ProgressPercent int
	WatchTime       int
}

