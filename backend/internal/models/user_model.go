package models

type UserModel struct {
	Id          int    `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	SocialId    string `json:"socialId"`
}
