package review

type StartReviewReq struct {
	VideoID uint `uri:"video_id" binding:"required"`
}

type PublishReq struct {
	VideoID uint `uri:"video_id" binding:"required"`
}

