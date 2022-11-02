package dto

import "github.com/suryaadi44/eAD-System/pkg/entity"

type UserSignUpRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	NIP      string `json:"nip" validate:"omitempty,len=18"`
	NIK      string `json:"nik" Validate:"required,len=16"`
	Name     string `json:"name" Validate:"required"`
	Position string `json:"position"`
	Telp     string `json:"telp" validate:"required"`
	Sex      string `json:"sex" validate:"required"`
	Address  string `json:"address" validate:"required"`
}

func (u *UserSignUpRequest) ToEntity() *entity.User {
	return &entity.User{
		Username: u.Username,
		Password: u.Password,
		NIP:      u.NIP,
		NIK:      u.NIK,
		Name:     u.Name,
		Position: u.Position,
		Telp:     u.Telp,
		Sex:      u.Sex,
		Address:  u.Address,
	}
}

type UserLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ApplicantResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

func NewApplicantResponse(user *entity.User) *ApplicantResponse {
	return &ApplicantResponse{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
	}
}

type EmployeeResponse struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Name     string `json:"name,omitempty"`
	NIP      string `json:"nip,omitempty"`
	Position string `json:"position,omitempty"`
}

func NewEmployeeResponse(user *entity.User) *EmployeeResponse {
	return &EmployeeResponse{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		NIP:      user.NIP,
		Position: user.Position,
	}
}
