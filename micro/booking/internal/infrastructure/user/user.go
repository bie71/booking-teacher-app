package user

import (
	"booking/internal/config"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type UserResponse struct {
	Id           uint      `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	ProfileImage string    `json:"profile_image"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserService struct {
	restyClient *resty.Client
	service     config.Service
}

func NewUserService(restyClient *resty.Client, service config.Service) *UserService {
	return &UserService{
		restyClient: restyClient,
		service:     service,
	}
}

func (s *UserService) GetUserById(userId uint) (UserResponse, error) {
	url := fmt.Sprintf("%s/api/v1/profile", s.service.Host)

	var userResponse UserResponse
	resp, err := s.restyClient.R().
		SetResult(&userResponse).
		SetQueryParam("user_id", fmt.Sprintf("%d", userId)).
		Get(url)

	if err != nil {
		return UserResponse{}, err
	}

	if resp.StatusCode() != http.StatusOK {
		return UserResponse{}, fmt.Errorf("user not found")
	}

	return userResponse, nil

}
