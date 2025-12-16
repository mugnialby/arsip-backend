package request

type LoginRequest struct {
	UserId   string `json:"userId" binding:"required"`
	Password string `json:"password" binding:"required"`
}
