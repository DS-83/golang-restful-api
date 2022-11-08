package models

type User struct {
	ID          int          `json:"id"`
	MongoID     string       `json:"-"`
	Username    string       `json:"username"`
	Password    string       `json:"-"`
	PhotoAlbums []PhotoAlbum `json:"photo_albums"`
}
