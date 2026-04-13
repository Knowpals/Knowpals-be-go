package domain

type PipelineJob struct {
	JobID        string
	VideoID      uint
	Status       string
	CurrentStage int
	TotalStage   int
}

type PipelineJobStage struct {
	ID         uint
	JobID      string
	Stage      string
	Status     string
	RetryCount int
	Output     string
}
