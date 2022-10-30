package models

type PhotoAlbum struct {
	Name     string   `json:"album_name"`
	UserID   int      `json:"-"`
	PhotosID []string `json:"photos"`
	Total    int      `json:"total"`
}

func NewAlbum(name string, userId int, photo []string, total int) PhotoAlbum {
	return PhotoAlbum{
		Name:     name,
		UserID:   userId,
		PhotosID: photo,
		Total:    total,
	}
}
