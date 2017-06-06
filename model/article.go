package model

type Article struct {
	Id 				string `json:"id" binding:"required"`
	ThumbnailUrl 	string `json:"thumbnailUrl" binding:"required"`
	Title 			string `json:"title" binding:"required"`
	CreatedDate		string `json:"createdDate" binding:"required"`
	Author 			string `json:"author" binding:"required"`
	Content 		string `json:"content" binding:"required"`
}

type Comment struct {
	Id 				string `json:"id"`
	Time 			string `json:"time"`
	Content 		string `json:"content" binding:"required"`
	Username 		string `json:"username"`
}