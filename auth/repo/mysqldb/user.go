package mysqldb

import (
	"context"
	"database/sql"
	"example-restful-api-server/models"
	"fmt"
)

type user struct {
	ID       int
	MongoID  string
	Username string
}

type UserRepo struct {
	db           *sql.DB
	defaultAlbum string
}

func NewUserRepo(db *sql.DB, defaultAlbum string) *UserRepo {
	return &UserRepo{
		db:           db,
		defaultAlbum: defaultAlbum,
	}
}

func (r UserRepo) CreateUser(ctx context.Context, u *models.User) (*models.User, error) {
	user := toMySQLUser(u)

	q := "INSERT INTO users (username, mongo_id) VALUES (?, ?)"

	if _, err := r.db.ExecContext(ctx, q, user.Username, user.MongoID); err != nil {
		return nil, err
	}
	q = fmt.Sprintf(
		"INSERT INTO albums (album_name, user_id) SELECT '%s', id FROM users WHERE username=?",
		r.defaultAlbum,
	)
	if _, err := r.db.ExecContext(ctx, q, u.Username); err != nil {
		return nil, err
	}
	return r.GetUser(ctx, user.Username)
}

func (r UserRepo) GetUser(ctx context.Context, username string) (*models.User, error) {
	u := new(user)

	q := "SELECT id, mongo_id, username FROM users WHERE username = ?"

	err := r.db.QueryRowContext(ctx, q, username).Scan(&u.ID, &u.MongoID, &u.Username)
	if err != nil {
		return nil, err
	}

	return toModelsUser(u), nil
}

func (r UserRepo) DeleteUser(ctx context.Context, u *models.User) error {
	tables := []string{"photos", "albums"}
	for _, t := range tables {
		q := fmt.Sprintf(
			"DELETE FROM %s WHERE user_id = (SELECT id FROM users WHERE username = ?)",
			t,
		)
		if _, err := r.db.ExecContext(ctx, q, u.Username); err != nil {
			return err
		}
	}
	q := "DELETE FROM users WHERE username = ?"
	if _, err := r.db.ExecContext(ctx, q, u.Username); err != nil {
		return err
	}
	return nil
}

func (r UserRepo) UpdateUser(c context.Context, f, u *models.User) error {
	filt := toMySQLUser(f)
	upd := toMySQLUser(u)

	q := "UPDATE users SET username = ? WHERE id = ?"

	if _, err := r.db.ExecContext(c, q, upd.Username, filt.ID); err != nil {
		return err
	}
	return nil
}

func toMySQLUser(u *models.User) *user {
	return &user{
		ID:       u.ID,
		MongoID:  u.MongoID,
		Username: u.Username,
	}
}

func toModelsUser(u *user) *models.User {
	return &models.User{
		ID:       u.ID,
		MongoID:  u.MongoID,
		Username: u.Username,
	}
}
