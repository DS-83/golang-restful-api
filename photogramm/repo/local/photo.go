package local

import (
	"context"
	e "example-restful-api-server/err"
	"example-restful-api-server/models"
	"io"
	"os"
	"path/filepath"
)

const defaultPerm = 0774

type PhotoRepo struct {
	basePath string
}

func NewPhotoRepo(basePath string) *PhotoRepo {
	return &PhotoRepo{
		basePath: basePath,
	}
}

func (r *PhotoRepo) CreatePhoto(ctx context.Context, p *models.Photo, s io.Reader) (string, error) {
	filePath := filepath.Join(r.basePath, p.Username)

	if err := os.MkdirAll(filePath, defaultPerm); err != nil {
		return "", err
	}

	fileName := p.ID
	if fileName == "" {
		return "", e.ErrIncorrectFileName
	}

	filePath = filepath.Join(filePath, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(file, s); err != nil {
		return "", err
	}
	return "", nil

}

func (r *PhotoRepo) GetPhoto(c context.Context, u *models.User, id string) (*models.Photo, error) {
	return nil, nil

}
func (r *PhotoRepo) RemovePhoto(ctx context.Context, p *models.Photo) error {
	filePath := filepath.Join(r.basePath, p.Username, p.ID)

	err := os.Remove(filePath)
	if os.IsNotExist(err) {
		return e.ErrNotFound
	}
	if err != nil {
		return err
	}

	return nil
}
