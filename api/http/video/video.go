package video

import (
	"time"

	"github.com/Knowpals/Knowpals-be-go/api/http/question"
)

type UploadVideoReq struct {
	Title string `form:"title" binding:"required"`
}

type UploadVideoResp struct {
	VideoID  uint   `json:"video_id"`
	JobID    string `json:"job_id"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	Duration int    `json:"duration"`
}

type GetVideoDetailReq struct {
	VideoID uint `json:"video_id" binding:"required"`
}

type Segment struct {
	ID    uint `json:"id"`
	Start int  `json:"start"`
	End   int  `json:"end"`
}

type GetVideoDetailResp struct {
	Segments  []Segment           `json:"segments"`
	Questions []question.Question `json:"questions"`
}

type PostVideoToClassReq struct {
	VideoID   uint   `json:"video_id"`
	ClassList []uint `json:"class_list"`
}

type GetClassVideosReq struct {
	ClassID uint `json:"class_id"`
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
