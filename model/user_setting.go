package model

type UserSetting struct {
	UserID string `json:"userId,omitempty"`
	Theme  string `json:"theme,omitempty"`
}
