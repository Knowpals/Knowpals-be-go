package domain

type StudentWeakKnowledgePoint struct {
	KnowledgeID uint
	Title       string
	MasterScore float64
}

type StudentVideoStat struct {
	Status              string
	CorrectRate         float64
	TimeCost            int
	PauseCount          int
	WeakKnowledgePoints []StudentWeakKnowledgePoint
}

type ClassOverview struct {
	AverageCorrectRate float64
	AverageTimeCost    int
	CompleteRate       float64
	TotalPauseCount    int
}

type ClassQuestionStat struct {
	QuestionID uint
	Content    string
	ErrorRate  float64
}

type ClassPauseAction struct {
	SegmentID  uint
	Start      int
	End        int
	PauseCount int
}

type ClassReplayAction struct {
	SegmentID   uint
	Start       int
	End         int
	ReplayCount int
}

type ClassVideoStat struct {
	Overview        ClassOverview
	TopQuestions    []ClassQuestionStat
	TopPauseAction  []ClassPauseAction
	TopReplayAction []ClassReplayAction
}

