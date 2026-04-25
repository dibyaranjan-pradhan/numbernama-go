package service

import "numbernama-go/model"

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) Me() model.User {
	return model.User{ID: "anonymous", Name: "stub"}
}
