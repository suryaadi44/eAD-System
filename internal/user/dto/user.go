package dto

import "github.com/suryaadi44/eAD-System/pkg/entity"

type UserSignUpRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	NIP      string `json:"nip" validate:"omitempty,len=18"`
	NIK      string `json:"nik" Validate:"required,len=16"`
	Name     string `json:"name" Validate:"required"`
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
		Telp:     u.Telp,
		Sex:      u.Sex,
		Address:  u.Address,
	}
}

type UserLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
