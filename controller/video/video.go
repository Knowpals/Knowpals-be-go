package video

import (
	errors1 "errors"
	"mime/multipart"

	"github.com/Knowpals/Knowpals-be-go/api/http"
	"github.com/Knowpals/Knowpals-be-go/api/http/video"
	"github.com/Knowpals/Knowpals-be-go/domain"
	"github.com/Knowpals/Knowpals-be-go/errors"
	"github.com/Knowpals/Knowpals-be-go/infra/cos"
	"github.com/Knowpals/Knowpals-be-go/pkg/ijwt"
	"github.com/Knowpals/Knowpals-be-go/service/pipeline"
	video2 "github.com/Knowpals/Knowpals-be-go/service/video"
	"github.com/gin-gonic/gin"
)

type VideoController interface {
	// UploadVideo 老师上传视频
	UploadVideo(c *gin.Context, req video.UploadVideoReq, file multipart.File, fileHeader *multipart.FileHeader, claim ijwt.UserClaim) (http.Response, error)
	GetTaskUploadingProcess(c *gin.Context, req video.UploadVideoReq) (http.Response, error)
	// GetVideoDetail 获取视频任务详情
	GetVideoDetail(c *gin.Context, req video.GetVideoDetailReq) (http.Response, error)
	// PostVideoToClass 给指定班级发送视频任务
	PostVideoToClass(c *gin.Context, req video.PostVideoToClassReq) (http.Response, error)
	// GetClassVideoTasks 学生获取班级任务信息
	GetClassVideoTasks(c *gin.Context, req video.GetClassVideosReq) (http.Response, error)
}

type videoController struct {
	cos *cos.COSClient
	vs  video2.VideoService
	ps  pipeline.PipelineService
}

func NewVideoController(cos *cos.COSClient, vs video2.VideoService, ps pipeline.PipelineService) VideoController {
	return &videoController{
		cos: cos,
		vs:  vs,
		ps:  ps,
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

	videoID, err := vc.vs.SaveVideo(c, domain.Video{Title: req.Title, TeacherID: claim.ID, FileKey: key})
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

func (vc *videoController) GetTaskUploadingProcess(c *gin.Context, req video.UploadVideoReq) (http.Response, error) {
	//TODO implement me
	panic("implement me")
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
	return http.Response{}, nil
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
	return http.Response{}, nil
}

// GetClassVideoTasks 学生获取班级任务信息
// @Summary 获取班级视频任务列表
// @Description 学生获取所在班级的所有视频任务
// @Tags video
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer )
// @Param class_id 	path  int 	true 		"班级id"
// @Success 200 {object} http.Response{data=video.GetClassVideosResp} "成功"
// @Failure 401 {object} http.Response "未授权"
// @Router	/api/v1/video/getTasks/{class_id} [get]
func (vc *videoController) GetClassVideoTasks(c *gin.Context, req video.GetClassVideosReq) (http.Response, error) {
	return http.Response{}, nil
}
