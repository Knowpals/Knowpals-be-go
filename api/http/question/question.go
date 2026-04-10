package question

type GenerateVideoExerciseReq struct {
	VideoID uint `json:"video_id"`
}

type KnowledgePoint struct {
	KnowledgeID uint   `json:"knowledge_id"`
	Title       string `json:"title"`
}

type Question struct {
	ID              uint             `json:"id"`
	Type            string           `json:"type"`
	Content         string           `json:"content"`
	Options         string           `json:"options"`
	Answer          string           `json:"answer"`
	Analysis        string           `json:"analysis"`
	SegmentID       *uint            `json:"segment_id"`
	KnowledgePoints []KnowledgePoint `json:"knowledge_points"`
}

type GenerateVideoExerciseResp struct {
	Questions []Question `json:"questions"`
}

type StudentAnswer struct {
	QuestionID uint   `json:"question_id"`
	Answer     string `json:"answer"`
	TimeCost   int    `json:"time_cost"`
}

type Result struct {
	QuestionID uint   `json:"question_id"`
	IsCorrect  bool   `json:"is_correct"`
	Answer     string `json:"answer"`
	Analysis   string `json:"analysis"`
}

type AnswerQuestionReq struct {
	VideoID        uint            `json:"video_id"`
	StudentAnswers []StudentAnswer `json:"studentanswers"`
}

type AnswerQuestionResp struct {
	Results []Result `json:"results"`
}
