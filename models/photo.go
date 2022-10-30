package models

import "example-restful-api-server/utils"

type Photo struct {
	ID        string `json:"photo_id"`
	Username  string
	UserID    int
	AlbumName string `json:"album_name"`
}

func NewPhoto(username string, userId int, albumName string) *Photo {
	id := utils.GetId()
	return &Photo{
		ID:        id,
		Username:  username,
		UserID:    userId,
		AlbumName: albumName,
	}
}
