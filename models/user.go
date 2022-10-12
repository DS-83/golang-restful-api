package models

import (
	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	Id          int          `json:"id"`
	Username    string       `json:"username"`
	Password    string       `json:"-"`
	PhotoAlbums []PhotoAlbum `json:"photo_albums"`
}

type NoAuthUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JwtClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type JwtToken struct {
	Token string `json:"token"`
}

func NewJwtToken(token string) JwtToken {
	return JwtToken{
		Token: token,
	}
}

func NewNoAuthUser(u, p string) NoAuthUser {
	return NoAuthUser{
		Username: u,
		Password: p,
	}
}

func NewUser(id int, uname string) User {
	return User{
		Id:          id,
		Username:    uname,
		Password:    "",
		PhotoAlbums: []PhotoAlbum{},
	}
}

// func (u *User) Info(db database.MysqlDB) (*[]byte, error) {
// 	albums, err := db.SelectAlbums(u.Id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	photos, err := db.SelectPhotos(u.Id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for k, v := range albums {
// 		pSlice := []string{}
// 		for i, j := range photos {
// 			if k == j {
// 				pSlice = append(pSlice, i)
// 			}
// 		}
// 		pAlbum := NewAlbum(v, u.Id, pSlice, len(pSlice))
// 		u.PhotoAlbums[k] = pAlbum
// 	}
// 	res, err := json.Marshal(u)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &res, nil
// }

// func (u *NoAuthUser) Auth(db database.MysqlDB) (*User, error) {
// 	id, uName, pass, err := db.SelectUser(u.Username)
// 	if err != nil {
// 		return nil, e.Wrap("can't auth user", err)
// 	}
// 	if err := bcrypt.CompareHashAndPassword([]byte(pass), []byte(u.Password)); err != nil {
// 		return nil, e.Wrap("invalid credentials", err)
// 	}
// 	return &User{
// 		Id:          id,
// 		Username:    uName,
// 		Password:    pass,
// 		PhotoAlbums: map[int]PhotoAlbum{},
// 	}, nil
// }

// func (j *JwtClaims) GetUser(db database.MysqlDB) (*User, error) {
// 	id, uName, _, err := db.SelectUser(j.Username)
// 	if err != nil {
// 		return nil, e.Wrap("getUser: can't get user from db", err)
// 	}
// 	return &User{
// 		Id:          id,
// 		Username:    uName,
// 		Password:    "",
// 		PhotoAlbums: map[int]PhotoAlbum{},
// 	}, nil
// }
