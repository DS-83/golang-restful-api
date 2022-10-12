package usecase

import (
	"context"
	"fotogramm/example-restful-api-server/models"
	"fotogramm/example-restful-api-server/photogramm"
	"io"
)

type PhotogrammUsecase struct {
	photoRepoDB    photogramm.PhotoRepo
	photoRepoLocal photogramm.PhotoRepo
	albumRepo      photogramm.PhotoAlbums
}

func NewPhotogrammUsecase(db photogramm.PhotoRepo, local photogramm.PhotoRepo, albumRepo photogramm.PhotoAlbums) *PhotogrammUsecase {
	return &PhotogrammUsecase{
		photoRepoDB:    db,
		photoRepoLocal: local,
		albumRepo:      albumRepo,
	}
}

func (uc *PhotogrammUsecase) UploadPhoto(ctx context.Context, u *models.User, albumName string, src io.Reader) error {
	p := models.NewPhoto(u.Username, u.Id, albumName)

	if err := uc.photoRepoLocal.CreatePhoto(ctx, p, src); err != nil {
		return err
	}
	return uc.photoRepoDB.CreatePhoto(ctx, p, src)
}

func (uc *PhotogrammUsecase) GetPhoto(ctx context.Context, u *models.User, id string) (*models.Photo, error) {
	return uc.photoRepoDB.GetPhoto(ctx, u, id)
}

func (uc *PhotogrammUsecase) RemovePhoto(ctx context.Context, u *models.User, id string) error {
	p := &models.Photo{
		Id:        id,
		Username:  u.Username,
		UserId:    u.Id,
		AlbumName: "",
	}

	if err := uc.photoRepoLocal.RemovePhoto(ctx, p); err != nil {
		return err
	}
	return uc.photoRepoDB.RemovePhoto(ctx, p)
}

func (uc *PhotogrammUsecase) CreateAlbum(c context.Context, u *models.User, name string) error {
	return uc.albumRepo.CreateAlbum(c, u, name)
}

func (uc *PhotogrammUsecase) GetAlbum(c context.Context, u *models.User, name string) (*models.PhotoAlbum, error) {
	return uc.albumRepo.GetAlbum(c, u, name)
}

func (uc *PhotogrammUsecase) RemoveAlbum(c context.Context, u *models.User, name string) error {
	// Get album with photo ids array
	a, err := uc.albumRepo.GetAlbum(c, u, name)
	if err != nil {
		return err
	}
	// Remove photos from local storage and DB
	for _, id := range a.Photos {
		p := &models.Photo{
			Id:        id,
			Username:  u.Username,
			UserId:    u.Id,
			AlbumName: name,
		}
		if err := uc.photoRepoLocal.RemovePhoto(c, p); err != nil {
			return err
		}
		if err := uc.photoRepoDB.RemovePhoto(c, p); err != nil {
			return err
		}
	}

	return uc.albumRepo.RemoveAlbum(c, u, name)
}

func (uc *PhotogrammUsecase) GetInfo(ctx context.Context, u *models.User) (*models.User, error) {
	return uc.albumRepo.GetInfo(ctx, u)
}
