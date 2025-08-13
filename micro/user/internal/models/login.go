package models

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	UserId string       `json:"user_id"`
	Token  string       `json:"token"`
	User   UserResponse `json:"user"`
}
