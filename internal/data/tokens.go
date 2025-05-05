package data

import (
	"database/sql"
	"time"

	"github.com/mf751/gocha/internal/validator"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

type Token struct {
	PlainText string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:""`
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
