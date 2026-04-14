package dao

import (
	"context"

	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/repository/model"
	"gorm.io/gorm"
)

type StatisticDao interface {
	GetStudentVideoStat(ctx context.Context, studentID uint, videoID uint) (domain.StudentVideoStat, error)
	GetClassVideoStat(ctx context.Context, classID uint, videoID uint) (domain.ClassVideoStat, error)
}

type statisticDao struct {
	db *gorm.DB
}

func NewStatisticDao(db *gorm.DB) StatisticDao {
	return &statisticDao{db: db}
}

func (d *statisticDao) GetStudentVideoStat(ctx context.Context, studentID uint, videoID uint) (domain.StudentVideoStat, error) {
	// progress
	status := "todo"
	timeCost := 0
	var p model.StudentVideoProgress
	if err := d.db.WithContext(ctx).Where("user_id=? AND video_id=?", studentID, videoID).First(&p).Error; err == nil {
		status = p.Status
		timeCost = p.WatchDuration
	}

	// correct rate
	type agg struct {
		Total   int64 `gorm:"column:total"`
		Correct int64 `gorm:"column:correct"`
	}
	var a agg
	_ = d.db.WithContext(ctx).
		Model(&model.StudentAnswer{}).
		Select("count(*) as total, sum(case when is_correct then 1 else 0 end) as correct").
		Where("student_id=? AND video_id=?", studentID, videoID).
		Scan(&a).Error
	correctRate := 0.0
	if a.Total > 0 {
		correctRate = float64(a.Correct) / float64(a.Total)
	}

	// pause count
	type pauseAgg struct {
		Pause int64 `gorm:"column:pause"`
	}
	var pa pauseAgg
	_ = d.db.WithContext(ctx).
		Model(&model.StudentBehavior{}).
		Select("sum(pause_count) as pause").
		Where("student_id=? AND video_id=?", studentID, videoID).
		Scan(&pa).Error

	// weak knowledge by error rate (top 5)
	type weakRow struct {
		KnowledgeID uint
		Title       string
		Total       int64
		Wrong       int64
	}
	var wr []weakRow
	_ = d.db.WithContext(ctx).
		Table("student_answers sa").
		Select("kp.id as knowledge_id, kp.title as title, count(*) as total, sum(case when sa.is_correct=0 then 1 else 0 end) as wrong").
		Joins("join question_to_knowledge qk on qk.question_id = sa.question_id").
		Joins("join knowledge_points kp on kp.id = qk.knowledge_id").
		Where("sa.student_id=? AND sa.video_id=?", studentID, videoID).
		Group("kp.id,kp.title").
		Order("wrong desc").
		Limit(5).
		Scan(&wr).Error
	weak := make([]domain.StudentWeakKnowledgePoint, 0, len(wr))
	for _, r := range wr {
		score := 1.0
		if r.Total > 0 {
			score = 1.0 - float64(r.Wrong)/float64(r.Total)
		}
		weak = append(weak, domain.StudentWeakKnowledgePoint{
			KnowledgeID: r.KnowledgeID,
			Title:       r.Title,
			MasterScore: score,
		})
	}

	return domain.StudentVideoStat{
		Status:              status,
		CorrectRate:         correctRate,
		TimeCost:            timeCost,
		PauseCount:          int(pa.Pause),
		WeakKnowledgePoints: weak,
	}, nil
}

func (d *statisticDao) GetClassVideoStat(ctx context.Context, classID uint, videoID uint) (domain.ClassVideoStat, error) {
	// student list
	var studentIDs []uint
	if err := d.db.WithContext(ctx).
		Model(&model.ClassStudent{}).
		Where("class_id = ?", classID).
		Pluck("student_id", &studentIDs).Error; err != nil {
		return domain.ClassVideoStat{}, err
	}
	if len(studentIDs) == 0 {
		return domain.ClassVideoStat{}, nil
	}

	// overview aggregates
	type ansAgg struct {
		Total   int64
		Correct int64
	}
	var aa ansAgg
	_ = d.db.WithContext(ctx).
		Model(&model.StudentAnswer{}).
		Select("count(*) as total, sum(case when is_correct then 1 else 0 end) as correct").
		Where("video_id=? AND student_id IN ?", videoID, studentIDs).
		Scan(&aa).Error
	avgCorrect := 0.0
	if aa.Total > 0 {
		avgCorrect = float64(aa.Correct) / float64(aa.Total)
	}

	type progAgg struct {
		AvgWatch float64
		Finished int64
		Total    int64
	}
	var pg progAgg
	_ = d.db.WithContext(ctx).
		Table("student_video_progresses").
		Select("avg(watch_duration) as avg_watch, sum(case when status='finished' then 1 else 0 end) as finished, count(*) as total").
		Where("video_id=? AND user_id IN ?", videoID, studentIDs).
		Scan(&pg).Error
	completeRate := 0.0
	if pg.Total > 0 {
		completeRate = float64(pg.Finished) / float64(pg.Total)
	}

	type pauseAgg struct {
		Pause int64
	}
	var pAgg pauseAgg
	_ = d.db.WithContext(ctx).
		Model(&model.StudentBehavior{}).
		Select("sum(pause_count) as pause").
		Where("video_id=? AND student_id IN ?", videoID, studentIDs).
		Scan(&pAgg).Error

	// top questions error rate
	type qRow struct {
		QuestionID uint
		Content    string
		Total      int64
		Wrong      int64
	}
	var qrows []qRow
	_ = d.db.WithContext(ctx).
		Table("student_answers sa").
		Select("sa.question_id as question_id, q.content as content, count(*) as total, sum(case when sa.is_correct=0 then 1 else 0 end) as wrong").
		Joins("join questions q on q.id = sa.question_id").
		Where("sa.video_id=? AND sa.student_id IN ?", videoID, studentIDs).
		Group("sa.question_id,q.content").
		Order("wrong desc").
		Limit(10).
		Scan(&qrows).Error
	topQs := make([]domain.ClassQuestionStat, 0, len(qrows))
	for _, r := range qrows {
		er := 0.0
		if r.Total > 0 {
			er = float64(r.Wrong) / float64(r.Total)
		}
		topQs = append(topQs, domain.ClassQuestionStat{QuestionID: r.QuestionID, Content: r.Content, ErrorRate: er})
	}

	// top pause/replay segments
	type segRow struct {
		SegmentID uint
		PauseSum  int64
		ReplaySum int64
	}
	var segs []segRow
	_ = d.db.WithContext(ctx).
		Model(&model.StudentBehavior{}).
		Select("segment_id, sum(pause_count) as pause_sum, sum(replay_count) as replay_sum").
		Where("video_id=? AND student_id IN ?", videoID, studentIDs).
		Group("segment_id").
		Order("pause_sum desc").
		Limit(10).
		Scan(&segs).Error
	var segModels []model.Segment
	segIDs := make([]uint, 0, len(segs))
	for _, s := range segs {
		segIDs = append(segIDs, s.SegmentID)
	}
	if len(segIDs) > 0 {
		_ = d.db.WithContext(ctx).Where("id IN ? AND video_id=?", segIDs, videoID).Find(&segModels).Error
	}
	segMap := make(map[uint]model.Segment, len(segModels))
	for _, s := range segModels {
		segMap[s.ID] = s
	}
	pauseActs := make([]domain.ClassPauseAction, 0, len(segs))
	replayActs := make([]domain.ClassReplayAction, 0, len(segs))
	for _, s := range segs {
		seg := segMap[s.SegmentID]
		pauseActs = append(pauseActs, domain.ClassPauseAction{SegmentID: s.SegmentID, Start: seg.Start, End: seg.End, PauseCount: int(s.PauseSum)})
		replayActs = append(replayActs, domain.ClassReplayAction{SegmentID: s.SegmentID, Start: seg.Start, End: seg.End, ReplayCount: int(s.ReplaySum)})
	}

	return domain.ClassVideoStat{
		Overview: domain.ClassOverview{
			AverageCorrectRate: avgCorrect,
			AverageTimeCost:    int(pg.AvgWatch),
			CompleteRate:       completeRate,
			TotalPauseCount:    int(pAgg.Pause),
		},
		TopQuestions:    topQs,
		TopPauseAction:  pauseActs,
		TopReplayAction: replayActs,
	}, nil
}
