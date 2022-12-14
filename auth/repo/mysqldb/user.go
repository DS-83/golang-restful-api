package mysqldb

import (
	"context"
	"database/sql"
	e "example-restful-api-server/err"
	"example-restful-api-server/models"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Username string
	Password string
}

func toMySQLUser(u *models.User) *User {
	return &User{
		Username: u.Username,
		Password: u.Password,
	}
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

func (r UserRepo) CreateUser(ctx context.Context, u *models.User) (err error) {
	sqlUser := toMySQLUser(u)

	q := "INSERT INTO users (username, password) VALUES (?, ?)"

	if _, err = r.db.ExecContext(ctx, q, sqlUser.Username, sqlUser.Password); err != nil {
		return err
	}
	q = fmt.Sprintf("INSERT INTO albums (album_name, user_id) SELECT '%s', id FROM users WHERE username=?", r.defaultAlbum)
	if _, err = r.db.ExecContext(ctx, q, u.Username); err != nil {
		return err
	}
	return nil
}

func (r UserRepo) GetUser(ctx context.Context, username string, pass string) (*models.User, error) {
	u := &models.User{}
	q := "SELECT id, username, password FROM users WHERE username = ?"

	err := r.db.QueryRowContext(ctx, q, username).Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		return nil, e.Wrap("can't auth user", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass))
	if err != nil {
		return nil, e.Wrap("invalid credentials", err)
	}
	return u, err
}

func (r UserRepo) DeleteUser(ctx context.Context, u *models.User) error {
	tables := []string{"photos", "albums"}
	for _, t := range tables {
		q := fmt.Sprintf("DELETE FROM %s WHERE user_id = (SELECT id FROM users WHERE username = ?)", t)
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
