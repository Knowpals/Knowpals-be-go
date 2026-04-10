package class

type CreateClassRequest struct {
	ClassName string `json:"class_name" binding:"required"`
}

type CreateClassResp struct {
	ClassName  string `json:"class_name"`
	ClassID    uint   `json:"class_id"`
	InviteCode string `json:"invite_code"`
}

type JoinClassRequest struct {
	InviteCode string `json:"invite_code" binding:"required"`
}

type QuitClassRequest struct {
	ClassID uint `json:"class_id" binding:"required"`
}

type GetClassInfoRequest struct {
	ClassID uint `json:"class_id" binding:"required"`
}

type ClassInfo struct {
	TeacherID   uint   `json:"teacher_id"`
	TeacherName string `json:"teacher_name"`
	ClassID     uint   `json:"class_id"`
	ClassName   string `json:"class_name"`
}

type GetClassInfoResp struct {
	ClassInfo ClassInfo `json:"class_info"`
}

type GetMyCreatedClassesResp struct {
	ClassList []ClassInfo `json:"class_list"`
}

type GetMyJoinedClassesResp struct {
	ClassList []ClassInfo `json:"class_list"`
}

type Student struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type GetClassStudentsResp struct {
	Students []Student `json:"students"`
}

type GetClassStudentsRequest struct {
	ClassID uint `json:"class_id" binding:"required"`
}
