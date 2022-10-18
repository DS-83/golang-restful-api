package models

import (
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"
)

type Photo struct {
	Id        string `json:"photo_id"`
	Username  string
	UserId    int
	AlbumName string `json:"album_name"`
}

func NewPhoto(username string, userId int, albumName string) *Photo {
	id := GetId()
	return &Photo{
		Id:        id,
		Username:  username,
		UserId:    userId,
		AlbumName: albumName,
	}
}

func GetId() string {
	n := 1000000
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(n) + rand.Intn(n*2)
	h := sha1.Sum([]byte(strconv.Itoa(i)))
	s := hex.EncodeToString(h[:])
	return s
}
