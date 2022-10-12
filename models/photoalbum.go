package models

type PhotoAlbum struct {
	Name   string   `json:"album_name"`
	UserId int      `json:"-"`
	Photos []string `json:"photos"`
	Total  int      `json:"total"`
}

func NewAlbum(name string, userId int, photo []string, total int) PhotoAlbum {
	return PhotoAlbum{
		Name:   name,
		UserId: userId,
		Photos: photo,
		Total:  total,
	}
}
