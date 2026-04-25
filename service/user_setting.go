package service

import "numbernama-go/model"

type UserSettingService struct{}

func NewUserSettingService() *UserSettingService {
	return &UserSettingService{}
}

func (s *UserSettingService) Get() model.UserSetting {
	return model.UserSetting{Theme: "default"}
}
