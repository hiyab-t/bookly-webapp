package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"time"

	"cafe_store.hiyabnako/internal/validator"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	ScopeAuthentication = "authentication"
)

type Token struct {
	Plaintext string `json:"token"`
	Hash []byte `json:"-"`
	UserID int64 `json:"-"`
	Expiry time.Time `json:"expiry"`
	Scope string `json:"-"`
}

type TokenModel struct {
	Dbpool *pgxpool.Pool
}

func ValidateToken(v *validator.Validator, tokenPlaintext string) {

	v.Check(tokenPlaintext != "", "token", "token must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "token must be 26 bytes long")

}

func generateToken(userID int64, ttl time.Duration, scope string) *Token{

	token := &Token{
		Plaintext: rand.Text(),
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope: scope,
	}

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token
}

func (db *TokenModel) InsertToken(token *Token) error {
	qry := `insert into tokens(hash, user_id, expiry, scope) values ($1, $2, $3, $4);`

	ctx, cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()

	_, err := db.Dbpool.Exec(ctx, qry, token.Hash, token.UserID, token.Expiry, token.Scope)

	return err

}

func (db *TokenModel) DeleteAllforUser(scope string, userID int64) error {
	qry := `delete from tokens
			where scope = $1 and user_id = $2;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.Dbpool.Exec(ctx, qry, scope, userID)

	return err
}

func (db TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := generateToken(userID, ttl, scope)

	err := db.InsertToken(token)

	return token, err

}