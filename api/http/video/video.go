package video

import (
	"time"

	"github.com/Knowpals/Knowpals-be-go/api/http/question"
)

type UploadVideoReq struct {
	Title    string `form:"title" binding:"required"`
	Deadline string `form:"deadline" binding:"required"`
}

type UploadVideoResp struct {
	VideoID  uint   `json:"video_id"`
	JobID    string `json:"job_id"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	Duration int    `json:"duration"`
}

type GetVideoDetailReq struct {
	VideoID uint `uri:"video_id" binding:"required"`
}

type Segment struct {
	ID        uint              `json:"id"`
	SegmentID string            `json:"segment_id"`
	Start     int               `json:"start"`
	End       int               `json:"end"`
	Text      string            `json:"text"`
	Question  question.Question `json:"question"`
}

type GetVideoDetailResp struct {
	VideoID   uint                 `json:"video_id"`
	Title     string               `json:"title"`
	Duration  int                  `json:"duration"`
	Segments  []Segment            `json:"segments"`
	Knowledge []KnowledgePointResp `json:"knowledge"`
	Questions []question.Question  `json:"questions"`
}

type KnowledgePointResp struct {
	ID          uint   `json:"id"`
	KnowledgeID string `json:"knowledge_id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
}

type PostVideoToClassReq struct {
	VideoID   uint   `json:"video_id"`
	ClassList []uint `json:"class_list"`
}

type GetClassVideosReq struct {
	ClassID uint `uri:"class_id" binding:"required"`
}

type VideoTask struct {
	VideoID   uint      `json:"video_id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Deadline  time.Time `json:"deadline"`
}

type GetClassVideosResp struct {
	VideoTasks []VideoTask `json:"video_tasks"`
}

type GetTaskUploadingProcessReq struct {
	JobID string `uri:"job_id" binding:"required"`
}

type GetTaskUploadingProcessResp struct {
	JobID  string `json:"job_id"`
	Status string `json:"status"`
	Stage  string `json:"stage"`
}

type GetMyUploadedVideosResp struct {
	Videos []VideoTask `json:"videos"`
}
