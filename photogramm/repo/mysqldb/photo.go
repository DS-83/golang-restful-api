package mysqldb

import (
	"context"
	"database/sql"
	e "example-restful-api-server/err"
	"example-restful-api-server/models"
	"io"
)

type PhotoRepo struct {
	db           *sql.DB
	defaultAlbum string
}

type photo struct {
	ID      string
	AlbumID int
	UserID  int
}

func NewPhotoRepo(db *sql.DB, defaultAlbum string) *PhotoRepo {
	return &PhotoRepo{
		db:           db,
		defaultAlbum: defaultAlbum,
	}
}

func (r *PhotoRepo) CreatePhoto(ctx context.Context, p *models.Photo, s io.Reader) (string, error) {
	q := "INSERT INTO photos (id, album_id, user_id) VALUES (?, ?, ?)"

	if p.AlbumName == "" {
		p.AlbumName = r.defaultAlbum
	}

	albumId, err := r.GetAlbumId(p.AlbumName, p.Username)
	if err != nil {
		return "", err
	}

	if _, err := r.db.ExecContext(ctx, q, p.ID, albumId, p.UserID); err != nil {
		return "", err
	}
	return p.ID, nil
}

func (r *PhotoRepo) GetPhoto(ctx context.Context, u *models.User, id string) (*models.Photo, error) {
	p := &photo{
		ID:     id,
		UserID: u.ID,
	}

	q := "SELECT album_id FROM photos WHERE id=?"
	err := r.db.QueryRowContext(ctx, q, p.ID).Scan(&p.AlbumID)
	if err == sql.ErrNoRows {
		return nil, e.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	q = "SELECT album_name FROM albums WHERE id=?"
	name := ""
	err = r.db.QueryRowContext(ctx, q, p.AlbumID).Scan(&name)
	if err == sql.ErrNoRows {
		return nil, e.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return toModelsPhoto(p, u.Username, name), nil

}

func (r *PhotoRepo) RemovePhoto(ctx context.Context, p *models.Photo) error {
	q := "DELETE FROM photos WHERE id=? AND user_id=?"

	_, err := r.db.ExecContext(ctx, q, p.ID, p.UserID)
	if err == sql.ErrNoRows {
		return e.ErrNotFound
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *PhotoRepo) GetAlbumId(albName, username string) (id int, err error) {
	q := "SELECT albums.id FROM albums INNER JOIN users ON albums.user_id=users.id WHERE username= ? AND album_name = ?"

	if err := r.db.QueryRow(q, username, albName).Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

func rowsToSlice(rows *sql.Rows) (*[]string, error) {
	s := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		s = append(s, id)
	}
	return &s, nil
}

func toModelsPhoto(p *photo, u string, a string) *models.Photo {
	return &models.Photo{
		ID:        p.ID,
		Username:  u,
		UserID:    p.UserID,
		AlbumName: a,
	}
}
