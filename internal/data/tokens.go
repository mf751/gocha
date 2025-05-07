package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"

	"github.com/google/uuid"

	"github.com/mf751/gocha/internal/validator"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

type Token struct {
	PlainText string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    uuid.UUID `json:""`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

type TokenModel struct {
	DB *sql.DB
}

func ValidateTokenPlainText(vdtr *validator.Validator, tokenPlainText string) {
	vdtr.Check(tokenPlainText != "", "token", "must be provided")
	vdtr.Check(len(tokenPlainText) == 26, "token", "must be 26 bytes long")
}

func generateToken(userID uuid.UUID, timeToLive time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(timeToLive),
		Scope:  scope,
	}

	// generate random 16 bytes
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]

	return token, nil
}

func (model TokenModel) New(
	userID uuid.UUID,
	timeToLive time.Duration,
	scope string,
) (*Token, error) {
	token, err := generateToken(userID, timeToLive, scope)
	if err != nil {
		return nil, err
	}

	err = model.Insert(token)
	return token, err
}

func (model TokenModel) Insert(token *Token) error {
	sqlQuery := `
INSERT INTO tokens (hash, user_id, expiry, scope)
VALUES($1, $2, $3, $4)
	`

	args := []interface{}{token.Hash, token.UserID, token.Expiry, token.Scope}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := model.DB.ExecContext(ctx, sqlQuery, args...)
	return err
}

func (model TokenModel) DeleteAllForUser(scope string, userID uuid.UUID) error {
	sqlQuery := `
DELETE FROM tokens
WHERE scope = $1 AND user_id = $2
  `

	args := []interface{}{scope, userID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := model.DB.ExecContext(ctx, sqlQuery, args...)
	return err
}
