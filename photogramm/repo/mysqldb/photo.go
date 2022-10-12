package mysqldb

import (
	"context"
	"database/sql"
	"fmt"
	"fotogramm/example-restful-api-server/e"
	"fotogramm/example-restful-api-server/models"
	"io"
)

type PhotoRepo struct {
	db           *sql.DB
	defaultAlbum string
}

type AlbumRepo struct {
	db *sql.DB
}

func NewPhotoRepo(db *sql.DB, defaultAlbum string) *PhotoRepo {
	return &PhotoRepo{
		db:           db,
		defaultAlbum: defaultAlbum,
	}
}

func NewAlbumRepo(db *sql.DB) *AlbumRepo {
	return &AlbumRepo{
		db: db,
	}
}

func (r *PhotoRepo) CreatePhoto(ctx context.Context, p *models.Photo, s io.Reader) error {
	q := "INSERT INTO photos (id, album_id, user_id) VALUES (?, ?, ?)"

	if p.AlbumName == "" {
		p.AlbumName = r.defaultAlbum
	}

	albumId, err := r.GetAlbumId(p.AlbumName, p.Username)
	if err != nil {
		return err
	}

	if _, err := r.db.ExecContext(ctx, q, p.Id, albumId, p.UserId); err != nil {
		return err
	}
	return nil
}

func (r *PhotoRepo) GetPhoto(ctx context.Context, u *models.User, id string) (*models.Photo, error) {
	p := &models.Photo{}
	q := "SELECT albums.album_name FROM photos JOIN albums ON photos.user_id=albums.user_id WHERE photos.id=?"

	err := r.db.QueryRowContext(ctx, q, id).Scan(&p.AlbumName)
	if err == sql.ErrNoRows {
		return nil, e.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	p.Id = id
	p.UserId = u.Id
	p.Username = u.Username

	return p, nil

}
func (r *PhotoRepo) RemovePhoto(ctx context.Context, p *models.Photo) error {
	q := "DELETE FROM photos WHERE id=? AND user_id=?"

	_, err := r.db.ExecContext(ctx, q, p.Id, p.UserId)
	if err == sql.ErrNoRows {
		fmt.Println(err)
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

func (r *AlbumRepo) CreateAlbum(ctx context.Context, u *models.User, name string) error {
	// Check album does not exist
	q := "SELECT id FROM albums WHERE album_name=? AND user_id=?"
	var id int
	err := r.db.QueryRowContext(ctx, q, name, u.Id).Scan(&id)
	if err != sql.ErrNoRows {
		return e.ErrAlreadyExist
	}

	q = "INSERT INTO albums (album_name, user_id) VALUES (?, ?)"

	if _, err := r.db.ExecContext(ctx, q, name, u.Id); err != nil {
		return err
	}

	return nil
}

func (r *AlbumRepo) GetAlbum(ctx context.Context, u *models.User, name string) (*models.PhotoAlbum, error) {
	a := &models.PhotoAlbum{}

	q := "SELECT id FROM albums WHERE album_name=? AND user_id=?"
	var id int
	err := r.db.QueryRowContext(ctx, q, name, u.Id).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, e.ErrNotFound
	}

	q = "SELECT id FROM photos WHERE album_id=?"
	rows, err := r.db.QueryContext(ctx, q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	s, err := rowsToSlice(rows)
	if err != nil {
		return nil, err
	}

	a.Name = name
	a.Photos = *s
	a.Total = len(*s)

	return a, nil
}

func (r *AlbumRepo) RemoveAlbum(ctx context.Context, u *models.User, name string) error {
	// Check album does exist
	q := "SELECT id FROM albums WHERE album_name=? AND user_id=?"
	var id int
	err := r.db.QueryRowContext(ctx, q, name, u.Id).Scan(&id)
	if err == sql.ErrNoRows {
		return e.ErrNotFound
	}

	q = "DELETE FROM albums WHERE album_name=? AND user_id=?"

	if _, err := r.db.ExecContext(ctx, q, name, u.Id); err != nil {
		return err
	}

	return nil
}

func (r *AlbumRepo) GetInfo(ctx context.Context, u *models.User) (*models.User, error) {
	q := "SELECT album_name FROM albums WHERE user_id=?"

	rows, err := r.db.QueryContext(ctx, q, u.Id)
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
