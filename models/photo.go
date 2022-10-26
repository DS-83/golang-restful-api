package models

import "example-restful-api-server/utils"

type Photo struct {
	Id        string `json:"photo_id"`
	Username  string
	UserId    int
	AlbumName string `json:"album_name"`
}

func NewPhoto(username string, userId int, albumName string) *Photo {
	id := utils.GetId()
	return &Photo{
		Id:        id,
		Username:  username,
		UserId:    userId,
		AlbumName: albumName,
	}
}
