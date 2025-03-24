package DTO

type LoginRequestDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponseDTO struct {
	AccessToken string `json:"access_token"`
	UserID      int64  `json:"user_id"`
	Username    string `json:"username"`
}

type SignUpRequestDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email" binding:"required,email"`
}
