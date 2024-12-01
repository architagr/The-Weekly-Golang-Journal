package dto

import "checking-error-types/pkg/entities"

type AuthRequest struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (res *AuthResponse) Init(userInfo *entities.User) *AuthResponse {
	res = new(AuthResponse)
	res.Id = userInfo.Id
	res.Name = userInfo.Name
	return res
}
