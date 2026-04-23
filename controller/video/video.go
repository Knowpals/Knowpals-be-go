package video

import (
	errors1 "errors"
	"fmt"
	"mime/multipart"

	"github.com/Knowpals/Knowpals-be-go/api/http"
	"github.com/Knowpals/Knowpals-be-go/api/http/question"
	"github.com/Knowpals/Knowpals-be-go/api/http/video"
	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/errors"
	"github.com/Knowpals/Knowpals-be-go/infra/cos"
	"github.com/Knowpals/Knowpals-be-go/pkg/ginx"
	"github.com/Knowpals/Knowpals-be-go/pkg/ijwt"
	"github.com/Knowpals/Knowpals-be-go/service/behavior"
	"github.com/Knowpals/Knowpals-be-go/service/pipeline"
	video2 "github.com/Knowpals/Knowpals-be-go/service/video"
	"github.com/Knowpals/Knowpals-be-go/tool"
	"github.com/gin-gonic/gin"
)

type VideoController interface {
	// UploadVideo 老师上传视频
	UploadVideo(c *gin.Context, req video.UploadVideoReq, file multipart.File, fileHeader *multipart.FileHeader, claim ijwt.UserClaim) (http.Response, error)
	// GetTaskUploadingProcess 获取视频上传进度
	GetTaskUploadingProcess(c *gin.Context, req video.GetTaskUploadingProcessReq) (http.Response, error)
	// GetVideoDetail 获取视频任务详情
	GetVideoDetail(c *gin.Context, req video.GetVideoDetailReq) (http.Response, error)
	// PostVideoToClass 给指定班级发送视频任务
	PostVideoToClass(c *gin.Context, req video.PostVideoToClassReq) (http.Response, error)
	// GetClassVideoTasks 学生获取班级任务信息
	GetClassVideoTasks(c *gin.Context, req video.GetClassVideosReq) (http.Response, error)
	// GetMyUploadedVideos 老师查询所有上传视频
	GetMyUploadedVideos(c *gin.Context, claim ijwt.UserClaim) (http.Response, error)
}

type videoController struct {
	cos *cos.COSClient
	vs  video2.VideoService
	ps  pipeline.PipelineService
	bs  behavior.BehaviorService
}

func NewVideoController(cos *cos.COSClient, vs video2.VideoService, ps pipeline.PipelineService, bs behavior.BehaviorService) VideoController {
	return &videoController{
		cos: cos,
		vs:  vs,
		ps:  ps,
		bs:  bs,
	}
}

// UploadVideo 老师上传视频
// @Summary 上传视频文件
// @Description 教师上传教学视频，支持 mp4 格式，form-data 提交
// @Tags video
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer 认证令牌" default(Bearer )
// @Param title formData string true "视频标题"
// @Param file formData file true "视频文件"
// @Param deadline formData string true "截止日期 格式：2025-12-31 23:59:59"
// @Success 200 {object} http.Response{data=video.UploadVideoResp} "上传成功"
// @Failure 400 {object} http.Response "参数错误"
// @Failure 500 {object} http.Response "服务器错误"
// @Router /api/v1/video/upload [post]
func (vc *videoController) UploadVideo(c *gin.Context, req video.UploadVideoReq, file multipart.File, fileHeader *multipart.FileHeader, claim ijwt.UserClaim) (http.Response, error) {
	//先判断是否是老师身份
	if domain.RoleType(claim.Role) != domain.Role_Teacher {
		return http.Response{}, errors.UploadVideoError(errors1.New("无上传视频权限"))
	}
	key, err := vc.cos.UploadFile(c, file, fileHeader)
	if err != nil {
		return http.Response{}, errors.UploadVideoError(err)
	}
	deadline, err := tool.ParseStringToTime(req.Deadline)
	if err != nil {
		return http.Response{}, errors.UploadVideoError(errors1.New(fmt.Sprintf("截止日期格式错误:%v", err)))
	}
	videoID, err := vc.vs.SaveVideo(c, domain.Video{Title: req.Title, TeacherID: claim.ID, FileKey: key, Deadline: deadline})
	if err != nil {
		return http.Response{}, errors.UploadVideoError(err)
	}

	jobID, err := vc.ps.CreateJob(c, videoID)
	if err != nil {
		return http.Response{}, errors.UploadVideoError(err)
	}

	url, err := vc.cos.SignUrl(c, key)
	if err != nil {
		return http.Response{}, errors.UploadVideoError(err)
	}

	resp := video.UploadVideoResp{
		VideoID: videoID,
		JobID:   jobID,
		Title:   req.Title,
		URL:     url,
	}
	return http.Success(resp), nil
}

// GetTaskUploadingProcess 获取视频任务发送进度
// @Summary 获取视频任务发送进度
// @Description 获取视频任务发送进度
// @Tags video
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param request body video.GetTaskUploadingProcessReq true "请求参数"
// @Success 200 {object} http.Response{data=video.GetTaskUploadingProcessResp} "成功"
// @Failure 401 {object} http.Response "未授权"
// @Router /api/v1/video/task/process [post]
func (vc *videoController) GetTaskUploadingProcess(c *gin.Context, req video.GetTaskUploadingProcessReq) (http.Response, error) {
	job, err := vc.ps.GetJob(c, req.JobID)
	if err != nil {
		return http.Response{}, err
	}
	stages, err := vc.ps.ListStages(c, req.JobID)
	if err != nil {
		return http.Response{}, err
	}
	stage := ""
	if len(stages) > 0 {
		stage = stages[len(stages)-1].Stage
	}
	_ = stages // 预留：后续可把 stages 全量返回给前端
	return http.Success(video.GetTaskUploadingProcessResp{
		JobID:  job.JobID,
		Status: job.Status,
		Stage:  stage,
	}), nil
}

// GetVideoDetail 获取视频任务详情
// @Summary 获取视频详情（分段+题目）
// @Description 获取视频详情、分段、视频内题目、课后习题
// @Tags video
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param video_id 	path  int 	true 		"视频id"
// @Success 200 {object} http.Response{data=video.GetVideoDetailResp} "成功"
// @Failure 401 {object} http.Response "未授权"
// @Router	/api/v1/video/getDetail/{video_id} [get]
func (vc *videoController) GetVideoDetail(c *gin.Context, req video.GetVideoDetailReq) (http.Response, error) {
	claim, err := ginx.GetClaim(c)
	if err != nil {
		return http.Response{}, errors.GetVideoDetailError(err)
	}
	if domain.RoleType(claim.Role) == domain.Role_Student {
		v0, err := vc.vs.GetVideo(c, req.VideoID)
		if err != nil {
			return http.Response{}, errors.GetVideoDetailError(err)
		}
		if v0.ReviewStatus != "published" {
			return http.Response{}, errors.VideoNotPublishedError(errors1.New("未发布"))
		}
	}

	v, segs, kps, qs, qkps, err := vc.vs.GetVideoDetail(c, req.VideoID)
	if err != nil {
		return http.Response{}, err
	}

	url, err := vc.cos.SignUrl(c, v.FileKey)
	if err != nil {
		return http.Response{}, errors.UploadVideoError(err)
	}

	kpResp := make([]video.KnowledgePointResp, 0, len(kps))
	for _, kp := range kps {
		kpResp = append(kpResp, video.KnowledgePointResp{
			ID:          kp.ID,
			KnowledgeID: kp.KnowledgeID,
			Title:       kp.Title,
			Content:     kp.Content,
		})
	}

	quizResp := make([]question.Question, 0, len(qs))
	segmentQuiz := map[uint]question.Question{}
	for _, q := range qs {
		var opts []string
		if q.Options != nil {
			opts = q.Options
		}
		kplist := qkps[q.ID]
		kpsOut := make([]question.KnowledgePoint, 0, len(kplist))
		for _, kp := range kplist {
			kpsOut = append(kpsOut, question.KnowledgePoint{
				KnowledgeID: kp.ID,
				Title:       kp.Title,
			})
		}
		qResp := question.Question{
			ID:              q.ID,
			Type:            q.Type,
			Content:         q.Content,
			Options:         opts,
			Answer:          q.Answer,
			Analysis:        q.Analysis,
			KnowledgePoints: kpsOut,
			SegmentID:       q.SegmentID,
		}

		if q.SegmentID != nil {
			segmentQuiz[*q.SegmentID] = qResp
		} else {
			quizResp = append(quizResp, qResp)
		}
	}

	segResp := make([]video.Segment, 0, len(segs))
	for _, s := range segs {
		segResp = append(segResp, video.Segment{
			ID:        s.ID,
			SegmentID: s.SegmentID,
			Start:     s.Start,
			End:       s.End,
			Text:      s.Text,
			Question:  segmentQuiz[s.ID],
		})
	}

	return http.Success(video.GetVideoDetailResp{
		VideoID:   v.ID,
		Title:     v.Title,
		Duration:  v.Duration,
		Url:       url,
		Segments:  segResp,
		Knowledge: kpResp,
		Questions: quizResp,
	}), nil
}

// PostVideoToClass 给指定班级发送视频任务
// @Summary 下发视频任务到班级
// @Description 老师将视频作为任务下发给多个班级
// @Tags video
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param request body video.PostVideoToClassReq true "请求参数"
// @Success 200 {object} http.Response "成功"
// @Failure 401 {object} http.Response "未授权"
// @Router /api/v1/video/post-to-class [post]
func (vc *videoController) PostVideoToClass(c *gin.Context, req video.PostVideoToClassReq) (http.Response, error) {
	claim, err := ginx.GetClaim(c)
	if err != nil {
		return http.Response{}, errors.UploadVideoError(err)
	}
	if domain.RoleType(claim.Role) != domain.Role_Teacher {
		return http.Response{}, errors.UploadVideoError(errors1.New("无下发任务权限"))
	}
	if err := vc.vs.AssignVideoToClasses(c, req.VideoID, req.ClassList); err != nil {
		return http.Response{}, err
	}
	return http.Success(nil), nil
}

// GetClassVideoTasks 老师获取班级任务信息
// @Summary 获取班级视频任务列表
// @Description 老师获取班级的所有视频任务
// @Tags video
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param class_id 	path  int 	true 		"班级id"
// @Success 200 {object} http.Response{data=video.GetClassVideosResp} "成功"
// @Failure 401 {object} http.Response "未授权"
// @Router	/api/v1/video/getTasks/{class_id} [get]
func (vc *videoController) GetClassVideoTasks(c *gin.Context, req video.GetClassVideosReq) (http.Response, error) {
	vs, err := vc.vs.ListClassVideoTasks(c, req.ClassID)
	if err != nil {
		return http.Response{}, err
	}
	out := make([]video.VideoTask, 0, len(vs))
	for _, v := range vs {
		out = append(out, video.VideoTask{
			VideoID:   v.ID,
			Title:     v.Title,
			CreatedAt: v.CreatedAt,
			Deadline:  v.Deadline,
		})
	}
	return http.Success(video.GetClassVideosResp{VideoTasks: out}), nil
}

// GetMyUploadedVideos 老师查询自己上传的所有视频
// @Summary 获取老师上传的视频列表
// @Tags video
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Success 200 {object} http.Response{data=video.GetMyUploadedVideosResp} "成功"
// @Failure 401 {object} http.Response "未授权"
// @Router /api/v1/video/my-uploaded [get]
func (vc *videoController) GetMyUploadedVideos(c *gin.Context, claim ijwt.UserClaim) (http.Response, error) {
	if domain.RoleType(claim.Role) != domain.Role_Teacher {
		return http.Response{}, errors.UploadVideoError(errors1.New("无权限"))
	}
	vs, err := vc.vs.ListMyUploadedVideos(c, claim.ID)
	if err != nil {
		return http.Response{}, err
	}
	out := make([]video.VideoTask, 0, len(vs))
	for _, v := range vs {
		out = append(out, video.VideoTask{
			VideoID:   v.ID,
			Title:     v.Title,
			Status:    "",
			CreatedAt: v.CreatedAt,
			Deadline:  v.Deadline,
		})
	}
	return http.Success(video.GetMyUploadedVideosResp{Videos: out}), nil
}
