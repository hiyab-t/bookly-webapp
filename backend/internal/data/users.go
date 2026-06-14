package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"

	"time"

	"cafe_store.hiyabnako/internal/validator"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	User_id int `json:"user_id"`
	Create_at time.Time `json:"create_at"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password password `json:"-"`
	Active bool `json:"active"`
	Version int `json:"-"`
}

type password struct {
	plaintext *string
	hashtext []byte
}

type UsersModel struct {
	DBpool *pgxpool.Pool
}

var (
	errDuplicateEmail = errors.New("duplicate email")
)

// declaring an anonymous user var
var AnonUser = &Users{}

// check if user is anonymous
func (u *Users) IsAnon() bool {
	return u == AnonUser
}

func (p *password) Set(plaintextPass string) error {

	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPass), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPass
	p.hashtext = hash

	return nil

}

func (p *password) Matches(plaintextPass string) (bool, error) {

	err := bcrypt.CompareHashAndPassword(p.hashtext, []byte(plaintextPass))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, err
		default:
			return false, err
		}
	}

	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {

	v.Check(email != "", "email", "cannot be empty")
	v.Check(validator.Match(email,validator.EmailRX), "email", "must be a valid email address")

}

func ValidatePassword(v *validator.Validator, password string) {

	v.Check(password != "", "password", "cannot be empty")
	v.Check(len(password) <= 72, "password", "cannot be more than 72 bytes long")
	v.Check(len(password) >= 8, "password", "cannot be less than 8 bytes long")
}



func ValidateUsers(v *validator.Validator, u *Users) {
	v.Check(u.Name != "", "name", "cannot be empty")
	v.Check(len(u.Name) <= 500, "name", "must be less than 500 bytes long")

	ValidateEmail(v, u.Email)
	
	if u.Password.plaintext != nil {
		ValidatePassword(v, *u.Password.plaintext)
	}

	if u.Password.hashtext == nil {
		panic("missing password hash for user")
	}
}

func (db *UsersModel) InsertUser(u *Users) error{
	qry := `insert into users (name, email, password_hashed, active) values ($1, $2, $3, $4) returning user_id, create_at, version;`

	ctx,cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{u.Name, u.Email, u.Password.hashtext, u.Active}

	err := db.DBpool.QueryRow(ctx, qry,args...).Scan(&u.User_id, &u.Create_at, &u.Version)
	if err != nil {
		switch {
		case err.Error() == `"pq: duplicate key value violates unique constraint "users_email_key"`:
			return errDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (db *UsersModel) GetByEmail(email string) (*Users, error){

	qry := `select user_id, create_at, name, email, password_hashed, active, version
			from users
			where email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	u := &Users{}

	err := db.DBpool.QueryRow(ctx, qry, email).Scan(
		&u.User_id,
		&u.Create_at, 
		&u.Name,
		&u.Email, 
		&u.Password.hashtext, 
		&u.Active, 
		&u.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return u, nil
}

func (db *UsersModel) Update(u *Users) error {

	qry := `update users
			set name = $1,
			email = $2,
			password_hashed = $3,
			active = $4
			where user_id = $5 and version = $6;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err :=db.DBpool.Exec(ctx, qry, u.Name, u.Email, u.Password.hashtext, u.Active, u.User_id, u.Version)
	if err != nil {
		switch {
			case err.Error() == `"pq: duplicate key value violates unique constraint "users_email_key"`:
				return errDuplicateEmail
			case errors.Is(err, sql.ErrNoRows):
				return ErrRecordNotFound
			default:
				return err
		}
	}

	return nil
}

func (db TokenModel) GetUserForToken(tokenScope string, tokenPlaintext string) (*Users,error) {

	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	qry := `select u.user_id, u.create_at, u.name, u.email, u.password_hashed, u.active, u.version
			from users u
			join tokens t on u.user_id = t.user_id
			where t.hash = $1
			and t.scope = $2
			and t.expiry > $3;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	u := &Users{}

	args := []any{tokenHash[:], tokenScope, time.Now()}

	err := db.Dbpool.QueryRow(ctx, qry, args...).Scan(
		&u.User_id, 
		&u.Create_at, 
		&u.Name, 
		&u.Email, 
		&u.Password.hashtext, 
		&u.Active, 
		&u.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return u, nil

}
