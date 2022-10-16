package mysqldb

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"example-restful-api-server/e"
	"example-restful-api-server/models"
	"fmt"

	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int
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

	err := r.db.QueryRowContext(ctx, q, username).Scan(&u.Id, &u.Username, &u.Password)
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

func (r UserRepo) RevokeToken(ctx context.Context, key []byte) error {
	// Create hash of key string
	hasher := sha1.New()
	if _, err := hasher.Write(key); err != nil {
		log.Println(err)
		return err
	}
	hash := hex.EncodeToString(hasher.Sum(nil))

	q := "SELECT jwt_token_id FROM revoked_tokens WHERE jwt_token_id = ?"
	if err := r.db.QueryRow(q, hash).Scan(); err != sql.ErrNoRows {
		log.Println(err)
		return e.ErrRevokedToken
	}

	q = "INSERT INTO revoked_tokens (jwt_token_id) VALUES (?)"
	if _, err := r.db.ExecContext(ctx, q, hash); err != nil {
		return err
	}
	return nil
}

func (r UserRepo) IsRevoked(key []byte) bool {
	// Create hash of key string
	hasher := sha1.New()
	if _, err := hasher.Write(key); err != nil {
		log.Println(err)
		return true
	}
	hash := hex.EncodeToString(hasher.Sum(nil))

	q := "SELECT jwt_token_id FROM revoked_tokens WHERE jwt_token_id = ?"
	if err := r.db.QueryRow(q, hash).Scan(); err != sql.ErrNoRows {
		return true
	}
	return false
}
