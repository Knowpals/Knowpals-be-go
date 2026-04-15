package behavior

import "time"

type RecordActionReq struct {
	VideoID   uint   `json:"video_id"`
	SegmentID uint   `json:"segment_id"`
	Event     string `json:"event"` //pause|replay
	Duration  int    `json:"duration"`
}

type UpdateProgressReq struct {
	VideoID    uint `json:"video_id"`
	CurrentSec int  `json:"current_sec"`
}

type VideoProgress struct {
	VideoID         uint       `json:"video_id"`
	Title           string     `json:"title"`
	Status          string     `json:"status"`           //finished | todo | expired
	ProgressPercent int        `json:"progress_percent"` //完成进度百分比
	WatchTime       int        `json:"watch_time"`       //总观看时长
	Deadline        *time.Time `json:"deadline"`
	CreatedAt       time.Time  `json:"created_at"`
}

type GetClassVideoProgressReq struct {
	ClassID uint   `uri:"class_id" binding:"required"`
	Status  string `uri:"status" binding:"required"`
}

type GetClassVideoProgressResp struct {
	ProgressList []VideoProgress `json:"progress_list"`
}
