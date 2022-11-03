// nolint
package usecase

import (
	"bytes"
	"context"
	"example-restful-api-server/models"
	"example-restful-api-server/photogramm/repo/local"
	"example-restful-api-server/photogramm/repo/mock"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	repoDB    = &mock.PhotoRepoMock{}
	repoLocal = &local.PhotoLocalStorageMock{}
	albumRepo = &mock.AlbumRepoMock{}
)

var user = &models.User{
	ID:          0,
	Username:    "test",
	Password:    "test",
	PhotoAlbums: []models.PhotoAlbum{},
}

var photo = models.NewPhoto(user.Username, user.ID, "test")

var album = &models.PhotoAlbum{
	Name:     "test",
	UserID:   0,
	PhotosID: []string{},
	Total:    0,
}

var uc = NewPhotogrammUsecase(repoDB, repoLocal, albumRepo)

func TestUsecase_UploadPhoto(t *testing.T) {

	repoLocal.On("CreatePhoto", photo).Return("", nil)
	repoDB.On("CreatePhoto", photo).Return(photo.ID, nil)

	res, err := uc.UploadPhoto(context.Background(), user, photo, &bytes.Buffer{})
	assert.NoError(t, err)
	assert.Equal(t, photo.ID, res)

}

func TestUsecase_GetPhoto(t *testing.T) {

	photo.AlbumName = ""

	repoLocal.On("RemovePhoto", photo).Return(nil)
	repoDB.On("RemovePhoto", photo).Return(nil)

	err := uc.RemovePhoto(context.Background(), user, photo.ID)

	assert.NoError(t, err)

}

func TestUsecase_RemovePhoto(t *testing.T) {

	uc := NewPhotogrammUsecase(repoDB, repoLocal, albumRepo)

	photo.AlbumName = ""
	repoLocal.On("RemovePhoto", photo).Return(nil)
	repoDB.On("RemovePhoto", photo).Return(nil)
	err := uc.RemovePhoto(context.Background(), user, photo.ID)
	assert.NoError(t, err)

}

func TestUsecase_CreateAlbum(t *testing.T) {
	albumRepo.On("CreateAlbum", user, "test").Return(nil)
	err := uc.CreateAlbum(context.Background(), user, "test")
	assert.NoError(t, err)

}
func TestUsecase_GetAlbum(t *testing.T) {
	albumRepo.On("GetAlbum", user, "test").Return(album, nil)
	res, err := uc.GetAlbum(context.Background(), user, "test")
	assert.NoError(t, err)
	assert.Equal(t, album, res)

}
func TestUsecase_RemoveAlbum(t *testing.T) {
	albumRepo.On("GetAlbum", user, "test").Return(album, nil)
	repoLocal.On("RemovePhoto", photo).Return(nil)
	repoDB.On("RemovePhoto", photo).Return(nil)

	albumRepo.On("RemoveAlbum", user, "test").Return(nil)
	err := uc.RemoveAlbum(context.Background(), user, "test")
	assert.NoError(t, err)
}
func TestUsecase_GetInfo(t *testing.T) {
	albumRepo.On("GetInfo", user).Return(user, nil)

	res, err := uc.GetInfo(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, user, res)
}
