package usecase

import (
	"context"
	"example-restful-api-server/models"
	"example-restful-api-server/photogramm"
	"io"
)

type PhotogrammUsecase struct {
	photoRepoDB    photogramm.PhotoRepo
	photoRepoLocal photogramm.PhotoRepo
	albumRepo      photogramm.AlbumsRepo
}

func NewPhotogrammUsecase(db photogramm.PhotoRepo, local photogramm.PhotoRepo, albumRepo photogramm.AlbumsRepo) *PhotogrammUsecase {
	return &PhotogrammUsecase{
		photoRepoDB:    db,
		photoRepoLocal: local,
		albumRepo:      albumRepo,
	}
}

func (uc *PhotogrammUsecase) UploadPhoto(c context.Context, u *models.User, p *models.Photo, src io.Reader) (string, error) {

	if _, err := uc.photoRepoLocal.CreatePhoto(c, p, src); err != nil {
		return "", err
	}
	return uc.photoRepoDB.CreatePhoto(c, p, src)
}

func (uc *PhotogrammUsecase) GetPhoto(ctx context.Context, u *models.User, id string) (*models.Photo, error) {
	return uc.photoRepoDB.GetPhoto(ctx, u, id)
}

func (uc *PhotogrammUsecase) RemovePhoto(ctx context.Context, u *models.User, id string) error {
	p := &models.Photo{
		ID:        id,
		Username:  u.Username,
		UserID:    u.ID,
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
	for _, id := range a.PhotosID {
		p := &models.Photo{
			ID:        id,
			Username:  u.Username,
			UserID:    u.ID,
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
