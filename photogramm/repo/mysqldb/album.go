package mysqldb

import (
	"context"
	"database/sql"
	e "example-restful-api-server/err"
	"example-restful-api-server/models"
)

type AlbumRepo struct {
	db *sql.DB
}

type album struct {
	ID        int
	AlbumName string
	UserID    int
}

func NewAlbumRepo(db *sql.DB) *AlbumRepo {
	return &AlbumRepo{
		db: db,
	}
}

func (r *AlbumRepo) CreateAlbum(ctx context.Context, u *models.User, name string) error {
	// Check album does not exist
	q := "SELECT id FROM albums WHERE album_name=? AND user_id=?"
	var id int
	err := r.db.QueryRowContext(ctx, q, name, u.ID).Scan(&id)
	if err != sql.ErrNoRows {
		return e.ErrAlreadyExist
	}

	q = "INSERT INTO albums (album_name, user_id) VALUES (?, ?)"

	if _, err := r.db.ExecContext(ctx, q, name, u.ID); err != nil {
		return err
	}

	return nil
}

func (r *AlbumRepo) GetAlbum(ctx context.Context, u *models.User, name string) (*models.PhotoAlbum, error) {
	a := &album{
		AlbumName: name,
		UserID:    u.ID,
	}

	q := "SELECT id FROM albums WHERE album_name=? AND user_id=?"
	err := r.db.QueryRowContext(ctx, q, name, a.UserID).Scan(&a.ID)
	if err == sql.ErrNoRows {
		return nil, e.ErrNotFound
	}

	q = "SELECT id FROM photos WHERE album_id=?"
	rows, err := r.db.QueryContext(ctx, q, a.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	s, err := rowsToSlice(rows)
	if err != nil {
		return nil, err
	}

	return toModelsAlbum(a, s), nil
}

func (r *AlbumRepo) RemoveAlbum(ctx context.Context, u *models.User, name string) error {
	// Check album does exist
	q := "SELECT id FROM albums WHERE album_name=? AND user_id=?"
	var id int
	err := r.db.QueryRowContext(ctx, q, name, u.ID).Scan(&id)
	if err == sql.ErrNoRows {
		return e.ErrNotFound
	}

	q = "DELETE FROM albums WHERE album_name=? AND user_id=?"

	if _, err := r.db.ExecContext(ctx, q, name, u.ID); err != nil {
		return err
	}

	return nil
}

func (r *AlbumRepo) GetInfo(ctx context.Context, u *models.User) (*models.User, error) {
	q := "SELECT album_name FROM albums WHERE user_id=?"

	rows, err := r.db.QueryContext(ctx, q, u.ID)
	if err != nil {
		return nil, err
	}
	a, err := rowsToSlice(rows)
	if err != nil {
		return nil, err
	}
	albums := []models.PhotoAlbum{}

	for _, name := range *a {
		album, err := r.GetAlbum(ctx, u, name)
		if err != nil {
			return nil, err
		}
		albums = append(albums, *album)
	}

	u.PhotoAlbums = albums
	return u, nil
}

func toModelsAlbum(a *album, s *[]string) *models.PhotoAlbum {
	return &models.PhotoAlbum{
		Name:     a.AlbumName,
		UserID:   a.UserID,
		PhotosID: *s,
		Total:    len(*s),
	}
}
