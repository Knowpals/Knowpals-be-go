package statistic

type GetStudentStatReq struct {
	VideoID uint `uri:"video_id" binding:"required"`
}

// 知识点的分数根据公式算(题目*0.6+暂停*0.4)
type KnowledgePoint struct {
	KnowledgeID uint    `json:"knowledge_id"`
	Title       string  `json:"title"`
	MasterScore float64 `json:"master_score"`
}

type GetStudentStatResp struct {
	//任务完成情况：finished | in_progress | todo | expired
	Status string `json:"status"`
	//题目正确率
	CorrectRate float64 `json:"correct_rate"`
	//视频观看时长
	TimeCost int `json:"time_cost"`
	//暂停总时长
	PauseCount int `json:"pause_count"`
	//薄弱知识点（薄弱知识点得分，要提前确定一个阈值）
	WeakKnowledgePoints []KnowledgePoint `json:"knowledge_points"`
}

type Overview struct {
	//平均答题正确率
	AverageCorrectRate float64 `json:"average_correct_rate"`
	//平均观看时长
	AverageTimeCost int `json:"average_time_cost"`
	//任务完成率
	CompleteRate float64 `json:"complete_rate"`
	//总暂停时长
	TotalPauseCount int `json:"total_pause_count"`
}

type GetClassStatReq struct {
	ClassID uint `json:"class_id"`
	VideoID uint `json:"video_id"`
}

type AverageKnowledgePoint struct {
	KnowledgeID        uint    `json:"knowledge_id"`
	Title              string  `json:"title"`
	AverageMasterScore float64 `json:"average_master_score"`
	//低于平均掌握程度分数的比例
	WeakRate float64 `json:"weak_rate"`
}

type PauseAction struct {
	SegmentID  uint `json:"segment_id"`
	Start      int  `json:"start"`
	End        int  `json:"end"`
	PauseCount int  `json:"pause_count"`
}
type ReplayAction struct {
	SegmentID   uint `json:"segment_id"`
	Start       int  `json:"start"`
	End         int  `json:"end"`
	ReplayCount int  `json:"replay_count"`
}

type QuestionStat struct {
	QuestionID uint    `json:"question_id"`
	Content    string  `json:"content"`
	ErrorRate  float64 `json:"error_rate"`
}

type GetClassStatResp struct {
	Overview           Overview                `json:"overview"`
	TopQuestions       []QuestionStat          `json:"top_questions"`
	WeakKnowledgePoint []AverageKnowledgePoint `json:"weak_knowledge_point"`
	TopPauseAction     []PauseAction           `json:"top_pause_action"`
	TopReplayAction    []ReplayAction          `json:"top_replay_action"`
}
