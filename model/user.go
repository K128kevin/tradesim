package model

type User struct {
	Id				string
	Username 		string
	Email			string
}

type NewUser struct {
	Username 		string `form:"Username" json:"Username" binding:"required"`
	Password		string `form:"Password" json:"Password" binding:"required"`
	Email			string `form:"Email" json:"Email" binding:"required"`
}

type UpdatePassword struct {
	OldPassword 	string `form:"OldPassword" json:"OldPassword" binding:"required"`
	NewPassword 	string `form:"NewPassword" json:"NewPassword" binding:"required"`
}

type Login struct {
	Username 		string `form:"Username" json:"Username" binding:"required"`
	Password 		string `form:"Password" json:"Password" binding:"required"`
}