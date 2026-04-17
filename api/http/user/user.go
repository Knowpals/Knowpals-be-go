package user

type SendCodeRequest struct {
	Email string `json:"email"`
}

type SendCodeResp struct {
	Code string `json:"code"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Code     string `json:"code"`
}

type RegisterResp struct {
	Token string `json:"token"`
}

type LoginByPassword struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginByVerifyCode struct {
	Email      string `json:"email"`
	VerifyCode string `json:"verify_code"`
}

type LoginResp struct {
	Token string `json:"token"`
}

type GetUserRequest struct {
	ID uint `uri:"id" biding:"required"`
}

type ForgotPasswordReq struct {
	Email       string `json:"email" binding:"required"`
	VerifyCode  string `json:"verify_code" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type GetUserInfoResp struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	CreatedAt   string `json:"created_at"`
	RegisterDay int    `json:"register_day"`
}
